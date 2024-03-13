package hotel

import (
	"context"
	"github.com/eniac/mucache/pkg/invoke"
	"github.com/eniac/mucache/pkg/state"
)

func Nearby(ctx context.Context, inDate string, outDate string, location string) []Rate {
	// Find the hotel ids in that location
	hotelIds := getHotelIdsForLocation(ctx, location)

	// Get the rates for these hotels
	req := GetRatesRequest{HotelIds: hotelIds}
	ratesRes := invoke.Invoke[GetRatesResponse](ctx, "rate", "ro_get_rates", req)
	return ratesRes.Rates
}

func StoreHotelLocation(ctx context.Context, hotelId string, location string) string {
	hotelIds := getHotelIdsForLocation(ctx, location)
	// Keep saved reviews bounded to 10 for consistent performance measurements
	if len(hotelIds) >= 10 {
		hotelIds = hotelIds[1:]
	}
	hotelIds = append(hotelIds, hotelId)
	state.SetState(ctx, location, hotelIds)
	return hotelId
}

func getHotelIdsForLocation(ctx context.Context, location string) []string {
	hotelIds, err := state.GetState[[]string](ctx, location)
	// If err != nil then the key does not exist
	if err != nil {
		return []string{}
	} else {
		return hotelIds
	}
}
