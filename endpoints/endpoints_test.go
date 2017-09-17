package endpoints

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// SchedulerTestSuite
type ErrorsTestSuite struct {
	suite.Suite
}

// TestNewContextError
func (e *ErrorsTestSuite) TestNewContextError() {
	ctx := make(map[string]interface{})
	ce := NewContextError("", 0, ctx)
	e.Require().NotEmpty(ce)
}

// TestContextError_Error
func (e *ErrorsTestSuite) TestContextError_Error() {
	ce := NewContextError("", 0, make(map[string]interface{}))
	e.Require().IsType("string", ce.Error())
}

// TestWriteJSON
func (e *ErrorsTestSuite) TestWriteJSON() {}

// TestNewForbiddenError
func (e *ErrorsTestSuite) TestNewForbiddenError() {
	ce := NewForbiddenError()
	e.Require().NotNil(ce)
}

// TestErrorsTestSuite
func TestErrorsTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorsTestSuite))
}
