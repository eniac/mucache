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

func charge(ctx context.Context, req *boutique.ChargeRequest) *boutique.ChargeResponse {
	uid, err := boutique.Charge(ctx, req.Amount, req.CreditCard)
	//fmt.Printf("Products read: %+v\n", products)
	resp := boutique.ChargeResponse{
		Uuid:  uid,
		Error: err,
	}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/charge", wrappers.NonROWrapper[boutique.ChargeRequest, boutique.ChargeResponse](charge))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
