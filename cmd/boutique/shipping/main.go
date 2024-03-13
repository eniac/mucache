package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/boutique"
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

func getQuote(ctx context.Context, req *boutique.GetQuoteRequest) *boutique.GetQuoteResponse {
	quote := boutique.GetQuote(ctx, req.Items)
	resp := boutique.GetQuoteResponse{CostUsd: quote}
	return &resp
}

func shipOrder(ctx context.Context, req *boutique.ShipOrderRequest) *boutique.ShipOrderResponse {
	id := boutique.ShipOrder(ctx, req.Address, req.Items)
	resp := boutique.ShipOrderResponse{TrackingId: id}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_get_quote", wrappers.ROWrapper[boutique.GetQuoteRequest, boutique.GetQuoteResponse](getQuote))
	http.HandleFunc("/ship_order", wrappers.NonROWrapper[boutique.ShipOrderRequest, boutique.ShipOrderResponse](shipOrder))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
