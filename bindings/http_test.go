package bindings

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/briandowns/hashknife/geo-api/config"
	"github.com/briandowns/hashknife/geo-api/endpoints"
	"github.com/briandowns/hashknife/geo-api/mocks"
	kitlog "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// HTTPBindingsSuite
type HTTPBindingsSuite struct {
	suite.Suite
	config             *config.Config
	logger             kitlog.Logger
	mockHealthCheckSvc *mocks.HealthCheckServicer
}

// TestHTTPBindingsSuite
func TestHTTPBindingsSuite(t *testing.T) {
	suite.Run(t, new(HTTPBindingsSuite))
}

// SetupSuite
func (s *HTTPBindingsSuite) SetupSuite() {
	statsdAddress := "127.0.0.1:8125"
	statsdInterval := int64(5)
	s.config = &config.Config{
		StatsdAddress:                  &statsdAddress,
		StatsdReportingIntervalSeconds: &statsdInterval,
	}
	s.logger = kitlog.NewJSONLogger(ioutil.Discard)
}

// SetupTest
func (s *HTTPBindingsSuite) SetupTest() {
	s.resetMocks()
}

// TearDownTest
func (s *HTTPBindingsSuite) TearDownTest() {
	s.assertMockExpectations()
}

// resetMocks
func (s *HTTPBindingsSuite) resetMocks() {
	s.mockHealthCheckSvc = new(mocks.HealthCheckServicer)
}

// assertMockExpectations
func (s *HTTPBindingsSuite) assertMockExpectations() {
	s.mockHealthCheckSvc.AssertExpectations(s.T())
}

// TestHealthCheckSuccess validates the healthcheck endpoint succeeds
func (s *HTTPBindingsSuite) TestHealthCheckSuccess() {
	// GIVEN
	s.mockHealthCheckSvc.On("Run", mock.Anything, mock.Anything).Return(&endpoints.HealthCheckResponse{}, nil)
	router := createHealthCheckRouter(context.Background(), s.logger, s.mockHealthCheckSvc)
	req, err := http.NewRequest("GET", "http://application.hashknife.io/geo-api/healthcheck", nil)
	s.Nil(err)
	var match mux.RouteMatch
	ok := router.Match(req, &match)
	s.Require().True(ok)
	s.Require().NotNil(match)

	// WHEN
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	// THEN
	s.Equal(http.StatusOK, response.Code)
	s.assertMockExpectations()
}

// TestHealthCheckError validates that the health check endpoint fails as expected
func (s *HTTPBindingsSuite) TestHealthCheckError() {
	// GIVEN
	s.mockHealthCheckSvc.On("Run", mock.Anything, mock.Anything).Return((*endpoints.HealthCheckResponse)(nil), errors.New("HealthCheck failed"))
	router := createHealthCheckRouter(context.Background(), s.logger, s.mockHealthCheckSvc)
	req, err := http.NewRequest("GET", "http://application.hashknife.io/geo-api/healthcheck", nil)
	s.Require().NoError(err)
	var match mux.RouteMatch
	ok := router.Match(req, &match)
	s.True(ok)
	s.Require().NotNil(match)

	// WHEN
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	// THEN
	s.Equal(http.StatusInternalServerError, response.Code)
	s.assertMockExpectations()
}
