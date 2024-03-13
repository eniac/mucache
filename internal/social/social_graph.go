package social

import (
	"context"
	"github.com/eniac/mucache/pkg/state"
)

func GetFollowers(ctx context.Context, userId string) []string {
	sg, err := state.GetState[SGVertex](ctx, userId)
	if err != nil {
		panic(err)
	}
	return sg.Followers
}

func GetFollowees(ctx context.Context, userId string) []string {
	sg, err := state.GetState[SGVertex](ctx, userId)
	if err != nil {
		panic(err)
	}
	return sg.Followees
}

func Follow(ctx context.Context, followerId string, followeeId string) {
	sg, err := state.GetState[SGVertex](ctx, followerId)
	if err != nil {
		sg = SGVertex{
			UserId:    followerId,
			Followers: []string{},
			Followees: []string{},
		}
	}
	sg.Followees = append(sg.Followees, followeeId)
	state.SetState(ctx, followerId, sg)

	sg, err = state.GetState[SGVertex](ctx, followeeId)
	if err != nil {
		sg = SGVertex{
			UserId:    followeeId,
			Followers: []string{},
			Followees: []string{},
		}
	}
	if err != nil {
		panic(err)
	}
	sg.Followers = append(sg.Followers, followerId)
	state.SetState(ctx, followeeId, sg)
}

// Only used for populating
func FollowMulti(ctx context.Context, userId string, followerIds []string, followeeIds []string) {
	sg := SGVertex{
		UserId:    userId,
		Followers: followerIds,
		Followees: followeeIds,
	}
	if len(sg.Followers) >= 10 {
		sg.Followers = sg.Followers[:10]
	}
	if len(sg.Followees) >= 10 {
		sg.Followees = sg.Followees[:10]
	}
	state.SetState(ctx, userId, sg)
}

func InsertUser(ctx context.Context, userId string) {
	sg := SGVertex{
		Followers: []string{},
		Followees: []string{},
	}
	state.SetState(ctx, userId, sg)
}
