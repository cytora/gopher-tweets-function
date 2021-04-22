package handlers

import (
	"net/http"

	"github.com/cytora/gopher-tweets-function/model"
)

const HealthCheckEndpoint = "/health"

// Health indicates if the lambda is alive
func Health(w http.ResponseWriter, _ *http.Request) {
	r := model.Response{
		StatusCode: http.StatusOK,
		Body:       model.AliveResponse{Alive: true},
		Error:      nil,
	}
	r.Write(w)
}
