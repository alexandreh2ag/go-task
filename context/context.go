package context

import (
	"alexandreh2ag/go-task/config"
	"github.com/spf13/afero"
	"io"
	"log/slog"

	"os"
)

type Context struct {
	Logger   *slog.Logger
	LogLevel *slog.LevelVar
	Config   *config.Config
	Fs       afero.Fs
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
		Fs:       afero.NewOsFs(),
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
		Fs:       afero.NewMemMapFs(),
	}
}
