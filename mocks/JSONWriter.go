package mocks

import "github.com/stretchr/testify/mock"

import "net/http"

type JSONWriter struct {
	mock.Mock
}

// WriteJSON provides a mock function with given fields: _a0
func (_m *JSONWriter) WriteJSON(_a0 http.ResponseWriter) {
	_m.Called(_a0)
}
