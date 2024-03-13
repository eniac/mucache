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

func storeMovieInfo(ctx context.Context, req *movie.StoreMovieInfoRequest) *movie.StoreMovieInfoResponse {
	movieId := movie.StoreMovieInfo(ctx, req.MovieId, req.Info, req.CastIds, req.PlotId)
	//fmt.Println("Movie info stored for id: " + movieId)
	resp := movie.StoreMovieInfoResponse{MovieId: movieId}
	return &resp
}

func readMovieInfo(ctx context.Context, req *movie.ReadMovieInfoRequest) *movie.ReadMovieInfoResponse {
	movieInfo := movie.ReadMovieInfo(ctx, req.MovieId)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := movie.ReadMovieInfoResponse{Info: movieInfo}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/store_movie_info", wrappers.NonROWrapper[movie.StoreMovieInfoRequest, movie.StoreMovieInfoResponse](storeMovieInfo))
	http.HandleFunc("/ro_read_movie_info", wrappers.ROWrapper[movie.ReadMovieInfoRequest, movie.ReadMovieInfoResponse](readMovieInfo))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
