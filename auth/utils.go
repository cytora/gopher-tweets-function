package auth

import (
	"errors"
	"fmt"
	"os"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
)

type Config struct {
	FrontendRedirect string
	SessionKey       string
	Authorization    *oauth1.Config
}

func NewConfig() (*Config, error) {
	frontendRedirect, ok := os.LookupEnv("FRONTEND_REDIRECT")
	if !ok {
		return nil, errors.New("frontend redirect not configured")
	}

	sessionKey, ok := os.LookupEnv("SESSION_SECRET_KEY")
	if !ok {
		return nil, errors.New("session key not configured")
	}

	oauthConfig, err := newAuthorizationConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create authorization config:%s", err.Error())
	}

	return &Config{
		FrontendRedirect: frontendRedirect,
		SessionKey:       sessionKey,
		Authorization:    oauthConfig,
	}, nil
}
func newAuthorizationConfig() (*oauth1.Config, error) {
	twitterConsumerKey, ok := os.LookupEnv("TWITTER_CONSUMER_KEY")
	if !ok {
		return nil, errors.New("twitter consumer key not found")
	}
	twitterConsumerSecret, ok := os.LookupEnv("TWITTER_CONSUMER_SECRET")
	if !ok {
		return nil, errors.New("twitter consumer secret not found")
	}
	twitterLoginCallback, ok := os.LookupEnv("TWITTER_LOGIN_CALLBACK")
	if !ok {
		return nil, errors.New("twitter login callback not found")
	}

	return &oauth1.Config{
		ConsumerKey:    twitterConsumerKey,
		ConsumerSecret: twitterConsumerSecret,
		CallbackURL:    twitterLoginCallback,
		Endpoint:       twitter.AuthorizeEndpoint,
	}, nil
}

func NewTweetAuthConfig() (*oauth1.Config, error) {
	twitterConsumerKey, ok := os.LookupEnv("TWITTER_CONSUMER_KEY")
	if !ok {
		return nil, errors.New("twitter consumer key not found")
	}
	twitterConsumerSecret, ok := os.LookupEnv("TWITTER_CONSUMER_SECRET")
	if !ok {
		return nil, errors.New("twitter consumer secret not found")
	}
	return oauth1.NewConfig(twitterConsumerKey, twitterConsumerSecret), nil
}
