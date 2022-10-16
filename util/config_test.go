package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	_, err := NewConfig()
	require.NoError(t, err)
}
