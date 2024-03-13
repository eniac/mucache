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

func writePlot(ctx context.Context, req *movie.WritePlotRequest) *movie.WritePlotResponse {
	plotId := movie.WritePlot(ctx, req.PlotId, req.Plot)
	//fmt.Println("Movie info stored for id: " + movieId)
	resp := movie.WritePlotResponse{PlotId: plotId}
	return &resp
}

func readPlot(ctx context.Context, req *movie.ReadPlotRequest) *movie.ReadPlotResponse {
	plot := movie.ReadPlot(ctx, req.PlotId)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := movie.ReadPlotResponse{Plot: plot}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/write_plot", wrappers.NonROWrapper[movie.WritePlotRequest, movie.WritePlotResponse](writePlot))
	http.HandleFunc("/ro_read_plot", wrappers.ROWrapper[movie.ReadPlotRequest, movie.ReadPlotResponse](readPlot))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
