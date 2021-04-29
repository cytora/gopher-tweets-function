package main

import (
	"os"

	"github.com/google/martian/log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cytora/gopher-tweets-function/auth"
	"github.com/cytora/gopher-tweets-function/handlers"
	"github.com/davyzhang/agw"
	twitterLogin "github.com/dghubble/gologin/v2/twitter"
	"github.com/gorilla/mux"
)

func main() {
	frontendRedirect, ok := os.LookupEnv("FRONTEND_REDIRECT")
	if !ok {
		log.Errorf("frontend redirect not configured")
		return
	}
	sessionKey, ok := os.LookupEnv("SESSION_SECRET_KEY")
	if !ok {
		log.Errorf("session key not configured")
		return
	}

	oauthConfig, err := auth.NewAuthorizationConfig()
	if err != nil {
		log.Errorf("failed to create authorization config:%s", err.Error())
		return
	}
	loginHandler := handlers.NewLoginHandler(frontendRedirect, sessionKey)
	tweetHandler := handlers.NewTweetHandler(sessionKey)

	mux := mux.NewRouter()
	mux.HandleFunc(handlers.HealthCheckEndpoint, handlers.Health)
	mux.Handle(handlers.TwitterAuthEndpoint, twitterLogin.LoginHandler(oauthConfig, nil))
	mux.Handle(handlers.TwitterCallbackEndpoint, twitterLogin.CallbackHandler(oauthConfig, loginHandler.Login(), nil))
	mux.HandleFunc(handlers.TweetEndpoint, tweetHandler.Tweet)

	lambda.Start(agw.Handler(mux))
}
