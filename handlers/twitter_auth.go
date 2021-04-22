package handlers

import (
	"errors"
	"net/http"
	"os"

	"github.com/cytora/gopher-tweets-function/auth"
	"github.com/cytora/gopher-tweets-function/model"
	oauth1Login "github.com/dghubble/gologin/v2/oauth1"
	twitterLogin "github.com/dghubble/gologin/v2/twitter"
)

const (
	TwitterAuthEndpoint     = "/twitter/login"
	TwitterCallbackEndpoint = "/twitter/callback"
)

type LoginHandler struct {
	cookieManager *auth.CookieManager
	redirect      string
}

func NewLoginHandler(cookieManager *auth.CookieManager) (*LoginHandler, error) {
	frontendRedirect, ok := os.LookupEnv("FRONTEND_REDIRECT")
	if !ok {
		return nil, errors.New("frontend redirect not configured")
	}
	return &LoginHandler{
		cookieManager: cookieManager,
		redirect:      frontendRedirect,
	}, nil
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

		if err := h.cookieManager.SaveCredentials(w, twitterUser.ScreenName, accessToken, accessSecret); err != nil {
			r := model.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       nil,
				Error:      err,
			}
			r.Write(w)
			return
		}

		auth.SaveLoginCookie(w, twitterUser.ScreenName)
		http.Redirect(w, req, h.redirect, http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
