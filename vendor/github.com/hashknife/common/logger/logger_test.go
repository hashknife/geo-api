package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestItCanLogWithoutBlowingUp
func TestValidateLogging(t *testing.T) {
	l := SetupJSONLogger("some-hashknife-service")
	err := l.Log("abc", "123")
	require.NoError(t, err)
}
