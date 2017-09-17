package models

import (
	"errors"
	"strconv"

	"github.com/jinzhu/gorm"
)

// PackageStatus contains the current status of a package
type PackageStatus string

const (
	// PackageDelivered represents the status of a delivered package
	PackageDelivered PackageStatus = "DELIVERED"
	// PackageInTransit represents the status of a package in transit
	PackageInTransit = "IN_TRANSIT"
	// PackageAwaitingCourier represents the status of a package waiting
	// to be picked up by a courier
	PackageAwaitingCourier = "AWAITING_COURIER"
	// PackageMissing represents a package whose location cannot be accounted for
	PackageMissing = "MISSING"
)

// Package reprsents and object sent for delivery
type Package struct {
	gorm.Model
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
	InsuranceID string `json:"insurance_id"`
	Status      string `json:"status"`
	Secure      bool   `json:"secure"`
}

// User
type User struct {
	gorm.Model
	Username  string `json:"username"`
	ProfileID string `json:"profile_id"`
}

// Courier
type Courier struct {
	gorm.Model
	ProfileID string `json:"profile_id"`
}

// Account
type Account struct {
	gorm.Model
}

// Sender
type Sender struct {
	gorm.Model
	ProfileID string `json:"profile_id"`
}

// Recipent
type Recipent struct {
	gorm.Model
	ProfileID string `json:"profile_id"`
}

// Profile
type Profile struct {
	gorm.Model
}

// Location represents a 2 point location represented by 64 bit values
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// NewLocation receives a string slice with position 0 being populated by longitude
// and position 1 being populated with latitidude
func NewLocation(points []float64) (*Location, error) {
	if len(points) < 1 {
		return nil, errors.New("ERR not enough points to create a location")
	}
	return &Location{
		Latitude:  points[1],
		Longitude: points[0],
	}, nil
}

// LatitudeToString converts latitude to a string
func (l *Location) LatitudeToString() string {
	return strconv.FormatFloat(l.Latitude, 'f', -1, 64)
}

// LongitudeToString converts longitude to a string
func (l *Location) LongitudeToString() string {
	return strconv.FormatFloat(l.Longitude, 'f', -1, 64)
}
