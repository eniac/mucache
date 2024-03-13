package cm

import (
	"github.com/eniac/mucache/pkg/cm"
	"github.com/golang/glog"
	"time"
)

type HttpSendSaveCallsRequest struct {
	Request cm.SaveCallsRequest
	Caller  cm.ServiceName
}

type HttpSendInvalidateCallRequest struct {
	InvSet map[cm.ServiceName]map[cm.CallArgs]struct{}
}

var HttpQueue = make(chan interface{}, cm.QueueSize)

// TODO: Move those to arguments of the cache manager
var HttpBatchLimit = 1
var HttpTimeoutDuration = time.Microsecond * 500

func HttpSender(cfg *Config) {
	glog.Infof("Started HttpSender\n")
	buffer := HttpSendBuffer{
		Callers: make(map[cm.ServiceName]HttpBufferElement),
	}
	sinceLastFlush := 0
	for {
		if len(HttpQueue) >= cm.QueueSize/2 {
			glog.Warningf("HttpSender queue size: %v\n", len(HttpQueue))
		}
		select {
		case request := <-HttpQueue:
			buffer.addRequest(request)
			sinceLastFlush++
			if sinceLastFlush > HttpBatchLimit {
				buffer.flushBuffer(cfg)
				sinceLastFlush = 0
			}
		case <-time.After(HttpTimeoutDuration):
			buffer.flushBuffer(cfg)
			sinceLastFlush = 0
		}
	}
}

type HttpBufferElement struct {
	InvSet  cm.CallArgSet
	SaveMap map[cm.CallArgs]cm.ReturnVal
}

type HttpSendBuffer struct {
	Callers map[cm.ServiceName]HttpBufferElement
}

func (buffer *HttpSendBuffer) addRequest(request interface{}) {
	switch request.(type) {
	case HttpSendSaveCallsRequest:
		req := request.(HttpSendSaveCallsRequest)
		buffer.addSaveCallsRequest(req)
	case HttpSendInvalidateCallRequest:
		req := request.(HttpSendInvalidateCallRequest)
		buffer.addInvalidateRequest(req)
	}
}

func (buffer *HttpSendBuffer) getOrInitBufEl(caller cm.ServiceName) HttpBufferElement {
	bufferEl, ok := buffer.Callers[caller]
	// If the caller does not exist in the buffer we need to add it
	if !ok {
		bufferEl = HttpBufferElement{
			InvSet:  cm.MakeCallArgSet(),
			SaveMap: make(map[cm.CallArgs]cm.ReturnVal),
		}
		buffer.Callers[caller] = bufferEl
	}
	return bufferEl
}

func (buffer *HttpSendBuffer) addSaveCallsRequest(request HttpSendSaveCallsRequest) {
	caller := request.Caller
	callArgList := request.Request.CallArgsList
	returnVals := request.Request.ReturnVals

	// Get the relevant buffer element
	bufferEl := buffer.getOrInitBufEl(caller)

	for i := range callArgList {
		ca := callArgList[i]
		ret := returnVals[i]

		// Add to the save map
		bufferEl.SaveMap[ca] = ret
		// Delete from the inv map
		bufferEl.InvSet.PopItemIfExists(ca)
	}
}

func (buffer *HttpSendBuffer) addInvalidateRequest(request HttpSendInvalidateCallRequest) {
	for caller, callArgsSet := range request.InvSet {
		// Get the relevant buffer element
		bufferEl := buffer.getOrInitBufEl(caller)
		for ca := range callArgsSet {
			// Add to the inv map
			bufferEl.InvSet.AddItem(ca)
			// Delete from the save map
			delete(bufferEl.SaveMap, ca)
		}
	}
}

func saveMapToLists(saveMap map[cm.CallArgs]cm.ReturnVal) ([]cm.CallArgs, []cm.ReturnVal) {
	callArgsList := make([]cm.CallArgs, len(saveMap))
	returnVals := make([]cm.ReturnVal, len(saveMap))
	i := 0
	for ca, ret := range saveMap {
		callArgsList[i] = ca
		returnVals[i] = ret
		i++
	}
	return callArgsList, returnVals
}

func (buffer *HttpSendBuffer) flushBuffer(cfg *Config) {
	for caller := range buffer.Callers {
		bufferEl := buffer.Callers[caller]
		// Remove this item from the buffer
		delete(buffer.Callers, caller)

		addr := cfg.GetCacheManagerAddress(caller)
		// The save calls and invalidate should always be disjoint
		// so it doesn't matter which is sent first
		invList := bufferEl.InvSet.ToList()
		if len(invList) > 0 {
			invReq := cm.InvalidateCallsRequest{Calls: invList}
			cm.SendInvCallsRequestHttp(&invReq, addr)
		}

		callArgsList, returnVals := saveMapToLists(bufferEl.SaveMap)
		if len(callArgsList) > 0 {
			saveReq := cm.SaveCallsRequest{
				CallArgsList: callArgsList,
				ReturnVals:   returnVals,
			}
			cm.SendSaveCallsRequestHttp(&saveReq, addr)
		}
	}
}
