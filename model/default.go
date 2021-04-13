package model

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

var (
	defaultHeaders = map[string]string{
		"Access-Control-Allow-Headers": "*",
		"Access-Control-Allow-Methods": "OPTIONS,POST,GET",
		"Access-Control-Allow-Origin":  "*",
	}
)

type Response struct {
	StatusCode int
	Body       interface{}
	Error      error
}

func (r Response) APIGatewayProxyResponse() (events.APIGatewayProxyResponse, error) {
	if r.Error != nil {
		return events.APIGatewayProxyResponse{}, r.Error
	}

	body, err := json.Marshal(r.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, r.Error
	}

	return events.APIGatewayProxyResponse{
		Headers:    defaultHeaders,
		StatusCode: r.StatusCode,
		Body:       string(body),
	}, nil
}
