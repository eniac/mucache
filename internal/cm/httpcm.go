package cm

import (
	"github.com/goccy/go-json"
	"fmt"
	"github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/common"
	"github.com/golang/glog"
	"net/http"
)

func startHandler(w http.ResponseWriter, req *http.Request) {
	var startReq cm.StartRequest
	err := json.NewDecoder(req.Body).Decode(&startReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cm.WQ <- startReq
}

func endHandler(w http.ResponseWriter, req *http.Request) {
	var endReq cm.EndRequest
	err := json.NewDecoder(req.Body).Decode(&endReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cm.WQ <- endReq
}

func invHandler(w http.ResponseWriter, req *http.Request) {
	var invReq cm.InvalidateKeyRequest
	err := json.NewDecoder(req.Body).Decode(&invReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cm.WQ <- invReq
}

func invCallsHandler(w http.ResponseWriter, req *http.Request) {
	var request cm.InvalidateCallsRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cm.WQ <- request
}

func saveCallHandler(w http.ResponseWriter, req *http.Request) {
	var request cm.SaveCallsRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cm.WQ <- request
}

func ServeHttp(cfg *Config) {
	if !common.ZMQ {
		http.HandleFunc(cm.HttpStartSuffix, startHandler)
		http.HandleFunc(cm.HttpEndSuffix, endHandler)
		//http.HandleFunc(cm.HttpInvSuffix, invHandler)
	}
	http.HandleFunc(cm.HttpInvSuffix, invHandler)
	http.HandleFunc(cm.HttpInvCallsSuffix, invCallsHandler)
	http.HandleFunc(cm.HttpSaveCallsSuffix, saveCallHandler)

	glog.Info("Listening to port: 80")
	panic(http.ListenAndServe(fmt.Sprintf(":%d", 80), nil))
}
