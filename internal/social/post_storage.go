package social

import (
	"context"
	"github.com/eniac/mucache/pkg/state"
	"github.com/lithammer/shortuuid"
	"time"
)

// StorePost return the post id
func StorePost(ctx context.Context, creatorId string, text string) string {
	postIds := StorePostMulti(ctx, creatorId, text, 1)
	//fmt.Printf("[StorePost] Stored: %+v, %+v, %+v\nReturned: %+v of size %+v\n", creatorId, text, 1, postIds, len(postIds))
	//fmt.Printf("[StorePost] Returned: %+v\n", postIds[0])
	return postIds[0]
}

func StorePostMulti(ctx context.Context, creatorId string, text string, number int) []string {
	posts := make(map[string]interface{}, number)
	postIds := make([]string, number)
	for i := 0; i < number; i++ {
		postId := shortuuid.New()
		timestamp := time.Now().Unix()
		posts[postId] = Post{
			PostId:    postId,
			CreatorId: creatorId,
			Text:      text,
			Timestamp: timestamp,
		}
		postIds[i] = postId
	}
	state.SetBulkState(ctx, posts)
	//fmt.Printf("[StorePostMulti] Returning %+v\n", postIds)
	return postIds
}

func ReadPost(ctx context.Context, postId string) Post {
	post, err := state.GetState[Post](ctx, postId)
	if err != nil {
		panic(err)
	}
	return post
}

func ReadPosts(ctx context.Context, postIds []string) []Post {
	//retPosts := make([]Post, 0)
	//for _, postId := range postIds {
	//	post, err := state.GetState[Post](ctx, postId)
	//	if err != nil {
	//		panic(err)
	//	}
	//	retPosts = append(retPosts, post)
	//}
	retPosts, err := state.GetBulkState[Post](ctx, postIds)
	if err != nil {
		panic(err)
	}
	return retPosts
}
