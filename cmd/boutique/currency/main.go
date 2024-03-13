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

func setCurrency(ctx context.Context, req *boutique.SetCurrencySupportRequest) *boutique.SetCurrencySupportResponse {
	ok := boutique.SetCurrencySupport(ctx, req.Currency)
	resp := boutique.SetCurrencySupportResponse{Ok: ok}
	return &resp
}

func getCurrencies(ctx context.Context, req *boutique.GetSupportedCurrenciesRequest) *boutique.GetSupportedCurrenciesResponse {
	currencies := boutique.GetSupportedCurrencies(ctx)
	resp := boutique.GetSupportedCurrenciesResponse{Currencies: currencies}
	return &resp
}

func convertCurrency(ctx context.Context, req *boutique.ConvertCurrencyRequest) *boutique.ConvertCurrencyResponse {
	amount := boutique.ConvertCurrency(ctx, req.Amount, req.ToCurrency)
	resp := boutique.ConvertCurrencyResponse{Amount: amount}
	return &resp
}

func initCurrencies(ctx context.Context, req *boutique.InitCurrencyRequest) *boutique.InitCurrencyResponse {
	boutique.InitCurrencies(ctx, req.Currencies)
	resp := boutique.InitCurrencyResponse{Ok: "OK"}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/set_currency", wrappers.NonROWrapper[boutique.SetCurrencySupportRequest, boutique.SetCurrencySupportResponse](setCurrency))
	http.HandleFunc("/init_currencies", wrappers.NonROWrapper[boutique.InitCurrencyRequest, boutique.InitCurrencyResponse](initCurrencies))
	http.HandleFunc("/ro_get_currencies", wrappers.ROWrapper[boutique.GetSupportedCurrenciesRequest, boutique.GetSupportedCurrenciesResponse](getCurrencies))
	http.HandleFunc("/ro_convert_currency", wrappers.ROWrapper[boutique.ConvertCurrencyRequest, boutique.ConvertCurrencyResponse](convertCurrency))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
