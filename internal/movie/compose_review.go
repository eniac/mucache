package movie

import (
	"context"
	"github.com/eniac/mucache/pkg/invoke"
)

// Note: ComposeReview is rearchitected from its original Deathstar version to not do this complex pushing of intermediate data.
//
//	This has nothing to do with our caching infrastructure and it was just done for simplicity and performance here.
func ComposeReview(ctx context.Context, review Review) {
	// TODO: Make invocations req2, req3, async
	req1 := StoreReviewRequest{Review: review}
	invoke.Invoke[StoreReviewResponse](ctx, "reviewstorage", "store_review", req1)
	req2 := UploadMovieReviewRequest{
		MovieId:   review.MovieId,
		ReviewId:  review.ReviewId,
		Timestamp: review.Timestamp,
	}
	invoke.Invoke[UploadMovieReviewResponse](ctx, "moviereviews", "upload_movie_review", req2)
	req3 := UploadUserReviewRequest{
		UserId:    review.UserId,
		ReviewId:  review.ReviewId,
		Timestamp: review.Timestamp,
	}
	invoke.Invoke[UploadUserReviewResponse](ctx, "userreviews", "upload_user_review", req3)
	//fmt.Printf("[ComposeReview] Successfully stored review: %v\n", review)
}
