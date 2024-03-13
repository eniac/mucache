package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/singleservice"
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

func read(ctx context.Context, req *singleservice.ReadRequest) *singleservice.ReadResponse {
	v, err := state.GetState[string](ctx, req.K)
	if err != nil {
		panic(err)
	}
	resp := singleservice.ReadResponse{V: v}
	return &resp
}

func write(ctx context.Context, req *singleservice.WriteRequest) *string {
	state.SetState(ctx, req.K, req.V)
	resp := "OK"
	return &resp
}

func echo(ctx context.Context, req *singleservice.ReadRequest) *singleservice.ReadResponse {
	resp := singleservice.ReadResponse{V: req.K}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(1))
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_read", wrappers.ROWrapper[singleservice.ReadRequest, singleservice.ReadResponse](read))
	http.HandleFunc("/echo", wrappers.ROWrapper[singleservice.ReadRequest, singleservice.ReadResponse](echo))
	http.HandleFunc("/write", wrappers.NonROWrapper[singleservice.WriteRequest, string](write))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
