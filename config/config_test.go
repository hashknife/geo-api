package config

import (
	"io/ioutil"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/suite"
)

// ConfigTestSuite
type ConfigTestSuite struct {
	suite.Suite
	logger kitlog.Logger
	config *Config
}

// SetupSuite
func (c *ConfigTestSuite) SetupSuite() {
	c.logger = kitlog.NewJSONLogger(ioutil.Discard)
	c.config = new(Config)
}

// SetupTest
func (c *ConfigTestSuite) SetupTest() {}

// TearDownTest
func (c *ConfigTestSuite) TearDownTest() {
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, &ConfigTestSuite{})
}

// TestLoad_Success
func (c *ConfigTestSuite) TestLoad_Success() {
	conf, err := Load("testdata/config.json", c.logger)
	c.Require().NoError(err)
	c.Require().NotNil(conf)
}

// TestString_Success
func (c *ConfigTestSuite) TestString_Success() {
	c.Require().NotEmpty(c.config.String())
}

// TestString_NullConfig_Failure
func (c *ConfigTestSuite) TestString_NullConfig_Failure() {
	var conf *Config
	conf = nil
	s := conf.String()
	c.Require().Equal("null", s)
}

// TestString_Empty_Failure
func (c *ConfigTestSuite) TestString_Empty_Failure() {
	var conf Config
	s := conf.String()
	c.Require().Equal("{}", s)
}
