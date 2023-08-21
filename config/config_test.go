package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	got := NewConfig()
	assert.Equal(t, Config{}, got)
}

func TestDefaultConfig(t *testing.T) {
	got := DefaultConfig()
	assert.Equal(t, Config{}, got)
}
