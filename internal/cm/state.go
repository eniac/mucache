package cm

import (
	"encoding/gob"
	"fmt"
	"github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/utility"
	"github.com/golang/glog"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"time"
)

// TODO: Maybe it makes sense to have both writes and callargs sets!
type HistoryItem interface {
	IsCallArgSet() bool
	IsWriteKey() bool
	IsInvCall() bool
}

type CallsAndCallers struct {
	Dict map[cm.ServiceName]map[cm.CallArgs]struct{}
}

func (cac1 *CallsAndCallers) union(cac2 *CallsAndCallers) {
	for caller, calls := range cac2.Dict {
		addCallSetToCAC(cac1, caller, calls)
	}
}

type State struct {
	history []HistoryItem // Contains a sequence of started calls and key writes
	//keyDeps  map[cm.Key]CallsAndCallers
	//callDeps map[cm.CallArgs]CallsAndCallers
	keyDeps  expirable.LRU[cm.Key, CallsAndCallers]
	callDeps expirable.LRU[cm.CallArgs, CallsAndCallers]
}

var ToEvict = CallsAndCallers{Dict: make(map[cm.ServiceName]map[cm.CallArgs]struct{})}
var ToEvictSize = 1000
var UserCacheSize = 80

//var ToEvict2 = CallsAndCallers{Dict: make(map[cm.ServiceName]map[cm.CallArgs]struct{})}

func (state *State) Init() {
	state.history = make([]HistoryItem, 0)
	evict := false
	if evict {
		onEvict1 := expirable.EvictCallback[cm.Key, CallsAndCallers](func(key cm.Key, value CallsAndCallers) {
			ToEvict.union(&value)
			if len(ToEvict.Dict) > ToEvictSize {
				sendInvsToCallers(nil, ToEvict)
				ToEvict = CallsAndCallers{Dict: make(map[cm.ServiceName]map[cm.CallArgs]struct{})}
			}
		})
		onEvict2 := expirable.EvictCallback[cm.CallArgs, CallsAndCallers](func(key cm.CallArgs, value CallsAndCallers) {
			ToEvict.union(&value)
			if len(ToEvict.Dict) > ToEvictSize {
				sendInvsToCallers(nil, ToEvict)
				ToEvict = CallsAndCallers{Dict: make(map[cm.ServiceName]map[cm.CallArgs]struct{})}
			}
		})
		state.keyDeps = *expirable.NewLRU[cm.Key, CallsAndCallers](UserCacheSize*100/2, onEvict1, 0)
		state.callDeps = *expirable.NewLRU[cm.CallArgs, CallsAndCallers](UserCacheSize*100/2, onEvict2, 0)
	} else {
		state.keyDeps = *expirable.NewLRU[cm.Key, CallsAndCallers](0, nil, 0)
		state.callDeps = *expirable.NewLRU[cm.CallArgs, CallsAndCallers](0, nil, 0)
	}
	//state.keyDeps = make(map[cm.Key]CallsAndCallers)
	//state.callDeps = make(map[cm.CallArgs]CallsAndCallers)
}

func NewState() *State {
	state := &State{}
	state.Init()
	gob.Register(cm.Key(""))
	gob.Register(cm.CallArgs(""))
	gob.Register(cm.CallArgSet{})
	//go GetSizeProcess(state)
	return state
}

// Just for debugging
func (state *State) String() string {
	return fmt.Sprintf("{ History: %v,\n  KeyInvSet: %v\n  CallInvSet: %v }", state.history, state.keyDeps, state.callDeps)
}

// TODO: What to do if call already exists?
func (state *State) appendCall(ca cm.CallArgs) {
	// If the list is empty
	if len(state.history) == 0 {
		cas := cm.MakeCallArgSet()
		cas.AddItem(ca)
		state.history = append(state.history, &cas)
	} else {
		last := state.history[len(state.history)-1]
		if last.IsCallArgSet() {
			// Modify the last set
			lastCas := last.(*cm.CallArgSet)
			lastCas.AddItem(ca)
		} else {
			// Last item was a write, add a new one
			cas := cm.MakeCallArgSet()
			cas.AddItem(ca)
			state.history = append(state.history, &cas)
		}
	}
}

