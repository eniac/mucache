package cm

import (
	"fmt"
	"strings"
)

// TODO: Add an interface that is common between CallArgs and Keys and that could be used
//       to store them in the cache manager state (since both keys and calls can be invalidated
//       and therefore need to be stored in pending and as dependencies in saved as well as
//       invalidate other calls).

// TODO: Should we use embeddings? (https://go.dev/doc/effective_go#embedding)

const (
	TypeStartRequest     = "S"
	TypeEndRequest       = "E"
	TypeInvRequest       = "I"
	TypeInvCallsRequest  = "IC"
	TypeSaveCallsRequest = "SC"
)

// Method call name and its arguments
type CallArgs string // TODO: What is the right type for this? Some form of hash?

func (ca *CallArgs) ToString() string {
	return string(*ca)
}

func (ca *CallArgs) IsCallArgSet() bool { return false }
func (ca *CallArgs) IsWriteKey() bool   { return false }
func (ca *CallArgs) IsInvCall() bool    { return true }

// ServiceName name (to be able to identify and send them invalidation and save requests)
type ServiceName string

// The return value of a function
// TODO: could we make this be arbitrary json/arbitrary type?
type ReturnVal string

func (retVal ReturnVal) ToByteArray() []byte {
	return []byte(retVal)
}

func ByteArrayToRetVal(bs []byte) ReturnVal {
	return ReturnVal(bs)
}

// Used to differentiate between calls with the same arguments
type CallId string // Do we need another type for this?

// Database keys (that are read and written by programs)
type Key string

func (key Key) IsCallArgSet() bool { return false }
func (key Key) IsWriteKey() bool   { return true }
func (key Key) IsInvCall() bool    { return false }

// A set of calls (with their arguments)
type CallArgSet struct {
	CAS map[CallArgs]struct{}
}

func (cas *CallArgSet) IsCallArgSet() bool { return true }
func (cas *CallArgSet) IsWriteKey() bool   { return false }
func (cas *CallArgSet) IsInvCall() bool    { return false }

func (cas *CallArgSet) String() string {
	keyStrings := []string{}
	for k, _ := range cas.CAS {
		keyString := fmt.Sprintf("%+v", k)
		keyStrings = append(keyStrings, keyString)
	}
	return fmt.Sprintf("{%s}", strings.Join(keyStrings, ","))
}

func (cas *CallArgSet) AddItem(ca CallArgs) {
	cas.CAS[ca] = struct{}{}
}

func (cas *CallArgSet) PopItemIfExists(ca CallArgs) {
	delete(cas.CAS, ca)
}

func (cas *CallArgSet) ToList() []CallArgs {
	keys := make([]CallArgs, len(cas.CAS))
	i := 0
	for k := range cas.CAS {
		keys[i] = k
		i++
	}
	return keys
}

func (cas *CallArgSet) Extend(callArgList []CallArgs) {
	for _, ca := range callArgList {
		cas.AddItem(ca)
	}
}

func MakeCallArgSet() CallArgSet {
	return CallArgSet{make(map[CallArgs]struct{})}
}

// TODO: Must create a type that is common for keys and callargs
//       and is used to store them everywhere in the state of the cache manager.
//       This must be a small, like a hash.

// The three types of requests that the Cache Manager can process
type StartRequest struct {
	CallArgs CallArgs `json:"callargs"`
}

func (request *StartRequest) String() string {
	return fmt.Sprintf("Start(%v)", request.CallArgs)
}

type EndRequest struct {
	CallArgs  CallArgs    `json:"callargs"`
	Caller    ServiceName `json:"caller"`
	KeyDeps   []Key       `json:"key_deps"`
	CallDeps  []CallArgs  `json:"call_deps"`
	ReturnVal ReturnVal   `json:"returnval"`
}

func (request *EndRequest) String() string {
	return fmt.Sprintf("End(%+v, %+v, %+v, %+v, %+v)", request.CallArgs, request.Caller, request.KeyDeps, request.CallDeps, request.ReturnVal)
}

// We don't need two invalidate requests (start and end),
// we just need to handle both the start and the end the exact same way, add them in pending,
// and check whether there are any saved that depend on them.
type InvalidateKeyRequest struct {
	Key    Key  `json:"key"`
	FromCM bool `json:"fromcm,omitempty"`
}

func (request *InvalidateKeyRequest) String() string {
	return fmt.Sprintf("Inv(%v)", request.Key)
}

type InvalidateCallsRequest struct {
	Calls []CallArgs `json:"calls"`
}

func (request *InvalidateCallsRequest) String() string {
	return fmt.Sprintf("InvCalls(%+v)", request.Calls)
}

type SaveCallsRequest struct {
	CallArgsList []CallArgs  `json:"callargslist"`
	ReturnVals   []ReturnVal `json:"returnvals"`
}

func (request *SaveCallsRequest) String() string {
	return fmt.Sprintf("SaveCall(%+v, %+v)", request.CallArgsList, request.ReturnVals)
}
