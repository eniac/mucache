package cm

import (
	"fmt"
	"github.com/eniac/mucache/pkg/utility"
	"github.com/goccy/go-json"
	"strings"
)

type Request interface {
	// TODO: Is there a better way (perf-wise) than to return a buffer?
	Unparse() []byte
	ToJson() []byte
}

// TODO: Can we use some standard serialization and deserialization instead of
//       the abominations below??
// The code below is both slow and non-general!

func (request *StartRequest) Unparse() []byte {
	string := fmt.Sprintf("S %s\n", request.CallArgs)
	return []byte(string)
}

func (request *StartRequest) ToJson() []byte {
	newData, err := json.Marshal(request)
	utility.Assert(err == nil)
	return newData
}

// TODO: Fix this hacky-af serialization!!
func unparseKeyDependencies(keys []Key) string {
	stringKeys := []string{}
	for _, key := range keys {
		stringKeys = append(stringKeys, string(key))
	}
	return strings.Join(stringKeys, ",")
}

func unparseCallDependencies(keys []CallArgs) string {
	stringKeys := []string{}
	for _, key := range keys {
		stringKeys = append(stringKeys, string(key))
	}
	return strings.Join(stringKeys, ",")
}

func (request *EndRequest) Unparse() []byte {
	string := fmt.Sprintf("E %s %s %s %s %s\n",
		request.CallArgs,
		request.Caller,
		unparseKeyDependencies(request.KeyDeps),
		unparseCallDependencies(request.CallDeps),
		request.ReturnVal)
	return []byte(string)
}

func (request *EndRequest) ToJson() []byte {
	newData, err := json.Marshal(request)
	utility.Assert(err == nil)
	return newData
}

func (request *InvalidateKeyRequest) Unparse() []byte {
	string := fmt.Sprintf("I %s\n", request.Key)
	return []byte(string)
}

func (request *InvalidateKeyRequest) ToJson() []byte {
	newData, err := json.Marshal(request)
	utility.Assert(err == nil)
	return newData
}

// TODO: Fix this hacky-af serialization!!
func unparseCalls(calls []CallArgs) string {
	stringCalls := []string{}
	for _, ca := range calls {
		stringCalls = append(stringCalls, string(ca))
	}
	return strings.Join(stringCalls, ",")
}

func unparseReturnVals(rets []ReturnVal) string {
	stringRets := []string{}
	for _, ca := range rets {
		stringRets = append(stringRets, string(ca))
	}
	return strings.Join(stringRets, ",")
}

func (request *InvalidateCallsRequest) Unparse() []byte {
	string := fmt.Sprintf("IC %s\n", unparseCalls(request.Calls))
	return []byte(string)
}

func (request *InvalidateCallsRequest) ToJson() []byte {
	newData, err := json.Marshal(request)
	utility.Assert(err == nil)
	return newData
}

func (request *SaveCallsRequest) Unparse() []byte {
	string := fmt.Sprintf("SC %s %s\n",
		unparseCalls(request.CallArgsList),
		unparseReturnVals(request.ReturnVals))
	return []byte(string)
}

func (request *SaveCallsRequest) ToJson() []byte {
	newData, err := json.Marshal(request)
	utility.Assert(err == nil)
	return newData
}

func parseKeyDependencies(keysString string) []Key {
	var keys []Key
	for _, stringKey := range strings.Split(keysString, ",") {
		keys = append(keys, Key(stringKey))
	}
	return keys
}

func parseCalls(callsString string) []CallArgs {
	calls := []CallArgs{}
	for _, stringCall := range strings.Split(callsString, ",") {
		calls = append(calls, CallArgs(stringCall))
	}
	return calls
}

func parseReturnVals(retValsString string) []ReturnVal {
	rets := []ReturnVal{}
	for _, stringCall := range strings.Split(retValsString, ",") {
		rets = append(rets, ReturnVal(stringCall))
	}
	return rets
}

func ParseBytes(buf []byte) interface{} {
	// TODO: Do the parsing more efficiently
	line := string(buf[:])
	lineNoSuffix := strings.TrimSuffix(line, "\n")
	tokens := strings.Split(lineNoSuffix, " ")
	switch tokens[0] {
	case "S":
		return StartRequest{CallArgs: CallArgs(tokens[1])}
	case "E":
		return EndRequest{
			CallArgs:  CallArgs(tokens[1]),
			Caller:    ServiceName(tokens[2]),
			KeyDeps:   parseKeyDependencies(tokens[3]),
			CallDeps:  parseCalls(tokens[4]),
			ReturnVal: ReturnVal(tokens[5])}
	case "I":
		return InvalidateKeyRequest{Key: Key(tokens[1])}
	case "IC":
		return InvalidateCallsRequest{Calls: parseCalls(tokens[1])}
	case "SC":
		return SaveCallsRequest{CallArgsList: parseCalls(tokens[1]), ReturnVals: parseReturnVals(tokens[2])}
	default:
		panic("Not known parsing")
	}
}
