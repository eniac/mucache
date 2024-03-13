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

func registerMovieId(ctx context.Context, req *movie.RegisterMovieIdRequest) *movie.RegisterMovieIdResponse {
	movie.RegisterMovieId(ctx, req.Title, req.MovieId)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := movie.RegisterMovieIdResponse{Ok: "OK"}
	return &resp
}

func getMovieId(ctx context.Context, req *movie.GetMovieIdRequest) *movie.GetMovieIdResponse {
	movieId := movie.GetMovieId(ctx, req.Title)
	resp := movie.GetMovieIdResponse{MovieId: movieId}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/register_movie_id", wrappers.NonROWrapper[movie.RegisterMovieIdRequest, movie.RegisterMovieIdResponse](registerMovieId))
	http.HandleFunc("/ro_get_movie_id", wrappers.ROWrapper[movie.GetMovieIdRequest, movie.GetMovieIdResponse](getMovieId))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
