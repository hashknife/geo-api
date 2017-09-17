package services

import (
	"testing"

	"github.com/hashknife/api/config"
	"github.com/hashknife/common/utils"
	"github.com/stretchr/testify/suite"
)

// Tile38TestSuite
type Tile38TestSuite struct {
	suite.Suite
	conf    *config.Config
	service Tile38er
}

// SetupSuite runs code needed for the test suite
func (t *Tile38TestSuite) SetupSuite() {
	t.conf = &config.Config{
		Tile38: &config.Tile38{
			Endpoint: utils.String("some-endpoint"),
		},
	}
	t.service = NewTile38(*t.conf.Tile38.Endpoint)
}

// TestTile38TestSuite
func TestTile38TestSuite(t *testing.T) {
	suite.Run(t, &Tile38TestSuite{})
}
