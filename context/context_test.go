package context

import (
	"alexandreh2ag/go-task/config"
	"alexandreh2ag/go-task/types"
	"github.com/jonboulle/clockwork"
	"github.com/spf13/afero"
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
	fs := afero.NewOsFs()
	want := &Context{
		Config:   cfg,
		Logger:   logger,
		LogLevel: level,
		Clock:    clockwork.NewRealClock(),
		Fs:       fs,
	}
	got := DefaultContext()
	assert.NotNil(t, got.done)
	got.done = nil
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
	fs := afero.NewMemMapFs()
	want := &Context{
		Config:   cfg,
		Logger:   logger,
		LogLevel: level,
		Clock:    clockwork.NewRealClock(),
		Fs:       fs,
	}
	got := TestContext(nil)
	assert.NotNil(t, got.done)
	got.done = nil
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
	fs := afero.NewMemMapFs()
	want := &Context{
		Config:   cfg,
		Logger:   logger,
		LogLevel: level,
		Clock:    clockwork.NewRealClock(),
		Fs:       fs,
	}
	got := TestContext(io.Discard)
	assert.NotNil(t, got.done)
	got.done = nil
	assert.Equal(t, want, got)
}

func TestContext_Cancel(t *testing.T) {
	ctx := &Context{}
	ctx.done = make(chan bool)
	running := true
	go func() {
		select {
		case <-ctx.done:
			running = false
		}
	}()
	ctx.Cancel()
	assert.Equal(t, false, running)
}

func TestContext_Done(t *testing.T) {
	ctx := &Context{}
	ctx.done = make(chan bool)
	running := true
	go func() {
		select {
		case <-ctx.Done():
			running = false
		}
	}()
	ctx.done <- true
	assert.Equal(t, false, running)
}
