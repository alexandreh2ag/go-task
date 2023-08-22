package config

import "alexandreh2ag/go-task/types"

type Config struct {
	//LogLevel string `mapstructure:"log_level"`
	Workers   types.WorkerTasks    `mapstructure:"workers" validate:"omitempty,required,unique=Id,dive"`
	Scheduled types.ScheduledTasks `mapstructure:"scheduled" validate:"omitempty,required,unique=Id,dive"`
}

func NewConfig() Config {
	return Config{
		Workers:   types.WorkerTasks{},
		Scheduled: types.ScheduledTasks{},
	}
}

func DefaultConfig() Config {
	cfg := NewConfig()

	return cfg
}
