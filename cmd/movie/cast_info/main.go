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

func storeCastInfo(ctx context.Context, req *movie.StoreCastInfoRequest) *movie.StoreCastInfoResponse {
	movieId := movie.StoreCastInfo(ctx, req.CastId, req.Name, req.Info)
	//fmt.Println("Movie info stored for id: " + movieId)
	resp := movie.StoreCastInfoResponse{CastId: movieId}
	return &resp
}

func readCastInfos(ctx context.Context, req *movie.ReadCastInfosRequest) *movie.ReadCastInfosResponse {
	castInfos := movie.ReadCastInfos(ctx, req.CastIds)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := movie.ReadCastInfosResponse{Infos: castInfos}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/store_cast_info", wrappers.NonROWrapper[movie.StoreCastInfoRequest, movie.StoreCastInfoResponse](storeCastInfo))
	http.HandleFunc("/ro_read_cast_infos", wrappers.ROWrapper[movie.ReadCastInfosRequest, movie.ReadCastInfosResponse](readCastInfos))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
