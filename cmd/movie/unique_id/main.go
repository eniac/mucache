package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/movie"
	"github.com/eniac/mucache/pkg/cm"
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

func getUniqueId(ctx context.Context, req *movie.GetUniqueIdRequest) *movie.GetUniqueIdResponse {
	reviewId := movie.GetUniqueId(ctx, req.ReqId)
	//fmt.Printf("Page read: %v\n", page)
	resp := movie.GetUniqueIdResponse{ReviewId: reviewId}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/get_unique_id", wrappers.NonROWrapper[movie.GetUniqueIdRequest, movie.GetUniqueIdResponse](getUniqueId))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