func (state *State) appendWrite(k cm.Key) {
	state.history = append(state.history, k)
}

func (state *State) appendInvCall(ca *cm.CallArgs) {
	state.history = append(state.history, ca)
}

// Q: Does this become slower if the slice and map become generic?
//
//	Maybe because of the map allocation
func depsSliceToSet(slice []cm.Key) map[cm.Key]struct{} {
	set := make(map[cm.Key]struct{})
	for _, el := range slice {
		set[el] = struct{}{}
	}
	return set
}

func callArgsSliceToSet(slice []cm.CallArgs) map[cm.CallArgs]struct{} {
	set := make(map[cm.CallArgs]struct{})
	for _, el := range slice {
		set[el] = struct{}{}
	}
	return set
}

var LogStart = time.Now()
var GcStart = time.Now()

// TODO: rename checkIfCallValidAndRemoveFromHistory
func (state *State) validCall(ca cm.CallArgs, keyDeps []cm.Key, callDeps []cm.CallArgs) bool {
	// Walk the history backwards until you find the call start
	// While walking, keep track of whether we have found a dependency invalidation
	history := state.history
	keyDepsSet := depsSliceToSet(keyDeps)
	callDepsSet := callArgsSliceToSet(callDeps)

	if time.Now().Sub(LogStart) >= 7*time.Second {
		glog.Info("History Size: ", len(state.history), utility.GetRealSizeOf(state.history))
		glog.Info("Deps Size: ", state.keyDeps.Len(), state.callDeps.Len(), utility.GetRealSizeOf(state.keyDeps.Values())+utility.GetRealSizeOf(state.callDeps.Values())+utility.GetRealSizeOf(state.keyDeps.Keys())+utility.GetRealSizeOf(state.callDeps.Keys()))
		LogStart = time.Now()
	}

	valid := true

	for i := len(history) - 1; i >= 0; i-- {
		historyItem := history[i]
		if historyItem.IsWriteKey() {
			writeKey := historyItem.(cm.Key)
			_, inSet := keyDepsSet[writeKey]
			// An invalidation of a dependency was found, and therefore the call is invalid
			if inSet {
				valid = false
				//return false
			}
		} else if historyItem.IsCallArgSet() {
			// Stop if we find the start of the call
			cas := historyItem.(*cm.CallArgSet)
			_, inSet := cas.CAS[ca]
			// The start of the call was found, and therefore we can
			// safely exit and the call is valid.
			if inSet {
				delete(cas.CAS, ca)
				//if len(cas.CAS) == 0 {
				//	GcCounter += 1
				//}
				//if GcCounter >= 1000 {
				if time.Now().Sub(GcStart) >= 5*time.Second {
					j := 0
					for ; j < len(history); j++ {
						if history[j].IsCallArgSet() {
							cas := history[j].(*cm.CallArgSet)
							if len(cas.CAS) != 0 {
								break
							}
						}
					}
					//glog.Infof("GC: %v items in history\n", j)
					state.history = history[j:]
					GcStart = time.Now()
				}
				// TODO: Before exiting, we also need to remove the call to ensure that memory stops growing
				//       To do that, we need to traverse history from the other way, to amortize the cost of deletions
				break
				//return true
			}
		} else {
			invCa := historyItem.(*cm.CallArgs)
			_, inSet := callDepsSet[*invCa]
			// An invalidation of a dependency was found, and therefore the call is invalid
			if inSet {
				valid = false
				//return false
			}
		}
	}
	// Note: Normally this should not be reachable, but leaving it like this
	//       so that it is easier to design a load generator.
	return valid
	//return true
	// This must never be reachable (except if we start dropping requests)
	//panic("Reached beginning of history without finding call!")
}

