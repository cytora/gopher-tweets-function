package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cytora/gopher-tweets-function/auth"
	"github.com/cytora/gopher-tweets-function/model"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const TweetEndpoint = "/tweets"
const TokenHeader = "X-Gopher-Token"

type TweetHandler struct {
	sessionKey string
}

func NewTweetHandler(sessionKey string) *TweetHandler {
	return &TweetHandler{
		sessionKey: sessionKey,
	}
}

// Tweet posts a tweet after the user has already been authenticated
func (h *TweetHandler) Tweet(w http.ResponseWriter, req *http.Request) {
	//These are preflight requests
	if req.Method == "OPTIONS" {
		r := model.Response{
			StatusCode: http.StatusOK,
			Body: model.PostTweetResponse{
				Message: "preflight",
			},
			Error: nil,
		}
		r.Write(w)
		return
	}

	var post model.PostTweetRequest
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	if err := json.NewDecoder(req.Body).Decode(&post); err != nil {
		r := model.Response{
			StatusCode: http.StatusBadRequest,
			Body:       nil,
			Error:      fmt.Errorf("error decoding body:%s", err.Error()),
		}
		r.Write(w)
		return
	}

	clt, err := h.getTwitterClient(req)
	if err != nil {
		r := model.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       nil,
			Error:      fmt.Errorf("error getting Twitter client:%s", err.Error()),
		}
		r.Write(w)
		return
	}
	t, _, err := clt.Statuses.Update(post.Tweet, nil)
	if err != nil {
		r := model.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       nil,
			Error:      err,
		}
		r.Write(w)
		return
	}

	r := model.Response{
		StatusCode: http.StatusOK,
		Body: model.PostTweetResponse{
			ID:      t.ID,
			Message: "tweet published! hooray!",
		},
		Error: nil,
	}
	r.Write(w)
	return
}

func (h *TweetHandler) getTwitterClient(req *http.Request) (*twitter.Client, error) {
	config, err := auth.NewTweetAuthConfig()
	if err != nil {
		return nil, err
	}
	internalToken := req.Header.Get(TokenHeader)
	creds, err := auth.Decode(h.sessionKey, internalToken)
	if err != nil {
		return nil, err
	}
	token := oauth1.NewToken(creds.ExtAccessKey, creds.ExtAccessSecret)
	httpClient := config.Client(req.Context(), token)
	return twitter.NewClient(httpClient), nil
}
