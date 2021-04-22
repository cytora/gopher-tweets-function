package auth

import (
	"errors"
	"os"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
)

func NewAuthorizationConfig() (*oauth1.Config, error) {
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
