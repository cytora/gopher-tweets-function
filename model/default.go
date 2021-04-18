package model

import (
	"net/http"

	"github.com/davyzhang/agw"
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

func (r Response) Write(w http.ResponseWriter) {
	resp := w.(*agw.LPResponse)
	resp.WriteHeader(r.StatusCode)
	h := resp.Header()
	for name, value := range defaultHeaders {
		h[name] = []string{value}
	}
	resp.WriteBody(r.Body, false)
}
