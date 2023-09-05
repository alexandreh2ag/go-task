package context

import (
	"github.com/alexandreh2ag/go-task/config"
	"github.com/jonboulle/clockwork"
	"github.com/spf13/afero"
	"io"
	"log/slog"

	"os"
)

type Context struct {
	Logger   *slog.Logger
	LogLevel *slog.LevelVar
	Config   *config.Config
	Clock    clockwork.Clock
	Fs       afero.Fs
	done     chan bool
}

func (c *Context) Cancel() {
	c.done <- true
}

func (c *Context) Done() <-chan bool {
	return c.done
}

func DefaultContext() *Context {
	cfg := config.DefaultConfig()
	level := &slog.LevelVar{}
	level.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{AddSource: false, Level: level}
	return &Context{
		Logger:   slog.New(slog.NewTextHandler(os.Stdout, opts)),
		LogLevel: level,
		Config:   &cfg,
		Clock:    clockwork.NewRealClock(),
		Fs:       afero.NewOsFs(),
		done:     make(chan bool),
	}
}

func TestContext(logBuffer io.Writer) *Context {
	if logBuffer == nil {
		logBuffer = io.Discard
	}
	cfg := config.DefaultConfig()
	level := &slog.LevelVar{}
	level.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{AddSource: false, Level: level}
	return &Context{
		Logger:   slog.New(slog.NewTextHandler(logBuffer, opts)),
		LogLevel: level,
		Config:   &cfg,
		Clock:    clockwork.NewRealClock(),
		Fs:       afero.NewMemMapFs(),
		done:     make(chan bool),
	}
}
