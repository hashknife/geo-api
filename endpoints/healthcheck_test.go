package endpoints

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

// HealthCheckTestSuite
type HealthCheckTestSuite struct {
	suite.Suite
	endpoint HealthCheckServicer
}

// SetupSuite runs code needed for the test suite
func (h *HealthCheckTestSuite) SetupSuite() {
	h.endpoint = NewHealthCheckEndpoint("version")
}

// TestHealthCheckTestSuite
func TestHealthCheckTestSuite(t *testing.T) {
	suite.Run(t, &HealthCheckTestSuite{})
}

// TestHealthCheck
func (h *HealthCheckTestSuite) TestHealthCheck() {
	resp, err := h.endpoint.Run(context.Background(), nil)
	h.NoError(err)
	h.NotNil(resp)
}
