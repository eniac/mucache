package movie

import (
	"context"
	"github.com/eniac/mucache/pkg/invoke"
	"github.com/lithammer/shortuuid"
	"time"
)

func Compose(ctx context.Context, username string, password string, title string, rating int, text string) bool {
	req1 := LoginRequest{
		Username: username,
		Password: password,
	}
	//fmt.Printf("[Page] Movie id asked: %v\n", req)
	tokenRes := invoke.Invoke[LoginResponse](ctx, "user", "login", req1)
	if tokenRes.Token != "OK" {
		return false
	}
	reqId := shortuuid.New()

	// TODO: Make them async
	req2 := GetUniqueIdRequest{ReqId: reqId}
	reviewIdRes := invoke.Invoke[GetUniqueIdResponse](ctx, "uniqueid", "get_unique_id", req2)
	req3 := GetUserIdRequest{Username: username}
	userIdRes := invoke.Invoke[GetUserIdResponse](ctx, "user", "ro_get_user_id", req3)
	req4 := GetMovieIdRequest{Title: title}
	movieIdRes := invoke.Invoke[GetMovieIdResponse](ctx, "movieid", "ro_get_movie_id", req4)
	//fmt.Printf("[Frontend] Title: %v was tied to id: %v\n", title, movieIdRes)
	ts := time.Now().Unix()
	review := Review{
		ReviewId:  reviewIdRes.ReviewId,
		UserId:    userIdRes.UserId,
		ReqId:     reqId,
		Text:      text,
		MovieId:   movieIdRes.MovieId,
		Rating:    rating,
		Timestamp: ts,
	}
	// This is the only sync call (all the previous ones can be async
	req5 := ComposeReviewRequest{Review: review}
	invoke.Invoke[ComposeReviewResponse](ctx, "composereview", "compose_review", req5)
	return true
}
