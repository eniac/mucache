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
		panic(err)
	}
	resp := twoserivces.ReadResponse{V: v}
	return &resp
}

func readBulk(ctx context.Context, req *twoserivces.ReadBulkRequest) *twoserivces.ReadBulkResponse {
	keys := make([]string, len(req.Ks))
	for i, k := range req.Ks {
		keys[i] = fmt.Sprint(k)
	}
	vs, err := state.GetBulkState[int](ctx, keys)
	if err != nil {
		panic(err)
	}
	resp := twoserivces.ReadBulkResponse{Vs: vs}
	return &resp
}

func write(ctx context.Context, req *twoserivces.WriteRequest) *string {
	state.SetState(ctx, fmt.Sprint(req.K), req.V)
	resp := "OK"
	return &resp
}

func writeBulk(ctx context.Context, req *twoserivces.WriteBulkRequest) *string {
	kvs := make(map[string]interface{}, len(req.Ks))
	for i, k := range req.Ks {
		kvs[fmt.Sprint(k)] = req.Vs[i]
	}
	state.SetBulkState(ctx, kvs)
	resp := "OK"
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(1))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_read", wrappers.ROWrapper[twoserivces.ReadRequest, twoserivces.ReadResponse](read))
	http.HandleFunc("/ro_read_bulk", wrappers.ROWrapper[twoserivces.ReadBulkRequest, twoserivces.ReadBulkResponse](readBulk))
	http.HandleFunc("/write", wrappers.NonROWrapper[twoserivces.WriteRequest, string](write))
	http.HandleFunc("/write_bulk", wrappers.NonROWrapper[twoserivces.WriteBulkRequest, string](writeBulk))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
