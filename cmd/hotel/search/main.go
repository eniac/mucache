package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/hotel"
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

func nearby(ctx context.Context, req *hotel.NearbyRequest) *hotel.NearbyResponse {
	rates := hotel.Nearby(ctx, req.InDate, req.OutDate, req.Location)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := hotel.NearbyResponse{Rates: rates}
	//fmt.Printf("[ReviewStorage] Response: %v\n", resp)
	return &resp
}

func storeHotelLocation(ctx context.Context, req *hotel.StoreHotelLocationRequest) *hotel.StoreHotelLocationResponse {
	hotelId := hotel.StoreHotelLocation(ctx, req.HotelId, req.Location)
	resp := hotel.StoreHotelLocationResponse{HotelId: hotelId}
	//fmt.Printf("[ReviewStorage] Response: %v\n", resp)
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_nearby", wrappers.ROWrapper[hotel.NearbyRequest, hotel.NearbyResponse](nearby))
	http.HandleFunc("/store_hotel_location", wrappers.NonROWrapper[hotel.StoreHotelLocationRequest, hotel.StoreHotelLocationResponse](storeHotelLocation))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
