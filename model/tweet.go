package model

type PostTweetRequest struct {
	Body string `json:"body"`
}

type PostTweetResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}
