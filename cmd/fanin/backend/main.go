package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/twoservices"
	"github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/state"
	"github.com/eniac/mucache/pkg/wrappers"
	"net/http"
	"runtime"
)

var MaxProcs = 8

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func read(ctx context.Context, req *twoserivces.ReadRequest) *twoserivces.ReadResponse {
	//req.K += 1
	v, err := state.GetState[int](ctx, fmt.Sprint(req.K))
	if err != nil {
		v = 0
	}
	resp := twoserivces.ReadResponse{V: v}
	return &resp
}

func write(ctx context.Context, req *twoserivces.WriteRequest) *string {
	state.SetState(ctx, fmt.Sprint(req.K), req.V)
	resp := "OK"
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(MaxProcs))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_read", wrappers.ROWrapper[twoserivces.ReadRequest, twoserivces.ReadResponse](read))
	http.HandleFunc("/write", wrappers.NonROWrapper[twoserivces.WriteRequest, string](write))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
