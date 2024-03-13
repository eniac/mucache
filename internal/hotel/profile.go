package hotel

import (
	"context"
	"github.com/eniac/mucache/pkg/state"
)

func StoreProfile(ctx context.Context, profile HotelProfile) string {
	state.SetState(ctx, profile.HotelId, profile)
	return profile.HotelId
}

func GetProfiles(ctx context.Context, hotelIds []string) []HotelProfile {
	//fmt.Printf("[ReviewStorage] Asked for: %v\n", reviewIds)
	//profiles := make([]HotelProfile, len(hotelIds))
	//for i, hotelId := range hotelIds {
	//	profile, err := state.GetState[HotelProfile](ctx, hotelId)
	//	if err != nil {
	//		panic(err)
	//	}
	//	profiles[i] = profile
	//}

	// Bulk
	var profiles []HotelProfile
	if len(hotelIds) > 0 {
		profiles = state.GetBulkStateDefault[HotelProfile](ctx, hotelIds, HotelProfile{})
	} else {
		profiles = make([]HotelProfile, len(hotelIds))
	}
	//fmt.Printf("[ReviewStorage] Returning: %v\n", reviews)
	return profiles
}
