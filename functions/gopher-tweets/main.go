package main

import (
	"github.com/google/martian/log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cytora/gopher-tweets-function/auth"
	"github.com/cytora/gopher-tweets-function/handlers"
	"github.com/davyzhang/agw"
	twitterLogin "github.com/dghubble/gologin/v2/twitter"
	"github.com/gorilla/mux"
)

func main() {
	cookieManager, err := auth.NewCookieManager()
	if err != nil {
		log.Errorf("failed to create cookie manager:%s", err.Error())
		return
	}
	oauthConfig, err := auth.NewAuthorizationConfig()
	if err != nil {
		log.Errorf("failed to create authorization config:%s", err.Error())
		return
	}
	loginHandler, err := handlers.NewLoginHandler(cookieManager)
	if err != nil {
		log.Errorf("failed to create login handler:%s", err.Error())
		return
	}
	tweetHandler := handlers.NewTweetHandler(cookieManager)

	mux := mux.NewRouter()
	mux.HandleFunc(handlers.HealthCheckEndpoint, handlers.Health)
	mux.Handle(handlers.TwitterAuthEndpoint, twitterLogin.LoginHandler(oauthConfig, nil))
	mux.Handle(handlers.TwitterCallbackEndpoint, twitterLogin.CallbackHandler(oauthConfig, loginHandler.Login(), nil))
	mux.HandleFunc(handlers.TweetEndpoint, tweetHandler.Tweet)
	mux.HandleFunc(handlers.ProfileEndpoint, tweetHandler.Profile)

	lambda.Start(agw.Handler(mux))
}
