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

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func read(ctx context.Context, req *twoserivces.ReadRequest) *twoserivces.ReadResponse {
	resp := invoke.Invoke[twoserivces.ReadResponse](ctx, "callee", "ro_read", req)
	return &resp
}

func readBulk(ctx context.Context, req *twoserivces.ReadBulkRequest) *twoserivces.ReadBulkResponse {
	resp := invoke.Invoke[twoserivces.ReadBulkResponse](ctx, "callee", "ro_read_bulk", req)
	return &resp
}

func write(ctx context.Context, req *twoserivces.WriteRequest) *string {
	resp := invoke.Invoke[string](ctx, "callee", "write", req)
	return &resp
}

func writeBulk(ctx context.Context, req *twoserivces.WriteBulkRequest) *string {
	resp := invoke.Invoke[string](ctx, "callee", "write_bulk", req)
	return &resp
}

func hit(ctx context.Context, req *twoserivces.ReadRequest) *string {
	invoke.InvokeHit(ctx, "callee", "ro_read", req)
	resp := "OK"
	return &resp
}

func miss(ctx context.Context, req *twoserivces.ReadRequest) *string {
	invoke.InvokeMiss[twoserivces.ReadResponse](ctx, "callee", "ro_read", req)
	resp := "OK"
	return &resp
}

func hitormiss(ctx context.Context, req *twoserivces.HitOrMissRequest) *string {
	dice := rand.Float32()
	if dice < req.HitRate {
		invoke.InvokeHit(ctx, "callee", "ro_read", req)
	} else {
		invoke.InvokeMiss[twoserivces.ReadResponse](ctx, "callee", "ro_read", req)
	}
	resp := "OK"
	return &resp
}

func invalidationExperiment(ctx context.Context, req *loadcm.InvalidationExperimentRequest) *string {
	fmt.Printf("Starting invalidation experiment for: %v \n", req.Times)
	go twoserivces.InvalidationExperiment(req.Times, req.Timeout, "callee", "ro_read", "callee", "write")
	resp := "OK"
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(1))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/read", wrappers.ROWrapper[twoserivces.ReadRequest, twoserivces.ReadResponse](read))
	http.HandleFunc("/read_bulk", wrappers.ROWrapper[twoserivces.ReadBulkRequest, twoserivces.ReadBulkResponse](readBulk))
	http.HandleFunc("/write", wrappers.NonROWrapper[twoserivces.WriteRequest, string](write))
	http.HandleFunc("/write_bulk", wrappers.NonROWrapper[twoserivces.WriteBulkRequest, string](writeBulk))
	http.HandleFunc("/hit", wrappers.ROWrapper[twoserivces.ReadRequest, string](hit))
	http.HandleFunc("/miss", wrappers.ROWrapper[twoserivces.ReadRequest, string](miss))
	http.HandleFunc("/ro_hitormiss", wrappers.ROWrapper[twoserivces.HitOrMissRequest, string](hitormiss))
	http.HandleFunc("/invalidation_experiment", wrappers.NonROWrapper[loadcm.InvalidationExperimentRequest, string](invalidationExperiment))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
