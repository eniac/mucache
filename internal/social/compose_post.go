package social

import (
	"context"
	"github.com/eniac/mucache/pkg/invoke"
)

func ComposePost(ctx context.Context, text string, creatorId string) {
	req1 := StorePostRequest{
		CreatorId: creatorId,
		Text:      text,
	}
	resp1 := invoke.Invoke[StorePostResponse](ctx, "poststorage", "store_post", req1)
	postId := resp1.PostId
	//fmt.Printf("Stored: %+v\nReturned: %+v\n", req1, resp1)
	req2 := WriteUserTimelineRequest{
		UserId:  creatorId,
		PostIds: []string{postId},
	}
	invoke.Invoke[string](ctx, "usertimeline", "write_user_timeline", req2)
	//fmt.Printf("Stored: %+v\n", req2)
	req3 := WriteHomeTimelineRequest{
		UserId:  creatorId,
		PostIds: []string{postId},
	}
	invoke.Invoke[string](ctx, "hometimeline", "write_home_timeline", req3)
	//fmt.Printf("Stored: %+v\n", req3)
}

func ComposeMulti(ctx context.Context, text string, number int, creatorId string) {
	req1 := StorePostMultiRequest{
		CreatorId: creatorId,
		Text:      text,
		Number:    number,
	}
	resp1 := invoke.Invoke[StorePostMultiResponse](ctx, "poststorage", "store_post_multi", req1)
	//fmt.Printf("Stored: %+v\nReturned: %+v\n", req1, resp1)
	postIds := resp1.PostIds
	req2 := WriteUserTimelineRequest{
		UserId:  creatorId,
		PostIds: postIds,
	}
	invoke.Invoke[string](ctx, "usertimeline", "write_user_timeline", req2)
	//fmt.Printf("Stored: %+v\n", req2)
	req3 := WriteHomeTimelineRequest{
		UserId:  creatorId,
		PostIds: postIds,
	}
	invoke.Invoke[string](ctx, "hometimeline", "write_home_timeline", req3)
	//fmt.Printf("Stored: %+v\n", req3)
}
