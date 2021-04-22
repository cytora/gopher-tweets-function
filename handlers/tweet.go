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
const ProfileEndpoint = "/twitter/profile"

type TweetHandler struct {
	cookieManager *auth.CookieManager
}

func NewTweetHandler(cookieManager *auth.CookieManager) *TweetHandler {
	return &TweetHandler{
		cookieManager: cookieManager,
	}
}

// Tweet posts a tweet after the user has already been authenticated
func (h *TweetHandler) Tweet(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Tweet REQUEST: %+v\n", req)
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
	var post model.PostTweetRequest
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	if err := json.NewDecoder(req.Body).Decode(&req); err != nil {
		r := model.Response{
			StatusCode: http.StatusBadRequest,
			Body:       nil,
			Error:      fmt.Errorf("error getting decoding body:%s", err.Error()),
		}
		r.Write(w)
		return
	}
	t, _, err := clt.Statuses.Update(post.Body, nil)
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

	authUser, err := h.cookieManager.GetAuthenticatedUser(req)
	if err != nil {
		return nil, err
	}
	token := oauth1.NewToken(authUser.AccessToken, authUser.AccessSecret)
	httpClient := config.Client(req.Context(), token)
	return twitter.NewClient(httpClient), nil
}

// Profile shows a personal profile - FOR BE DEBUG PURPOSES
func (h *TweetHandler) Profile(w http.ResponseWriter, req *http.Request) {
	userProfile, err := h.cookieManager.GetUserProfile(req)
	if err != nil {
		r := model.Response{
			StatusCode: http.StatusNotFound,
			Body:       nil,
			Error:      err,
		}
		r.Write(w)
		return
	}

	// authenticated profile
	r := model.Response{
		StatusCode: http.StatusOK,
		Body:       userProfile,
		Error:      nil,
	}
	r.Write(w)
}
