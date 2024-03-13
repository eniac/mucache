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

func composeReview(ctx context.Context, req *movie.ComposeReviewRequest) *movie.ComposeReviewResponse {
	movie.ComposeReview(ctx, req.Review)
	//fmt.Printf("Page read: %v\n", page)
	resp := movie.ComposeReviewResponse{Ok: "OK"}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/compose_review", wrappers.NonROWrapper[movie.ComposeReviewRequest, movie.ComposeReviewResponse](composeReview))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
