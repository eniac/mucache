package movie

import (
	"context"
	"github.com/eniac/mucache/pkg/state"
)

func StoreMovieInfo(ctx context.Context, movieId string, info string, castIds []string, plotId string) string {
	movieInfo := MovieInfo{
		MovieId: movieId,
		Info:    info,
		CastIds: castIds,
		PlotId:  plotId,
	}
	state.SetState(ctx, movieId, movieInfo)
	return movieId
}

func ReadMovieInfo(ctx context.Context, movieId string) MovieInfo {
	movieInfo, err := state.GetState[MovieInfo](ctx, movieId)
	if err != nil {
		panic(err)
	}
	return movieInfo
}
