package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSONWriter
type JSONWriter interface {
	WriteJSON(http.ResponseWriter)
}

type ContextError struct {
	Message    string                 `json:"message"`
	Context    map[string]interface{} `json:"context"`
	StatusCode int                    `json:"status_code"`
}

func NewContextError(msg string, status int, ctx map[string]interface{}) *ContextError {
	return &ContextError{
		Message:    msg,
		StatusCode: status,
		Context:    ctx,
	}
}

// Error provides
func (e *ContextError) Error() string {
	return fmt.Sprintf("Message: %s Context: %v StatusCode: %d", e.Message, e.Context, e.StatusCode)
}

// WriteJSON
func (e *ContextError) WriteJSON(rw http.ResponseWriter) {
	err := map[string]interface{}{
		"error": e.Message,
	}
	if e.Context != nil {
		for key, value := range e.Context {
			err[key] = value
		}
	}
	WriteJSON(e.StatusCode, rw, err)
}

// WriteJSON
func WriteJSON(status int, rw http.ResponseWriter, v interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	if err := json.NewEncoder(rw).Encode(v); err != nil {
		panic(err)
	}
}

// NewForbiddenError is run when a user tries to access an endpoint
// but isn't allowed to do so
func NewForbiddenError() *ContextError {
	return &ContextError{
		Message:    "Authorization Token Incorrect",
		StatusCode: 403,
	}
}
