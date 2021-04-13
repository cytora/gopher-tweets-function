package main

import (
	"encoding/json"
	"os"

	"github.com/cytora/gopher-tweets-function/model"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var (
	twitterConsumerKey       = os.Getenv("TWITTER_CONSUMER_KEY")
	twitterConsumerSecret    = os.Getenv("TWITTER_CONSUMER_SECRET")
	twitterAccessToken       = os.Getenv("TWITTER_ACCESS_TOKEN")
	twitterAccessTokenSecret = os.Getenv("TWITTER_TOKEN_ACCESS_TOKEN_SECRET")

	clt *twitter.Client
)

func init() {
	if clt == nil {
		config := oauth1.NewConfig(twitterConsumerKey, twitterConsumerSecret)
		token := oauth1.NewToken(twitterAccessToken, twitterAccessTokenSecret)
		httpClient := config.Client(oauth1.NoContext, token)
		clt = twitter.NewClient(httpClient)
	}
}

func handler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req model.PostTweetRequest
	_ = json.Unmarshal([]byte(r.Body), &req)

	t, _, err := clt.Statuses.Update(req.Body, nil)
	if err != nil {
		return model.Response{
			Body:  nil,
			Error: err,
		}.APIGatewayProxyResponse()
	}

	res := model.PostTweetResponse{
		ID:      t.ID,
		Message: "tweet published",
	}

	return model.Response{
		StatusCode: 200,
		Body:       res,
		Error:      nil,
	}.APIGatewayProxyResponse()
}

func main() {
	lambda.Start(handler)
}
