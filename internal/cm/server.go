package cm

import (
	"github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/common"
)

// These two functions do the communication with the caller cache manager
func sendSaveToCaller(cfg *Config, caller cm.ServiceName, callArgs cm.CallArgs, retVal cm.ReturnVal) {
	//glog.Infof(" -- Sending to: %v (Save-Cache-Line, ca: %v -> %v)\n", caller, callArgs, retVal)
	req := cm.SaveCallsRequest{CallArgsList: []cm.CallArgs{callArgs}, ReturnVals: []cm.ReturnVal{retVal}}
	request := HttpSendSaveCallsRequest{
		Request: req,
		Caller:  caller,
	}
	HttpQueue <- request
}

func sendInvsToCallers(cfg *Config, saved CallsAndCallers) {
	//glog.Infof(" -- Sending Invalidate line: %v\n", saved)
	request := HttpSendInvalidateCallRequest{InvSet: saved.Dict}
	HttpQueue <- request
}

func Process(cfg *Config, state *State) {
	//profileState := initProfile(cfg.printTimeFreq)
	for {
		//if len(cm.WQ) >= cm.QueueSize/2 {
		//	glog.Warningf("Queue size: %v\n", len(cm.WQ))
		//}
		//profileState.profileProcRequest()
		request := <-cm.WQ
		switch request.(type) {
		case cm.StartRequest:
			startRequest := request.(cm.StartRequest)
			//glog.Infof("Processing: %v\n", &startRequest)
			state.appendCall(startRequest.CallArgs)
		case cm.EndRequest:
			endRequest := request.(cm.EndRequest)
			//glog.Infof("Processing: %v\n", &endRequest)
			if state.validCall(endRequest.CallArgs, endRequest.KeyDeps, endRequest.CallDeps) {
				sendSaveToCaller(cfg, endRequest.Caller, endRequest.CallArgs, endRequest.ReturnVal)
				state.storeDeps(cfg, endRequest.Caller, endRequest.CallArgs, endRequest.KeyDeps, endRequest.CallDeps)
			} else {
				// Only for debugging
				//glog.Infof("Call", endRequest.CallArgs, "not valid!")
			}
		case cm.InvalidateKeyRequest:
			invRequest := request.(cm.InvalidateKeyRequest)
			if common.ShardEnabled && !invRequest.FromCM {
				neighbors := cfg.GetNeighbors()
				//glog.Infof("neighbors: %v\n", neighbors)
				invRequest.FromCM = true
				for _, neighbor := range neighbors {
					cm.SendInvRequestHttp(&invRequest, neighbor)
				}
			}
			//glog.Infof("Processing: %v\n", &invRequest)
			state.appendWrite(invRequest.Key)
			saved, exists := state.popKeyDeps(invRequest.Key)
			if exists {
				sendInvsToCallers(cfg, saved)
			}
		case cm.InvalidateCallsRequest:
			invRequest := request.(cm.InvalidateCallsRequest)
			//glog.Infof("Processing: %v\n", &invRequest)
			for _, ca := range invRequest.Calls {
				state.appendInvCall(&ca)
			}

			// Transitively invalidate upstream entries
			callDeps, exists := state.popCallDeps(invRequest.Calls)
			if exists {
				sendInvsToCallers(cfg, callDeps)
			}

			// Remove the call entries from the local cache
			cm.CacheRemoveCalls(cfg.cacheClient, invRequest.Calls)
		case cm.SaveCallsRequest:
			saveRequest := request.(cm.SaveCallsRequest)
			//glog.Infof("Processing: %v, Len: %v\n", &saveRequest, len(saveRequest.CallArgsList))
			cm.CacheSaveCalls(cfg.cacheClient, saveRequest.CallArgsList, saveRequest.ReturnVals)
		default:
			panic("Unreachable")
		}
		//fmt.Println("State after processing:", request)
		//fmt.Println(state.String())
	}
}
