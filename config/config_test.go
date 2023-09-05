package config

import (
	"github.com/alexandreh2ag/go-task/types"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	got := NewConfig()
	assert.Equal(t, Config{Workers: types.WorkerTasks{}, Scheduled: types.ScheduledTasks{}}, got)
}

func TestDefaultConfig(t *testing.T) {
	got := DefaultConfig()
	assert.Equal(t, Config{Workers: types.WorkerTasks{}, Scheduled: types.ScheduledTasks{}}, got)
}

func Test_ConfigWorkers_SuccessValidateEmpty(t *testing.T) {
	cfg := DefaultConfig()
	validate := validator.New()
	err := validate.Struct(cfg)

	assert.NoError(t, err)
}

func Test_ConfigWorkers_SuccessValidate(t *testing.T) {
	cfg := DefaultConfig()
	validate := validator.New()
	cfg.Workers = types.WorkerTasks{
		{Id: "test", Command: "fake"},
		{Id: "test2", Command: "fake"},
	}
	cfg.Scheduled = types.ScheduledTasks{
		{Id: "test", CronExpr: "* * * * *", Command: "fake"},
		{Id: "test2", CronExpr: "* * * * *", Command: "fake"},
	}
	err := validate.Struct(cfg)

	assert.NoError(t, err)
}

func Test_ConfigWorkers_ErrorValidateDivet(t *testing.T) {
	cfg := DefaultConfig()
	validate := validator.New()
	cfg.Workers = types.WorkerTasks{
		{Id: "test", Command: "fake"},
		{Id: "test2"},
	}
	cfg.Scheduled = types.ScheduledTasks{
		{Id: "test", CronExpr: "* * * * *", Command: "fake"},
		{Id: "test2", CronExpr: "wrong", Command: "fake"},
	}
	err := validate.Struct(cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Config.Workers[1].Command' Error:Field validation for 'Command' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Config.Scheduled[1].CronExpr' Error:Field validation for 'CronExpr' failed on the 'cron' tag")
}

func Test_ConfigWorkers_ErrorValidateDuplicateId(t *testing.T) {
	cfg := DefaultConfig()
	validate := validator.New()
	cfg.Workers = types.WorkerTasks{
		{Id: "test", Command: "fake"},
		{Id: "test", Command: "fake"},
	}
	cfg.Scheduled = types.ScheduledTasks{
		{Id: "test", CronExpr: "* * * * *", Command: "fake"},
		{Id: "test", CronExpr: "* * * * *", Command: "fake"},
	}
	err := validate.Struct(cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error:Field validation for 'Workers' failed on the 'unique' tag")
}
