package mocks

import "github.com/stretchr/testify/mock"

import "golang.org/x/net/context"

type HealthCheckServicer struct {
	mock.Mock
}

// Run provides a mock function with given fields: _a0, _a1
func (_m *HealthCheckServicer) Run(_a0 context.Context, _a1 interface{}) (interface{}, error) {
	ret := _m.Called(_a0, _a1)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) interface{}); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
