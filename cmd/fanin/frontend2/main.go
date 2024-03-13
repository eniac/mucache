package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/loadcm"
	"github.com/eniac/mucache/internal/twoservices"
	"github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/invoke"
	"github.com/eniac/mucache/pkg/wrappers"
	"math/rand"
	"net/http"
	"runtime"
)

var Callee = "backend"
var MaxProcs = 8

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func read(ctx context.Context, req *twoserivces.ReadRequest) *twoserivces.ReadResponse {
	var resp twoserivces.ReadResponse
	resp = invoke.Invoke[twoserivces.ReadResponse](ctx, Callee, "ro_read", req)
	return &resp
}

func write(ctx context.Context, req *twoserivces.WriteRequest) *string {
	resp := invoke.Invoke[string](ctx, Callee, "write", req)
	return &resp
}

func hitormiss(ctx context.Context, req *twoserivces.HitOrMissRequest) *string {
	dice := rand.Float32()
	if dice < req.HitRate {
		invoke.InvokeHit(ctx, Callee, "ro_read", req)
	} else {
		invoke.InvokeMiss[twoserivces.ReadResponse](ctx, Callee, "ro_read", req)
	}
	resp := "OK"
	return &resp
}

func invalidationExperiment(ctx context.Context, req *loadcm.InvalidationExperimentRequest) *string {
	// Start running the zmqfeeder
	fmt.Printf("Starting experiment for: %v \n", req.Times)
	// TODO: Fix to invalidate a right service
	//go twoserivces.InvalidationExperiment(req.Times, req.Timeout, Callee, "ro_read", "backend", "write")
	resp := "OK"
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(MaxProcs))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_read", wrappers.ROWrapper[twoserivces.ReadRequest, twoserivces.ReadResponse](read))
	http.HandleFunc("/write", wrappers.NonROWrapper[twoserivces.WriteRequest, string](write))
	http.HandleFunc("/ro_hitormiss", wrappers.ROWrapper[twoserivces.HitOrMissRequest, string](hitormiss))
	http.HandleFunc("/invalidation_experiment", wrappers.NonROWrapper[loadcm.InvalidationExperimentRequest, string](invalidationExperiment))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
