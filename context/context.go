package context

import (
	"alexandreh2ag/go-task/config"
	"log/slog"
	"os"
)

type Context struct {
	Logger   *slog.Logger
	LogLevel *slog.LevelVar
	Config   *config.Config
}

func NewContext(logger *slog.Logger, logLevel *slog.LevelVar, cfg *config.Config) *Context {
	return &Context{Logger: logger, LogLevel: logLevel, Config: cfg}
}

func DefaultContext() *Context {
	cfg := config.DefaultConfig()
	level := &slog.LevelVar{}
	level.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{AddSource: false, Level: level}
	return NewContext(slog.New(slog.NewTextHandler(os.Stdout, opts)), level, &cfg)
}
