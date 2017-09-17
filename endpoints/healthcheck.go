package endpoints

import (
	"golang.org/x/net/context"
)

// HealthCheckServicer has functions for performing healthchecks
type HealthCheckServicer interface {
	Run(context.Context, interface{}) (interface{}, error)
}

// HealthCheckResponse with status
type HealthCheckResponse struct {
	Status string `json:"ok"`
	GitSHA string `json:"git_sha"`
}

// HealthCheckEndpoint runs healthcheck
type HealthCheckEndpoint struct {
	gitSHA string // SHA1 hash of version this code
}

// NewHealthCheckEndpoint returns a new healthcheck endpoint
func NewHealthCheckEndpoint(gs string) HealthCheckServicer {
	return &HealthCheckEndpoint{gitSHA: gs}
}

// Run the healthcheck
func (h *HealthCheckEndpoint) Run(ctx context.Context, i interface{}) (interface{}, error) {
	return &HealthCheckResponse{Status: "OK", GitSHA: h.gitSHA}, nil
}
