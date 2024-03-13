package wrappers

import (
	"fmt"
	"github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/utility"
	"sync"
)

// Global variables can be accessed by all threads in Go
// Use a sync.Map to make it thread-safe
type Deps struct {
	// We are now using CallIds instead of call arguments to index this
	// dependency map. Using callArgs complicates reasoning when there are multiple concurrent calls
	// with the same arguments.
	deps sync.Map
}

// Just for debugging
func (deps *Deps) String() string {
	var s string
	deps.deps.Range(func(key, value interface{}) bool {
		s += fmt.Sprintf("%v: %v\n", key, value)
		return true
	})
	return s
}

func (deps *Deps) InitDep(id cm.CallId) {
	//deps.deps[id] = make(map[cm.Key]struct{})
	//deps.deps.Store(id, make(map[cm.Key]struct{}))
	m := &sync.Map{}
	deps.deps.Store(id, m)
}

func (deps *Deps) AddKeyDep(id cm.CallId, key cm.Key) {
	//idDeps, ok := deps.deps[id]
	//// The call must exist already (it has to have been initialized
	//utility.Assert(ok)
	//idDeps[key] = struct{}{}
	idDeps, ok := deps.deps.Load(id)
	// The call must exist already (it has to have been initialized
	utility.Assert(ok)
	//idDeps.(map[cm.Key]struct{})[key] = struct{}{}
	idDeps.(*sync.Map).Store(key, struct{}{})
}

func (deps *Deps) AddCallDep(id cm.CallId, ca cm.CallArgs) {
	idDeps, ok := deps.deps.Load(id)
	// The call must exist already (it has to have been initialized
	utility.Assert(ok)
	idDeps.(*sync.Map).Store(ca, struct{}{})
}

func (deps *Deps) PopDeps(id cm.CallId) ([]cm.Key, []cm.CallArgs) {
	//idDeps, ok := deps.deps[id]
	//// The call must exist already (it has to have been initialized
	//utility.Assert(ok)
	//// Pop this call id from the dependency map
	//delete(deps.deps, id)
	idDeps, ok := deps.deps.Load(id)
	// The call must exist already (it has to have been initialized
	utility.Assert(ok)
	// Pop this call id from the dependency map
	deps.deps.Delete(id)
	//return idDeps.(map[cm.Key]struct{})
	var keyDeps []cm.Key
	var callDeps []cm.CallArgs
	idDeps.(*sync.Map).Range(func(key, value interface{}) bool {
		switch key.(type) {
		case cm.Key:
			keyDeps = append(keyDeps, key.(cm.Key))
		case cm.CallArgs:
			callDeps = append(callDeps, key.(cm.CallArgs))
		default:
			panic(fmt.Sprintf("Unknown type: %+v in deps map", key))
		}
		return true
	})
	return keyDeps, callDeps
}

//func DepsMapToSlice(deps map[cm.Key]struct{}) []cm.Key {
//	keys := make([]cm.Key, len(deps))
//
//	i := 0
//	for k := range deps {
//		keys[i] = k
//		i++
//	}
//	return keys
//}

// make it a pointer so that it's not accidentally copied
var deps = &Deps{}
