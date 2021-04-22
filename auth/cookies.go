package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/cytora/gopher-tweets-function/model"

	"github.com/dghubble/sessions"
)

type CookieManager struct {
	store *sessions.CookieStore
}

const (
	sessionCredentialsName = "gopher-tweets-credentials"
	sessionLoginName       = "gopher-tweets-login"
	sessionAccessKey       = "accessKey"
	sessionAccessSecret    = "accessSecret"
	sessionScreenName      = "screenName"
)

func NewCookieManager() (*CookieManager, error) {
	sessionSecret, ok := os.LookupEnv("SESSION_SECRET_KEY")
	if !ok {
		return nil, errors.New("session secret key not found")
	}

	// sessionStore encodes and decodes session data stored in signed cookies
	store := sessions.NewCookieStore([]byte(sessionSecret), nil)
	return &CookieManager{
		store: store,
	}, nil
}

func (c *CookieManager) SaveCredentials(w http.ResponseWriter, screenName, accessToken, accessSecret string) error {
	session := c.store.New(sessionCredentialsName)
	session.Values[sessionAccessKey] = accessToken
	session.Values[sessionAccessSecret] = accessSecret
	session.Values[sessionScreenName] = screenName
	return session.Save(w)
}

func (c *CookieManager) GetAuthenticatedUser(req *http.Request) (*model.AuthenticatedUser, error) {
	session, err := c.store.Get(req, sessionCredentialsName)
	if err != nil {
		return nil, fmt.Errorf("missing session store:%s", err.Error())
	}
	access := session.Values[sessionAccessKey]
	secret := session.Values[sessionAccessSecret]
	if access == nil || secret == nil {
		return nil, errors.New("access key or secret access key not set")
	}
	return &model.AuthenticatedUser{
		AccessToken:  access.(string),
		AccessSecret: secret.(string),
	}, nil
}

func (c *CookieManager) GetUserProfile(req *http.Request) (*model.UserProfile, error) {
	session, err := c.store.Get(req, sessionCredentialsName)
	if err != nil {
		return nil, fmt.Errorf("missing session store:%s", err.Error())
	}
	screenName := session.Values[sessionScreenName]
	if screenName == nil {
		return &model.UserProfile{IsLoggedIn: false}, nil
	}
	return &model.UserProfile{
		IsLoggedIn: true,
		ScreenName: screenName.(string),
	}, nil
}

func SaveLoginCookie(w http.ResponseWriter, screenName string) {
	loginCookie := &http.Cookie{
		Name:     sessionLoginName,
		Value:    screenName,
		Path:     "/",
		HttpOnly: false,
		Secure:   false,
	}
	w.Header().Set("Set-Cookie", loginCookie.String())
}
