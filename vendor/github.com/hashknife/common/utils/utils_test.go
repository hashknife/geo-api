package utils

import (
	"testing"
)

// TestDeliveryPin
func TestDeliveryPin(t *testing.T) {
	for _, tc := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12} {
		dp := DeliveryPin(tc)
		if len(dp) != tc {
			t.Errorf("delivery pin not legnth specified: %d", tc)
		}
	}
}
