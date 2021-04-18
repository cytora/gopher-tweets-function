package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cytora/gopher-tweets-function/model"
	"github.com/davyzhang/agw"
	"github.com/dghubble/go-twitter/twitter"
	oauth1Login "github.com/dghubble/gologin/v2/oauth1"
	twitterLogin "github.com/dghubble/gologin/v2/twitter"
	"github.com/dghubble/oauth1"
	twitterOAuth1 "github.com/dghubble/oauth1/twitter"
	"github.com/dghubble/sessions"
	"github.com/gorilla/mux"
)

var (
	twitterConsumerKey    = os.Getenv("TWITTER_CONSUMER_KEY")
	twitterConsumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")
	twitterLoginCallback  = os.Getenv("TWITTER_LOGIN_CALLBACK")
	sessionSecret         = os.Getenv("SESSION_SECRET_KEY")
	// sessionStore encodes and decodes session data stored in signed cookies
	store = sessions.NewCookieStore([]byte(sessionSecret), nil)
	clt   *twitter.Client
)

const (
	sessionName         = "gopher-tweets"
	sessionAccessKey    = "accessKey"
	sessionAccessSecret = "accessSecret"
	sessionUserKey      = "twitterID"
	sessionUsername     = "twitterUsername"
	authEndpoint        = "/twitter/login"
)

func getTwitterClient(req *http.Request) (*twitter.Client, error) {
	config := oauth1.NewConfig(twitterConsumerKey, twitterConsumerSecret)
	session, err := store.Get(req, sessionName)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	access := session.Values[sessionAccessKey]
	secret := session.Values[sessionAccessSecret]
	if access == nil || secret == nil {
		return nil, errors.New("access key or secret access key not set")
	}
	token := oauth1.NewToken(access.(string), secret.(string))
	httpClient := config.Client(req.Context(), token)
	return twitter.NewClient(httpClient), nil
}

func tweetHandler(w http.ResponseWriter, req *http.Request) {
	clt, err := getTwitterClient(req)
	if err != nil {
		r := model.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       nil,
			Error:      errors.New("error getting Twitter client"),
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
			Error:      errors.New("error getting Twitter client"),
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
			Message: "tweet published",
		},
		Error: nil,
	}
	r.Write(w)
	return
}

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/hello", helloHandler)
	// 1. Register Twitter login and callback handlers
	oauth1Config := &oauth1.Config{
		ConsumerKey:    twitterConsumerKey,
		ConsumerSecret: twitterConsumerSecret,
		CallbackURL:    twitterLoginCallback,
		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
	}

	mux.Handle(authEndpoint, twitterLogin.LoginHandler(oauth1Config, nil))
	mux.Handle("/twitter/callback", twitterLogin.CallbackHandler(oauth1Config, issueSession(), nil))
	mux.HandleFunc("/profile", profileHandler)
	mux.HandleFunc("/twitter/reverse", reverseHandler)
	mux.HandleFunc("/tweet", tweetHandler)

	lambda.Start(agw.Handler(mux))
}

// helloHandler indicates if the lambda is alive
func helloHandler(w http.ResponseWriter, _ *http.Request) {
	r := model.Response{
		StatusCode: http.StatusOK,
		Body:       model.AliveResponse{Alive: true},
		Error:      nil,
	}
	r.Write(w)
}

// issueSession issues a cookie session after successful Twitter login
func issueSession() http.Handler {
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

		// 2. Implement a success handler to issue some form of session
		session := store.New(sessionName)
		session.Values[sessionAccessKey] = accessToken
		session.Values[sessionAccessSecret] = accessSecret
		session.Values[sessionUserKey] = twitterUser.ID
		session.Values[sessionUsername] = twitterUser.ScreenName
		if err := session.Save(w); err != nil {
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
			Body:       model.TwitterAuthMessage{Message: "twitter auth success"},
			Error:      nil,
		}
		r.Write(w)

	}
	return http.HandlerFunc(fn)
}

// profileHandler shows a personal profile or a login button.
func profileHandler(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, sessionName)
	if err != nil {
		r := model.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       model.TwitterProfile{Error: err.Error()},
			Error:      nil,
		}
		r.Write(w)
		return
	}
	user := session.Values[sessionUsername]
	if user == nil {
		r := model.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       model.TwitterProfile{Error: "nil user"},
			Error:      nil,
		}
		r.Write(w)
		return
	}

	// authenticated profile
	r := model.Response{
		StatusCode: http.StatusOK,
		Body:       model.TwitterProfile{UserName: user.(string)},
		Error:      nil,
	}
	r.Write(w)
}

// reverseHandler shows a personal profile or a login button.
func reverseHandler(w http.ResponseWriter, req *http.Request) {
	//session, err := store.Get(req, sessionName)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	w.(*agw.LPResponse).WriteBody(model.Response{
	//		StatusCode: http.StatusInternalServerError,
	//		Body:       model.TwitterProfile{Error: err.Error()},
	//		Error:      nil,
	//	}, false)
	//	return
	//}
	//user := session.Values[sessionUsername]
	//if user == nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	w.(*agw.LPResponse).WriteBody(model.Response{
	//		StatusCode: http.StatusInternalServerError,
	//		Body:       model.TwitterProfile{Error: "nil user"},
	//		Error:      nil,
	//	}, false)
	//	return
	//}

	// authenticated profile
	r := model.Response{
		StatusCode: http.StatusOK,
		Body:       nil,
		Error:      nil,
	}
	r.Write(w)
}
