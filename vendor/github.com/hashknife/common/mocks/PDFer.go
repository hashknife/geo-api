package mocks

import "github.com/stretchr/testify/mock"

type PDFer struct {
	mock.Mock
}

// PDF provides a mock function with given fields:
func (_m *PDFer) PDF() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
