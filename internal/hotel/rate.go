package hotel

import (
	"context"
	"github.com/eniac/mucache/pkg/state"
)

func StoreRate(ctx context.Context, rate Rate) string {
	state.SetState(ctx, rate.HotelId, rate)
	return rate.HotelId
}

func GetRates(ctx context.Context, hotelIds []string) []Rate {
	//fmt.Printf("[ReviewStorage] Asked for: %v\n", reviewIds)
	rates := make([]Rate, len(hotelIds))
	for i, hotelId := range hotelIds {
		rate, err := state.GetState[Rate](ctx, hotelId)
		if err != nil {
			panic(err)
		}
		rates[i] = rate
	}
	//fmt.Printf("[ReviewStorage] Returning: %v\n", reviews)
	return rates
}
