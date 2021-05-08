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
	cfg, err := auth.NewConfig()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	loginHandler := handlers.NewLoginHandler(cfg.FrontendRedirect, cfg.SessionKey)
	tweetHandler := handlers.NewTweetHandler(cfg.SessionKey)

	mux := mux.NewRouter()
	mux.HandleFunc(handlers.HealthCheckEndpoint, handlers.Health)
	mux.Handle(handlers.TwitterAuthEndpoint, twitterLogin.LoginHandler(cfg.Authorization, nil))
	mux.Handle(handlers.TwitterCallbackEndpoint, twitterLogin.CallbackHandler(cfg.Authorization, loginHandler.Login(), nil))
	mux.HandleFunc(handlers.TweetEndpoint, tweetHandler.Tweet)

	lambda.Start(agw.Handler(mux))
}
