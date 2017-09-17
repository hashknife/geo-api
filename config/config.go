package config

import (
	"encoding/json"
	"os"

	kitlog "github.com/go-kit/kit/log"
)

// Tile38
type Tile38 struct {
	Hostname *string `json:"hostname"`
}

// Config
type Config struct {
	HTTPAddress                    *string `json:"http_address,omitempty"`
	HealthCheckAddress             *string `json:"healthcheck_address,omitempty"`
	HashknifeAuthToken             *string `json:"hashknife_auth_token,omitempty"`
	StatsdAddress                  *string `json:"statsd_address,omitempty"`
	StatsdReportingIntervalSeconds *int64  `json:"statsd_reporting_interval_seconds,omitempty"`
	Tile38                         *Tile38 `json:"tile38,omitempty"`
}

// String has Config implement the Stringer interface and allows for nicer/easier
// printing of the configuration
func (c *Config) String() string {
	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return ""
	}

	return string(j)
}

// Load reads the config file and creates a new value of type
// Config pointer
func Load(cf string, l kitlog.Logger) (*Config, error) {
	f, err := os.Open(cf)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}
	return &c, err
}
