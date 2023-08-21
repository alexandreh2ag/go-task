package context

import (
	"alexandreh2ag/go-task/config"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {

	cfg := &config.Config{}
	level := &slog.LevelVar{}
	level.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{AddSource: false, Level: level}
	logger := slog.New(slog.NewTextHandler(io.Discard, opts))
	want := &Context{
		Config:   cfg,
		Logger:   logger,
		LogLevel: level,
	}
	got := NewContext(logger, level, cfg)

	assert.Equal(t, want, got)
}
