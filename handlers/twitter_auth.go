package handlers

import (
	"fmt"
	"net/http"

	"github.com/cytora/gopher-tweets-function/auth"
	oauth1Login "github.com/dghubble/gologin/v2/oauth1"
	twitterLogin "github.com/dghubble/gologin/v2/twitter"
)

const (
	TwitterAuthEndpoint     = "/twitter/login"
	TwitterCallbackEndpoint = "/twitter/callback"
)

type LoginHandler struct {
	redirect   string
	sessionKey string
}

func NewLoginHandler(redirect, sessionKey string) *LoginHandler {
	return &LoginHandler{
		redirect:   redirect,
		sessionKey: sessionKey,
	}
}

// Login issues a cookie session after successful Twitter login
func (h *LoginHandler) Login() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		accessToken, accessSecret, err := oauth1Login.AccessTokenFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		twitterUser, err := twitterLogin.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		creds := &auth.Credentials{
			ExtAccessKey:    accessToken,
			ExtAccessSecret: accessSecret,
		}
		http.Redirect(w, req, fmt.Sprintf("%s/?token=%s&user=%s", h.redirect, creds.Encode(h.sessionKey), twitterUser.ScreenName), http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
