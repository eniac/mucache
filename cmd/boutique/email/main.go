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

func sendEmail(ctx context.Context, req *boutique.SendOrderConfirmationRequest) *boutique.SendOrderConfirmationResponse {
	ok := boutique.SendConfirmation(ctx, req.Email, req.Order)
	resp := boutique.SendOrderConfirmationResponse{Ok: ok}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_send_email", wrappers.ROWrapper[boutique.SendOrderConfirmationRequest, boutique.SendOrderConfirmationResponse](sendEmail))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
