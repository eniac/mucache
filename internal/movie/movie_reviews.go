package movie

import (
	"context"
	"github.com/eniac/mucache/pkg/invoke"
	"github.com/eniac/mucache/pkg/state"
)

// This service simply keeps indexes of reviews for each movie

func UploadMovieReview(ctx context.Context, movieId string, reviewId string, timestamp int64) string {
	reviews := getMovieReviews(ctx, movieId)
	// Keep saved reviews bounded to 10 for consistent performance measurements
	if len(reviews) >= 10 {
		reviews = reviews[1:]
	}
	reviews = append(reviews, MovieReview{ReviewId: reviewId, Timestamp: timestamp})
	state.SetState(ctx, movieId, reviews)
	return movieId
}

func getMovieReviews(ctx context.Context, movieId string) []MovieReview {
	reviews, err := state.GetState[[]MovieReview](ctx, movieId)
	// If err != nil then the key does not exist
	if err != nil {
		return []MovieReview{}
	} else {
		return reviews
	}
}

func ReadMovieReviews(ctx context.Context, movieId string) []Review {
	movieReviews := getMovieReviews(ctx, movieId)
	reviewIds := make([]string, len(movieReviews))
	for i, movieReview := range movieReviews {
		reviewIds[i] = movieReview.ReviewId
	}
	req := ReadReviewsRequest{ReviewIds: reviewIds}
	//fmt.Printf("[MovieReviews] Asking review ids: %v\n", req)
	reviewStorageResp := invoke.Invoke[ReadReviewsResponse](ctx, "reviewstorage", "ro_read_reviews", req)
	reviews := reviewStorageResp.Reviews
	//fmt.Printf("[MovieReviews] Reviews read: %v\n", reviewStorageResp)
	return reviews
}
