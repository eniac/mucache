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

var Callees = []string{"backend1", "backend2", "backend3", "backend4"}
var MaxProcs = 8

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func read(ctx context.Context, req *twoserivces.ReadRequest) *twoserivces.ReadResponse {
	var resp twoserivces.ReadResponse
	for _, callee := range Callees {
		resp = invoke.Invoke[twoserivces.ReadResponse](ctx, callee, "ro_read", req)
	}
	return &resp
}

func asyncRead(ctx context.Context, req *twoserivces.ReadRequest) *twoserivces.ReadResponse {
	ch := make(chan twoserivces.ReadResponse, len(Callees))
	for _, callee := range Callees {
		go func() {
			r := invoke.Invoke[twoserivces.ReadResponse](ctx, callee, "ro_read", req)
			ch <- r
		}()
	}
	var resp twoserivces.ReadResponse
	for _ = range Callees {
		resp = <-ch
	}
	return &resp
}

func write(ctx context.Context, req *twoserivces.WriteRequest) *string {
	var resp string
	for _, callee := range Callees {
		// TODO: Make them async
		resp = invoke.Invoke[string](ctx, callee, "write", req)
	}
	return &resp
}

func hitormiss(ctx context.Context, req *twoserivces.HitOrMissRequest) *string {
	for _, callee := range Callees {
		dice := rand.Float32()
		if dice < req.HitRate {
			invoke.InvokeHit(ctx, callee, "ro_read", req)
		} else {
			invoke.InvokeMiss[twoserivces.ReadResponse](ctx, callee, "ro_read", req)
		}
	}
	resp := "OK"
	return &resp
}

func asyncHitormiss(ctx context.Context, req *twoserivces.HitOrMissRequest) *string {
	ch := make(chan string, len(Callees))
	for _, callee := range Callees {
		go func() {
			dice := rand.Float32()
			if dice < req.HitRate {
				invoke.InvokeHit(ctx, callee, "ro_read", req)
			} else {
				invoke.InvokeMiss[twoserivces.ReadResponse](ctx, callee, "ro_read", req)
			}
			ch <- "OK"
		}()
	}
	var resp string
	for _ = range Callees {
		resp = <-ch
	}
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
	http.HandleFunc("/ro_async_read", wrappers.ROWrapper[twoserivces.ReadRequest, twoserivces.ReadResponse](asyncRead))
	http.HandleFunc("/write", wrappers.NonROWrapper[twoserivces.WriteRequest, string](write))
	http.HandleFunc("/ro_hitormiss", wrappers.ROWrapper[twoserivces.HitOrMissRequest, string](hitormiss))
	http.HandleFunc("/ro_async_hitormiss", wrappers.ROWrapper[twoserivces.HitOrMissRequest, string](asyncHitormiss))
	http.HandleFunc("/invalidation_experiment", wrappers.NonROWrapper[loadcm.InvalidationExperimentRequest, string](invalidationExperiment))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
