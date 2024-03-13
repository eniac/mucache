package movie

import (
	"context"
	"github.com/eniac/mucache/pkg/invoke"
	"github.com/eniac/mucache/pkg/state"
)

// This service simply keeps indexes of reviews for each movie

func UploadUserReview(ctx context.Context, userId string, reviewId string, timestamp int64) string {
	reviews := getUserReviews(ctx, userId)
	// Keep saved reviews bounded to 10 for consistent performance measurements
	if len(reviews) >= 10 {
		reviews = reviews[1:]
	}
	newReviews := append(reviews, UserReview{ReviewId: reviewId, Timestamp: timestamp})
	state.SetState(ctx, userId, newReviews)
	return userId
}

func getUserReviews(ctx context.Context, userId string) []UserReview {
	reviews, err := state.GetState[[]UserReview](ctx, userId)
	// If err != nil then the key does not exist
	if err != nil {
		return []UserReview{}
	} else {
		return reviews
	}
}

func ReadUserReviews(ctx context.Context, userId string) []Review {
	userReviews := getUserReviews(ctx, userId)
	reviewIds := make([]string, len(userReviews))
	for i, userReview := range userReviews {
		reviewIds[i] = userReview.ReviewId
	}
	req := ReadReviewsRequest{ReviewIds: reviewIds}
	//fmt.Printf("[UserReviews] Asking review ids: %v\n", req)
	reviewStorageResp := invoke.Invoke[ReadReviewsResponse](ctx, "reviewstorage", "ro_read_reviews", req)
	reviews := reviewStorageResp.Reviews
	//fmt.Printf("[UserReviews] Reviews read: %v\n", reviewStorageResp)
	return reviews
}
