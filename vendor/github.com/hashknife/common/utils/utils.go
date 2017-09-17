package utils

import (
	"math/rand"
	"time"
)

// init runs at package initialization time
func init() {
	rand.Seed(time.Now().Unix())
}

var letterRunes = []rune("0123456789")

// DeliveryPin generates a pin to be entered by a package
// recipient at completion of a secure delivery
func DeliveryPin(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
