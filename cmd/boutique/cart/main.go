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

func addItemToCart(ctx context.Context, req *boutique.AddItemRequest) *boutique.AddItemResponse {
	ok := boutique.AddItem(ctx, req.UserId, req.ProductId, req.Quantity)
	resp := boutique.AddItemResponse{Ok: ok}
	return &resp
}

func getCart(ctx context.Context, req *boutique.GetCartRequest) *boutique.GetCartResponse {
	cart := boutique.GetCart(ctx, req.UserId)
	resp := boutique.GetCartResponse{Cart: cart}
	return &resp
}

func emptyCart(ctx context.Context, req *boutique.EmptyCartRequest) *boutique.EmptyCartResponse {
	ok := boutique.EmptyCart(ctx, req.UserId)
	resp := boutique.EmptyCartResponse{Ok: ok}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/add_item", wrappers.NonROWrapper[boutique.AddItemRequest, boutique.AddItemResponse](addItemToCart))
	http.HandleFunc("/ro_get_cart", wrappers.ROWrapper[boutique.GetCartRequest, boutique.GetCartResponse](getCart))
	http.HandleFunc("/empty_cart", wrappers.NonROWrapper[boutique.EmptyCartRequest, boutique.EmptyCartResponse](emptyCart))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
