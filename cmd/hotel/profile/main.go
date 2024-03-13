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

func storeProfile(ctx context.Context, req *hotel.StoreProfileRequest) *hotel.StoreProfileResponse {
	hotelId := hotel.StoreProfile(ctx, req.Profile)
	//fmt.Println("Movie info stored for id: " + movieId)
	resp := hotel.StoreProfileResponse{HotelId: hotelId}
	return &resp
}

func getProfiles(ctx context.Context, req *hotel.GetProfilesRequest) *hotel.GetProfilesResponse {
	hotels := hotel.GetProfiles(ctx, req.HotelIds)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := hotel.GetProfilesResponse{Profiles: hotels}
	//fmt.Printf("[ReviewStorage] Response: %v\n", resp)
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/store_profile", wrappers.NonROWrapper[hotel.StoreProfileRequest, hotel.StoreProfileResponse](storeProfile))
	http.HandleFunc("/ro_get_profiles", wrappers.ROWrapper[hotel.GetProfilesRequest, hotel.GetProfilesResponse](getProfiles))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
