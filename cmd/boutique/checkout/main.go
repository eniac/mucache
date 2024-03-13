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

func placeOrder(ctx context.Context, req *boutique.PlaceOrderRequest) *boutique.PlaceOrderResponse {
	order := boutique.PlaceOrder(ctx, req.UserId, req.UserCurrency, req.Address, req.Email, req.CreditCard)
	resp := boutique.PlaceOrderResponse{Order: order}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/place_order", wrappers.NonROWrapper[boutique.PlaceOrderRequest, boutique.PlaceOrderResponse](placeOrder))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
