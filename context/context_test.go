package context

import (
	"alexandreh2ag/go-task/config"
	"alexandreh2ag/go-task/types"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultContext(t *testing.T) {

	cfg := &config.Config{
		Workers:   types.WorkerTasks{},
		Scheduled: types.ScheduledTasks{},
	}
	level := &slog.LevelVar{}
	level.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{AddSource: false, Level: level}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	want := &Context{
		Config:   cfg,
		Logger:   logger,
		LogLevel: level,
	}
	got := DefaultContext()

	assert.Equal(t, want, got)
}

func TestTestContext(t *testing.T) {

	cfg := &config.Config{
		Workers:   types.WorkerTasks{},
		Scheduled: types.ScheduledTasks{},
	}
	level := &slog.LevelVar{}
	level.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{AddSource: false, Level: level}
	logger := slog.New(slog.NewTextHandler(io.Discard, opts))
	want := &Context{
		Config:   cfg,
		Logger:   logger,
		LogLevel: level,
	}
	got := TestContext(nil)

	assert.Equal(t, want, got)
}
func TestTestContext_WithLogBuffer(t *testing.T) {

	cfg := &config.Config{
		Workers:   types.WorkerTasks{},
		Scheduled: types.ScheduledTasks{},
	}
	level := &slog.LevelVar{}
	level.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{AddSource: false, Level: level}
	logger := slog.New(slog.NewTextHandler(io.Discard, opts))
	want := &Context{
		Config:   cfg,
		Logger:   logger,
		LogLevel: level,
	}
	got := TestContext(io.Discard)

	assert.Equal(t, want, got)
}
