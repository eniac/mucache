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

func compose(ctx context.Context, req *movie.ComposeRequest) *movie.ComposeResponse {
	ok := movie.Compose(ctx, req.Username, req.Password, req.Title, req.Rating, req.Text)
	//fmt.Printf("Page read: %v\n", page)
	resp := movie.ComposeResponse{Ok: ok}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/compose", wrappers.NonROWrapper[movie.ComposeRequest, movie.ComposeResponse](compose))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
