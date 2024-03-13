package hotel

import (
	"context"
	"github.com/eniac/mucache/pkg/invoke"
)

func SearchHotels(ctx context.Context, inDate string, outDate string, location string) []HotelProfile {
	req1 := NearbyRequest{InDate: inDate, OutDate: outDate, Location: location}
	hotelIdsRes := invoke.Invoke[NearbyResponse](ctx, "search", "ro_nearby", req1)
	rates := hotelIdsRes.Rates

	hotelIds := make([]string, len(rates))
	for i, rate := range rates {
		hotelIds[i] = rate.HotelId
	}

	req2 := CheckAvailabilityRequest{
		CustomerName: "",
		HotelIds:     hotelIds,
		InDate:       inDate,
		OutDate:      outDate,
		RoomNumber:   1,
	}
	availableHotelIdsRes := invoke.Invoke[CheckAvailabilityResponse](ctx, "reservation", "ro_check_availability", req2)

	//fmt.Printf("[Frontend] Location: %v -- Available hotel ids: %v\n", location, availableHotelIdsRes)

	req3 := GetProfilesRequest{HotelIds: availableHotelIdsRes.HotelIds}
	profilesRes := invoke.Invoke[GetProfilesResponse](ctx, "profile", "ro_get_profiles", req3)
	//fmt.Printf("[Frontend] Reviews read: %v\n", reviewsRes)
	return profilesRes.Profiles
}

func StoreHotel(ctx context.Context, hotelId string, name string, phone string, location string, rate int, capacity int, info string) string {
	req1 := StoreHotelLocationRequest{Location: location, HotelId: hotelId}
	invoke.Invoke[StoreHotelLocationResponse](ctx, "search", "store_hotel_location", req1)

	req2 := StoreRateRequest{Rate: Rate{HotelId: hotelId, Price: rate}}
	invoke.Invoke[StoreRateResponse](ctx, "rate", "store_rate", req2)

	req3 := AddHotelAvailabilityRequest{
		HotelId:  hotelId,
		Capacity: capacity,
	}
	invoke.Invoke[AddHotelAvailabilityResponse](ctx, "reservation", "add_hotel_availability", req3)

	hotelProfile := HotelProfile{
		HotelId: hotelId,
		Name:    name,
		Phone:   phone,
		Info:    info,
	}
	req4 := StoreProfileRequest{Profile: hotelProfile}
	invoke.Invoke[StoreProfileRequest](ctx, "profile", "store_profile", req4)
	return hotelId
}

func FrontendReservation(ctx context.Context, hotelId string, inDate string, outDate string, rooms int, username string, password string) bool {
	req1 := LoginRequest{
		Username: username,
		Password: password,
	}
	tokenRes := invoke.Invoke[LoginResponse](ctx, "user", "login", req1)
	if tokenRes.Token != "OK" {
		return false
	}

	req2 := MakeReservationRequest{
		CustomerName: username,
		HotelId:      hotelId,
		InDate:       inDate,
		OutDate:      outDate,
		RoomNumber:   rooms,
	}
	successRes := invoke.Invoke[MakeReservationResponse](ctx, "reservation", "make_reservation", req2)
	return successRes.Success
}
