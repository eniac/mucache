package movie

import (
	"context"
	"github.com/eniac/mucache/pkg/state"
)

func StoreReview(ctx context.Context, review Review) string {
	state.SetState(ctx, review.ReviewId, review)
	return review.ReviewId
}

func ReadReviews(ctx context.Context, reviewIds []string) []Review {
	//fmt.Printf("[ReviewStorage] Asked for: %v\n", reviewIds)
	//reviews := make([]Review, len(reviewIds))
	//for i, reviewId := range reviewIds {
	//	review, err := state.GetState[Review](ctx, reviewId)
	//	if err != nil {
	//		panic(err)
	//	}
	//	reviews[i] = review
	//}
	//fmt.Printf("[ReviewStorage] Returning: %v\n", reviews)
	// Bulk
	var reviews []Review
	if len(reviewIds) > 0 {
		reviews = state.GetBulkStateDefault[Review](ctx, reviewIds, Review{})
	} else {
		reviews = make([]Review, len(reviews))
	}
	return reviews
}
