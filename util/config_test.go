package util

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	_, err := NewConfig()
	if err != nil {
		t.Errorf("NewConfig() error = %v", err)
		return
	}
}
