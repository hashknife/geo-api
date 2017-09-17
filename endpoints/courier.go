package endpoints

import (
	"context"
	"errors"
	"strconv"

	"github.com/briandowns/hashknife/common/services"
	"github.com/briandowns/hashknife/geo-api/config"
)

// CourierServicer
type CourierServicer interface {
	Location(ctx context.Context, i interface{}) (interface{}, error)
	UpdateLocation(ctx context.Context, i interface{}) (interface{}, error)
}

// CourierService
type CourierService struct {
	conf *config.Config
	tc   services.Tile38er
}

// compile time validation
var _ CourierServicer = (*CourierService)(nil)

// NewCourierService
func NewCourierService(c *config.Config, t services.Tile38er) CourierServicer {
	return &CourierService{conf: c, tc: t}
}

// CourierServiceLocationRequest
type CourierServiceLocationRequest struct {
	AccountID string `json:"account_id"`
	CourierID string `json:"courier_id"`
}

// CourierServiceLocationResponse
type CourierServiceLocationResponse struct {
	CourierID string  `json:"courier_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// CourierServiceUpdateRequest
type CourierServiceUpdateRequest struct {
	AccountID string  `json:"account_id"`
	CourierID string  `json:"courier_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// CourierServiceUpdateResponse
type CourierServiceUpdateResponse struct {
	CourierID string `json:"courier_id"`
	Status    string `json:"status"`
}

// Location gets the given courier's location
func (c *CourierService) Location(ctx context.Context, i interface{}) (interface{}, error) {
	req, ok := i.(*CourierServiceLocationRequest)
	if !ok {
		return nil, errors.New("unable to convert request to CourierServiceLocationRequest type")
	}
	l, err := c.tc.CourierLocation(req.AccountID, req.CourierID)
	if err != nil {
		return nil, err
	}
	return &CourierServiceLocationResponse{req.CourierID, l.Latitude, l.Longitude}, nil
}

// UpdateLocation updates the given courier's location
func (c *CourierService) UpdateLocation(ctx context.Context, i interface{}) (interface{}, error) {
	req, ok := i.(*CourierServiceUpdateRequest)
	if !ok {
		return nil, errors.New("unable to convert request to CourierServiceUpdateRequest type")
	}
	lat := strconv.FormatFloat(req.Latitude, 'f', -1, 64)
	lon := strconv.FormatFloat(req.Longitude, 'f', -1, 64)
	if err := c.tc.UpdateCourierLocaton(req.AccountID, req.CourierID, lat, lon); err != nil {
		return nil, err
	}
	return &CourierServiceUpdateResponse{req.CourierID, "OK"}, nil
}
