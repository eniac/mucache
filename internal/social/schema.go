package social

type Media struct {
	MediaId   string `json:"media_id"`
	MediaType string `json:"media_type"`
}

type User struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	//FirstName string `json:"first_name"`
	//LastName  string `json:"last_name"`
	//Password  string `json:"password"`
	//Salt      string `json:"salt"`
}

type UserMention struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type Url struct {
	ShortenedUrl string `json:"shortened_url"`
	ExpandedUrl  string `json:"expanded_url"`
}

type Post struct {
	PostId    string `json:"post_id,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Text      string `json:"text"`
	CreatorId string `json:"creator_id"`
	//UserMentions []UserMention `json:"user_mentions"`
	//Media        []Media       `json:"media"`
	//Urls         []Url         `json:"urls"`
}

type SGVertex struct {
	UserId    string   `json:"userId"`
	Followers []string `json:"followers"`
	Followees []string `json:"followed"`
}

// post_storage

type ReadPostRequest struct {
	PostId string `json:"post_id"`
}

type ReadPostResponse struct {
	Post Post `json:"post"`
}

type ReadPostsRequest struct {
	PostIds []string `json:"post_ids"`
}

type ReadPostsResponse struct {
	Posts []Post `json:"posts"`
}

type StorePostRequest struct {
	CreatorId string `json:"creator_id"`
	Text      string `json:"text"`
}

type StorePostResponse struct {
	PostId string `json:"post_id"`
}

type StorePostMultiRequest struct {
	CreatorId string `json:"creator_id"`
	Text      string `json:"text"`
	Number    int    `json:"number"`
}

type StorePostMultiResponse struct {
	PostIds []string `json:"post_ids"`
}

// home_timeline

type ReadHomeTimelineRequest struct {
	UserId string `json:"user_id"`
}

type ReadHomeTimelineResponse struct {
	Posts []Post `json:"posts"`
}

type WriteHomeTimelineRequest struct {
	UserId  string   `json:"user_id"`
	PostIds []string `json:"post_ids"`
}

// user_timeline

type ReadUserTimelineRequest struct {
	UserId string `json:"user_id"`
}

type ReadUserTimelineResponse struct {
	Posts []Post `json:"posts"`
}

type WriteUserTimelineRequest struct {
	UserId  string   `json:"user_id"`
	PostIds []string `json:"post_ids"`
}

// social_graph

type InsertUserRequest struct {
	UserId string `json:"user_id"`
}

type GetFollowersRequest struct {
	UserId string `json:"user_id"`
}

type GetFollowersResponse struct {
	Followers []string `json:"followers"`
}

type GetFolloweesRequest struct {
	UserId string `json:"user_id"`
}

type GetFolloweesResponse struct {
	Followees []string `json:"followees"`
}

type FollowRequest struct {
	FollowerId string `json:"follower_id"`
	FolloweeId string `json:"followee_id"`
}

type FollowManyRequest struct {
	UserId      string   `json:"user_id"`
	FollowerIds []string `json:"follower_ids"`
	FolloweeIds []string `json:"followee_ids"`
}

// compose_post

type ComposePostRequest struct {
	CreatorId string `json:"creator_id"`
	Text      string `json:"text"`
}

type ComposePostMultiRequest struct {
	CreatorId string `json:"creator_id"`
	Text      string `json:"text"`
	Number    int    `json:"number"`
}
