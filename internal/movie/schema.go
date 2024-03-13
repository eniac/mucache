package movie

type MovieInfo struct {
	MovieId string   `json:"movie_id"`
	Info    string   `json:"info"`
	CastIds []string `json:"cast_ids"`
	PlotId  string   `json:"plot_id"`
}

type CastInfo struct {
	CastId string `json:"cast_id"`
	Name   string `json:"name"`
	Info   string `json:"info"`
}

type User struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Password []byte `json:"password"`
	Salt     string `json:"salt"`
}

type Review struct {
	ReviewId  string `json:"review_id"`
	UserId    string `json:"user_id"`
	ReqId     string `json:"req_id"`
	Text      string `json:"text"`
	MovieId   string `json:"movie_id"`
	Rating    int    `json:"rating"`
	Timestamp int64  `json:"timestamp"`
}

type MovieReview struct {
	ReviewId  string `json:"review_id"`
	Timestamp int64  `json:"timestamp"`
}

type UserReview struct {
	ReviewId  string `json:"review_id"`
	Timestamp int64  `json:"timestamp"`
}

type Page struct {
	MovieInfo MovieInfo  `json:"movie_info"`
	Reviews   []Review   `json:"reviews"`
	CastInfos []CastInfo `json:"cast_infos"`
	Plot      string     `json:"plot"`
}

// movie_info

type ReadMovieInfoRequest struct {
	MovieId string `json:"movie_id"`
}

type ReadMovieInfoResponse struct {
	Info MovieInfo `json:"movie_info"`
}

type StoreMovieInfoRequest struct {
	MovieId string   `json:"movie_id"`
	Info    string   `json:"movie_info"`
	CastIds []string `json:"cast_ids"`
	PlotId  string   `json:"plot_id"`
}

type StoreMovieInfoResponse struct {
	MovieId string `json:"movie_id"`
}

// cast_info

type ReadCastInfosRequest struct {
	CastIds []string `json:"cast_ids"`
}

type ReadCastInfosResponse struct {
	Infos []CastInfo `json:"cast_infos"`
}

type StoreCastInfoRequest struct {
	CastId string `json:"cast_id"`
	Name   string `json:"name"`
	Info   string `json:"info"`
}

type StoreCastInfoResponse struct {
	CastId string `json:"cast_id"`
}

// page
type ReadPageRequest struct {
	MovieId string `json:"movie_id"`
}

type ReadPageResponse struct {
	Page Page `json:"page"`
}

// user
type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	Ok bool `json:"ok"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type GetUserIdRequest struct {
	Username string `json:"username"`
}

type GetUserIdResponse struct {
	UserId string `json:"user_id"`
}

// review_storage
type ReadReviewsRequest struct {
	ReviewIds []string `json:"review_ids"`
}

type ReadReviewsResponse struct {
	Reviews []Review `json:"reviews"`
}

type StoreReviewRequest struct {
	Review Review `json:"review"`
}

type StoreReviewResponse struct {
	ReviewId string `json:"review_id"`
}

// plot
type ReadPlotRequest struct {
	PlotId string `json:"plot_id"`
}

type ReadPlotResponse struct {
	Plot string `json:"plot"`
}

type WritePlotRequest struct {
	PlotId string `json:"plot_id"`
	Plot   string `json:"plot"`
}

type WritePlotResponse struct {
	PlotId string `json:"plot_id"`
}

// movie_review
type ReadMovieReviewsRequest struct {
	MovieId string `json:"movie_id"`
}

type ReadMovieReviewsResponse struct {
	Reviews []Review `json:"reviews"`
}

type UploadMovieReviewRequest struct {
	MovieId   string `json:"movie_id"`
	ReviewId  string `json:"review_id"`
	Timestamp int64  `json:"timestamp"`
}

type UploadMovieReviewResponse struct {
	ReviewId string `json:"review_id"`
}

// user_reviews
type ReadUserReviewsRequest struct {
	UserId string `json:"user_id"`
}

type ReadUserReviewsResponse struct {
	Reviews []Review `json:"reviews"`
}

type UploadUserReviewRequest struct {
	UserId    string `json:"user_id"`
	ReviewId  string `json:"review_id"`
	Timestamp int64  `json:"timestamp"`
}

type UploadUserReviewResponse struct {
	ReviewId string `json:"review_id"`
}

// frontend
type ComposeRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Title    string `json:"title"`
	Rating   int    `json:"rating"`
	Text     string `json:"text"`
}

type ComposeResponse struct {
	Ok bool `json:"ok"`
}

// compose review
type ComposeReviewRequest struct {
	Review Review `json:"review"`
}

type ComposeReviewResponse struct {
	Ok string `json:"ok"`
}

// unique_id
type GetUniqueIdRequest struct {
	ReqId string `json:"req_id"`
}

type GetUniqueIdResponse struct {
	ReviewId string `json:"review_id"`
}

// movie_id
type GetMovieIdRequest struct {
	Title string `json:"title"`
}

type GetMovieIdResponse struct {
	MovieId string `json:"movie_id"`
}

type RegisterMovieIdRequest struct {
	Title   string `json:"title"`
	MovieId string `json:"movie_id"`
}

type RegisterMovieIdResponse struct {
	Ok string `json:"ok"`
}
