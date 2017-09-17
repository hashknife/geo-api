package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/hashknife/common/models"
)

// Geometry
type Geometry struct {
	Type        string    `json:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
}

// Properties
type Properties struct {
	Name string `json:"name,omitempty"`
}

// GeoJSONResponse
type GeoJSONResponse struct {
	Type        string     `json:"type,omitempty"`
	Coordinates []float64  `json:"coordinates,omitempty"`
	Geometry    Geometry   `json:"geometry,omitempty"`
	Properties  Properties `json:"properties,omitempty"`
}

// Tile38er
type Tile38er interface {
	RemoveCourier(accoutID, courierID string) error
	RemoveAccount(accoutID string) error
	UpdateCourierLocaton(accoutID, courierID, lat, lon string) error
	AddFence(accoutID, courierID, lat, lon string) error
	CourierLocation(accoutID, courierID string) (*models.Location, error)
	ServerStatistics() error
}

// Tile38
type Tile38 struct {
	pool *redis.Pool
}

// compile time validation
var _ Tile38er = (*Tile38)(nil)

// NewTile38
func NewTile38(host string) *Tile38 {
	return &Tile38{
		pool: newPool(host),
	}
}

// newPool
func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// RemoveCourier
func (t *Tile38) RemoveCourier(accoutID, courierID string) error {
	conn := t.pool.Get()
	defer conn.Close()
	_, err := conn.Do("GET", accoutID, courierID)
	if err != nil {
		return err
	}
	return nil
}

// RemoveAccount
func (t *Tile38) RemoveAccount(accoutID string) error {
	conn := t.pool.Get()
	defer conn.Close()
	_, err := conn.Do("GET", accoutID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateCourierLocaton
func (t *Tile38) UpdateCourierLocaton(accoutID, courierID, lat, lon string) error {
	conn := t.pool.Get()
	defer conn.Close()
	_, err := conn.Do("GET", accoutID, courierID, lat, lon)
	if err != nil {
		return err
	}
	return nil
}

// UpdateCourierLocaton
func (t *Tile38) AddFence(accoutID, courierID, lat, lon string) error {
	conn := t.pool.Get()
	defer conn.Close()
	_, err := conn.Do("GET", accoutID, courierID, lat, lon)
	if err != nil {
		return err
	}
	return nil
}

// UpdateCourierLocaton
func (t *Tile38) CourierLocation(accoutID, courierID string) (*models.Location, error) {
	conn := t.pool.Get()
	defer conn.Close()
	r, err := conn.Do("GET", accoutID, courierID)
	if err != nil {
		return nil, err
	}
	var res GeoJSONResponse
	if v, ok := r.([]byte); ok {
		if err := json.Unmarshal(v, &res); err != nil {
			return nil, err
		}
	}
	location, err := models.NewLocation(res.Coordinates)
	if err != nil {
		return nil, err
	}
	return location, nil
}

// ServerStatistics
func (t *Tile38) ServerStatistics() error {
	//return t.call("SERVER")
	return nil
}

// runOp runs a the given command and args against the configured cluster
func (t *Tile38) call(res interface{}, cmd string, args ...interface{}) error {
	conn := t.pool.Get()
	defer conn.Close()
	switch cmd {
	case "GET":
		fmt.Printf("%s %+v\n", cmd, args)
		r, err := conn.Do(cmd, args)
		if err != nil {
			return err
		}
		if v, ok := r.([]byte); ok {
			if err := json.Unmarshal(v, res); err != nil {
				return err
			}
		}
		return nil
	}
	_, err := conn.Do(cmd, args)
	if err != nil {
		return err
	}
	return nil
}
