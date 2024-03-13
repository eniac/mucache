package movie

import (
	"context"
	"github.com/eniac/mucache/pkg/state"
)

func StoreCastInfo(ctx context.Context, castId string, name string, info string) string {
	castInfo := CastInfo{
		CastId: castId,
		Name:   name,
		Info:   info,
	}
	state.SetState(ctx, castId, castInfo)
	return castId
}

func ReadCastInfos(ctx context.Context, castIds []string) []CastInfo {
	//fmt.Printf("Keys: %+v\n", castIds)
	//castInfos := make([]CastInfo, len(castIds))
	//for i, castId := range castIds {
	//	castInfo, err := state.GetState[CastInfo](ctx, castId)
	//	if err != nil {
	//		// If we don't find the cast info, we can simply return an empty struct for that one
	//		// This might sometimes happen because we haven't populated the actor dictionary with all actors.
	//		castInfos[i] = CastInfo{}
	//	} else {
	//		castInfos[i] = castInfo
	//	}
	//}

	// Bulk
	var castInfos []CastInfo
	if len(castIds) > 0 {
		castInfos = state.GetBulkStateDefault[CastInfo](ctx, castIds, CastInfo{})
	} else {
		castInfos = make([]CastInfo, len(castIds))
	}
	return castInfos
}
