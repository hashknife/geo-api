package mocks

import "github.com/stretchr/testify/mock"

import "github.com/hashknife/common/models"

type Tile38er struct {
	mock.Mock
}

// RemoveCourier provides a mock function with given fields: accoutID, courierID
func (_m *Tile38er) RemoveCourier(accoutID string, courierID string) error {
	ret := _m.Called(accoutID, courierID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(accoutID, courierID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveAccount provides a mock function with given fields: accoutID
func (_m *Tile38er) RemoveAccount(accoutID string) error {
	ret := _m.Called(accoutID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(accoutID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateCourierLocaton provides a mock function with given fields: accoutID, courierID, lat, lon
func (_m *Tile38er) UpdateCourierLocaton(accoutID string, courierID string, lat string, lon string) error {
	ret := _m.Called(accoutID, courierID, lat, lon)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string) error); ok {
		r0 = rf(accoutID, courierID, lat, lon)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddFence provides a mock function with given fields: accoutID, courierID, lat, lon
func (_m *Tile38er) AddFence(accoutID string, courierID string, lat string, lon string) error {
	ret := _m.Called(accoutID, courierID, lat, lon)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string) error); ok {
		r0 = rf(accoutID, courierID, lat, lon)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CourierLocation provides a mock function with given fields: accoutID, courierID
func (_m *Tile38er) CourierLocation(accoutID string, courierID string) (*models.Location, error) {
	ret := _m.Called(accoutID, courierID)

	var r0 *models.Location
	if rf, ok := ret.Get(0).(func(string, string) *models.Location); ok {
		r0 = rf(accoutID, courierID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Location)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(accoutID, courierID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ServerStatistics provides a mock function with given fields:
func (_m *Tile38er) ServerStatistics() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
