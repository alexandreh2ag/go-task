package types

import (
	"bytes"
	"errors"
	"github.com/alexandreh2ag/go-task/log"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"testing"
	"time"
)

func Test_ScheduledTask_SuccessValidate(t *testing.T) {
	validate := validator.New()
	scheduled := ScheduledTask{
		Id:       "test",
		CronExpr: "* * * * *",
		Command:  "fake",
	}
	err := validate.Struct(scheduled)

	assert.NoError(t, err)
}

func Test_ScheduledTask_SuccessValidateWithOptionalData(t *testing.T) {
	validate := validator.New()
	scheduled := ScheduledTask{
		Id:        "test",
		CronExpr:  "* * * * *",
		Command:   "fake",
		Directory: "/tmp/test/",
	}
	err := validate.Struct(scheduled)

	assert.NoError(t, err)
}

func Test_ScheduledTask_ErrorValidate(t *testing.T) {
	validate := validator.New()
	scheduled := ScheduledTask{
		Id:       "test",
		CronExpr: "* * * * *",
	}
	err := validate.Struct(scheduled)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'Command' failed on the 'required' tag")
}

func Test_ScheduledTask_ErrorValidateComplex(t *testing.T) {
	validate := validator.New()
	scheduled := ScheduledTask{
		Id:        "test",
		CronExpr:  "wrong",
		Command:   "fake",
		Directory: "wrong",
	}
	err := validate.Struct(scheduled)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'CronExpr' failed on the 'cron' tag")
	assert.Contains(t, err.Error(), "Field validation for 'Directory' failed on the 'dirpath' tag")
}

func TestPrepareScheduledTasks(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	type args struct {
		tasks      ScheduledTasks
		logger     *slog.Logger
		user       string
		workingDir string
	}
	tests := []struct {
		name string
		args args
		want ScheduledTasks
	}{
		{
			name: "SuccessEmptyTasks",
			args: args{
				tasks:      ScheduledTasks{},
				logger:     logger,
				user:       "foo",
				workingDir: "/app/foo/",
			},
			want: ScheduledTasks{},
		},
		{
			name: "SuccessMultipleTasks",
			args: args{
				tasks: ScheduledTasks{
					&ScheduledTask{Id: "test", Command: "cmd", CronExpr: "* * * * *"},
					&ScheduledTask{Id: "test2", Command: "cmd", CronExpr: "* * * * *", Directory: "/app/bar/"},
				},
				logger:     logger,
				user:       "foo",
				workingDir: "/app/foo/",
			},
			want: ScheduledTasks{
				&ScheduledTask{Id: "test", Command: "cmd", CronExpr: "* * * * *", Directory: "/app/foo/", Logger: logger.With(log.TaskKey, "test")},
				&ScheduledTask{Id: "test2", Command: "cmd", CronExpr: "* * * * *", Directory: "/app/bar/", Logger: logger.With(log.TaskKey, "test2")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrepareScheduledTasks(tt.args.tasks, tt.args.logger, tt.args.workingDir)
			assert.Equal(t, tt.want, tt.args.tasks)
		})
	}
}

func TestTaskResult_StatusString(t1 *testing.T) {

	tests := []struct {
		name   string
		Status int
		want   string
	}{
		{
			name:   "SuccessWithPending",
			Status: Pending,
			want:   "pending",
		},
		{
			name:   "SuccessWithSucceed",
			Status: Succeed,
			want:   "succeed",
		},
		{
			name:   "SuccessWithFailed",
			Status: Failed,
			want:   "failed",
		},
		{
			name:   "SuccessWithUnknown",
			Status: -1,
			want:   "unknown",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TaskResult{
				Status: tt.Status,
			}
			assert.Equalf(t1, tt.want, t.StatusString(), "StatusString()")
		})
	}
}

func TestScheduledTask_Execute(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	type fields struct {
		Command string
		Logger  *slog.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   *TaskResult
	}{
		{
			name: "SuccessWithSimpleCommand",
			fields: fields{
				Command: "echo",
				Logger:  logger,
			},
			want: &TaskResult{
				Status:   Succeed,
				Output:   *bytes.NewBuffer([]byte("\n")),
				StartAt:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
				FinishAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "SuccessWithComplexCommand",
			fields: fields{
				Command: "echo test",
				Logger:  logger,
			},
			want: &TaskResult{
				Status:   Succeed,
				Output:   *bytes.NewBuffer([]byte("test\n")),
				StartAt:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
				FinishAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "FailedWithWrongCommand",
			fields: fields{
				Command: "wrong",
				Logger:  logger,
			},
			want: &TaskResult{
				Status:   Failed,
				Error:    &exec.Error{Name: "wrong", Err: errors.New("executable file not found in $PATH")},
				Output:   *bytes.NewBuffer(nil),
				StartAt:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
				FinishAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ScheduledTask{
				Command: tt.fields.Command,
				Logger:  tt.fields.Logger,
			}
			tt.want.Task = s
			res := s.Execute()
			assert.WithinDuration(t, s.TaskResult.StartAt, s.TaskResult.FinishAt, time.Second)
			assert.NotNil(t, s.TaskResult.Task)
			s.TaskResult.StartAt = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
			s.TaskResult.FinishAt = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
			assert.Equalf(t, tt.want, res, "Execute()")
		})
	}
}

func Test_splitCommand(t *testing.T) {

	tests := []struct {
		name    string
		command string
		want    []string
	}{
		{
			name:    "SuccessSimpleCommand",
			command: "echo",
			want:    []string{"echo"},
		},
		{
			name:    "SuccessComplexCommand",
			command: "echo test",
			want:    []string{"echo", "test"},
		},
		{
			name:    "SuccessComplexCommandWithDoubleQuote",
			command: "echo \"test\"",
			want:    []string{"echo", "\"test\""},
		},
		{
			name:    "SuccessComplexCommandWithQuote",
			command: "echo 'test'",
			want:    []string{"echo", "'test'"},
		},
		{
			name:    "SuccessComplexCommandWithQuoteAndDoubleQuote",
			command: "echo '\"test\"'",
			want:    []string{"echo", "'\"test\"'"},
		},
		{
			name:    "SuccessComplexCommandWithDoubleQuoteAndSpace",
			command: "echo ' test '",
			want:    []string{"echo", " test "},
		},
		{
			name:    "SuccessComplexCommandWithQuoteAndSpace",
			command: "echo ' test '",
			want:    []string{"echo", " test "},
		},
		{
			name:    "SuccessBashWrappedCommand",
			command: "bash -c 'echo test'",
			want:    []string{"bash", "-c", "echo test"},
		},
		{
			name:    "SuccessBashWrappedCommandAndSpace",
			command: "bash -c ' echo test '",
			want:    []string{"bash", "-c", " echo test "},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, splitCommand(tt.command), "splitCommand(%v)", tt.command)
		})
	}
}

func Test_getEnvVars(t *testing.T) {
	key := "GTASK_TESTING_GETENVVARS"
	value := "foo"
	_ = os.Setenv(key, value)
	tests := []struct {
		name      string
		key       string
		extraVars map[string]string
		want      string
	}{
		{
			name: "SuccessVarNotExist",
			key:  key + "WRONG",
			want: "",
		},
		{
			name: "SuccessOSVar",
			key:  key,
			want: "foo",
		},
		{
			name:      "SuccessGtaskVar",
			key:       GtaskIDKey,
			extraVars: map[string]string{GtaskIDKey: "bar"},
			want:      "bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getEnvVars(tt.extraVars)(tt.key), "getEnvVars(%v)", tt.extraVars)
		})
	}
}
