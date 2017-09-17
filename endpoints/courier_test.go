package endpoints

import (
	"context"
	"errors"
	"testing"

	"github.com/hashknife/common/models"
	"github.com/hashknife/geo-api/config"
	"github.com/hashknife/geo-api/mocks"
	"github.com/stretchr/testify/suite"
)

// CourierTestSuite
type CourierServiceTestSuite struct {
	suite.Suite
	conf       *config.Config
	endpoint   CourierServicer
	mockTile38 *mocks.Tile38er
	accountID  string
	courierID  string
}

// SetupSuite runs code needed for the test suite
func (c *CourierServiceTestSuite) SetupSuite() {
	hn := "localhost:1234"
	c.conf = &config.Config{
		Tile38: &config.Tile38{
			Hostname: &hn,
		},
	}
	c.mockTile38 = &mocks.Tile38er{}
	c.endpoint = NewCourierService(c.conf, c.mockTile38)
	c.accountID = "asdf"
	c.courierID = "1234"
}

// resetMocks
func (c *CourierServiceTestSuite) resetMocks() {
	c.mockTile38 = &mocks.Tile38er{}
}

// assertMockExpectations
func (c *CourierServiceTestSuite) assertMockExpectations() {
	c.mockTile38.AssertExpectations(c.T())
}

// SetupTest
func (c *CourierServiceTestSuite) SetupTest() {
	c.resetMocks()
}

// TearDownTest
func (c *CourierServiceTestSuite) TearDownTest() {
	c.assertMockExpectations()
}

// TestCourierServiceTestSuite
func TestCourierServiceTestSuite(t *testing.T) {
	suite.Run(t, &CourierServiceTestSuite{})
}

// TestRetrieve_Success
func (c *CourierServiceTestSuite) TestLocation_Success() {
	loc, err := models.NewLocation([]float64{66.7, 55.6})
	c.NoError(err)
	c.mockTile38.On("CourierLocation", c.accountID, c.courierID).Return(loc, nil)
	req := &CourierServiceLocationRequest{
		AccountID: c.accountID,
		CourierID: c.courierID,
	}
	ce := NewCourierService(c.conf, c.mockTile38)
	resp, err := ce.Location(context.Background(), req)
	c.NoError(err)
	c.NotNil(resp)
}

// TestLocation_Failure
func (c *CourierServiceTestSuite) TestLocation_Failure() {
	c.mockTile38.On("CourierLocation", c.accountID, c.courierID).Return(nil, errors.New("errored"))
	req := &CourierServiceLocationRequest{
		AccountID: c.accountID,
		CourierID: c.courierID,
	}
	ce := NewCourierService(c.conf, c.mockTile38)
	resp, err := ce.Location(context.Background(), req)
	c.Error(err)
	c.Nil(resp)
}

// TestUpdateLocation_Success
func (c *CourierServiceTestSuite) TestUpdateLocation_Success() {
	c.mockTile38.On("UpdateCourierLocaton", c.accountID, c.courierID, "11.1", "12.2").Return(nil)
	req := &CourierServiceUpdateRequest{
		AccountID: c.accountID,
		CourierID: c.courierID,
		Latitude:  11.1,
		Longitude: 12.2,
	}
	ce := NewCourierService(c.conf, c.mockTile38)
	resp, err := ce.UpdateLocation(context.Background(), req)
	c.NoError(err)
	c.NotNil(resp)
}

// TestUpdateLocation_Failure
func (c *CourierServiceTestSuite) TestUpdateLocation_Failure() {
	c.mockTile38.On("UpdateCourierLocaton", c.accountID, c.courierID, "11.1", "12.2").Return(errors.New("errored"))
	req := &CourierServiceUpdateRequest{
		AccountID: c.accountID,
		CourierID: c.courierID,
		Latitude:  11.1,
		Longitude: 12.2,
	}
	ce := NewCourierService(c.conf, c.mockTile38)
	resp, err := ce.UpdateLocation(context.Background(), req)
	c.Error(err)
	c.Nil(resp)
}