func addCallSetToCAC(callsAndCallers *CallsAndCallers, caller cm.ServiceName, calls map[cm.CallArgs]struct{}) {
	// If the caller does not exist, create a set for it
	_, ok2 := callsAndCallers.Dict[caller]
	if !ok2 {
		callsAndCallers.Dict[caller] = make(map[cm.CallArgs]struct{})
	}

	// Insert the new call args to the set
	for ca, _ := range calls {
		callsAndCallers.Dict[caller][ca] = struct{}{}
	}
}

func addCalltoCAC(callsAndCallers *CallsAndCallers, caller cm.ServiceName, ca cm.CallArgs) {
	callSet := make(map[cm.CallArgs]struct{})
	callSet[ca] = struct{}{}
	addCallSetToCAC(callsAndCallers, caller, callSet)
}

func (state *State) storeKeyDeps(cfg *Config, caller cm.ServiceName, ca cm.CallArgs, keys []cm.Key) {
	for _, key := range keys {
		//deps, ok := state.keyDeps[key]
		deps, ok := state.keyDeps.Get(key)
		// If the key does not exist, create a dictionary for it
		if !ok {
			deps = CallsAndCallers{make(map[cm.ServiceName]map[cm.CallArgs]struct{})}
		}
		addCalltoCAC(&deps, caller, ca)
		//state.keyDeps[key] = deps
		state.keyDeps.Add(key, deps)
	}
}

func (state *State) storeCallDeps(cfg *Config, caller cm.ServiceName, ca cm.CallArgs, callDeps []cm.CallArgs) {
	for _, callDep := range callDeps {
		//deps, ok := state.callDeps[callDep]
		deps, ok := state.callDeps.Get(callDep)
		// TODO: This initialization is copy-pasted from the one above, but I can't make it cleaner
		// If the key does not exist, create a dictionary for it
		if !ok {
			deps = CallsAndCallers{make(map[cm.ServiceName]map[cm.CallArgs]struct{})}
		}
		addCalltoCAC(&deps, caller, ca)
		//state.callDeps[callDep] = deps
		state.callDeps.Add(callDep, deps)
	}
}

func (state *State) storeDeps(cfg *Config, caller cm.ServiceName, ca cm.CallArgs, keys []cm.Key, callDeps []cm.CallArgs) {
	state.storeKeyDeps(cfg, caller, ca, keys)
	state.storeCallDeps(cfg, caller, ca, callDeps)
}

// Returns a map of keyDeps calls (the caller services and the call_args)
func (state *State) popKeyDeps(key cm.Key) (CallsAndCallers, bool) {
	//deps, ok := state.keyDeps[key]
	deps, ok := state.keyDeps.Get(key)
	if ok {
		// Note: This delete potentially leaves inside instances of keyDeps.
		//       For example, if a call read keys k1, k2 and it was keyDeps,
		//       and later k1 was invalidated, the keyDeps part of k2 will always stay in the keyDeps
		//       dictionary until k1 is also invalidated.
		//
		// TODO: We might want to clean the dependencies dictionary from these stale keyDeps at some point.
		//delete(state.keyDeps, key)
		state.keyDeps.Remove(key)
	}
	return deps, ok
}

func (state *State) popCallDeps(calls []cm.CallArgs) (CallsAndCallers, bool) {
	totalCallsAndCallers := CallsAndCallers{
		Dict: make(map[cm.ServiceName]map[cm.CallArgs]struct{}),
	}
	notEmpty := false
	for _, call := range calls {
		//callsAndCallers, ok := state.callDeps[call]
		callsAndCallers, ok := state.callDeps.Get(call)
		if ok {
			notEmpty = true
			//delete(state.callDeps, call)
			state.callDeps.Remove(call)
			totalCallsAndCallers.union(&callsAndCallers)
		}
	}
	return totalCallsAndCallers, notEmpty
}
