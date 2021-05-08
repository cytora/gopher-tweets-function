package model

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TweetRequest struct {
	Tweet string `json:"tweet"`
}

type PostTweetResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

type AliveResponse struct {
	Alive   bool   `json:"alive"`
	Message string `json:"message"`
}

type TwitterAuthMessage struct {
	Message string `json:"message"`
}

func NewTweetRequest(req *http.Request) (*TweetRequest, error) {
	var post TweetRequest
	if err := json.NewDecoder(req.Body).Decode(&post); err != nil {
		return nil, fmt.Errorf("error during tweet decode:%s", err.Error())
	}
	return &post, nil
}
