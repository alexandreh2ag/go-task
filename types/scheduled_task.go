package types

import (
	"alexandreh2ag/go-task/log"
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"
)

const (
	Pending int = iota
	Succeed
	Failed
)

type ScheduledTasks = []*ScheduledTask

type ScheduledTask struct {
	Id         string `mapstructure:"id" validate:"required,alphanumunicode"`
	CronExpr   string `mapstructure:"expr" validate:"required,cron"`
	Command    string `mapstructure:"command" validate:"required"`
	User       string `mapstructure:"user" validate:"omitempty,required,alphanum"`
	Directory  string `mapstructure:"directory" validate:"omitempty,required,dirpath"`
	TaskResult *TaskResult

	Logger *slog.Logger
}

func (s *ScheduledTask) Execute() *TaskResult {
	var cmd *exec.Cmd
	result := &TaskResult{Status: Pending, Task: s}
	s.TaskResult = result

	args := strings.Fields(s.Command)
	if len(args) > 1 {
		cmd = exec.Command(args[0], args[1:]...)
	} else {
		cmd = exec.Command(args[0])
	}
	cmd.Dir = s.Directory
	cmd.Stdout = &result.Output
	cmd.Stderr = &result.Output

	s.Logger.Debug(fmt.Sprintf("Command (id: %s) run `%s` in %s by user %s", s.Id, s.Command, s.Directory, s.User))
	result.StartAt = time.Now()
	result.Error = cmd.Run()
	result.FinishAt = time.Now()

	if result.Error != nil {
		result.Status = Failed
	} else {
		result.Status = Succeed
	}
	s.Logger.Debug(fmt.Sprintf("Command (id: %s) end with status %s (%d)", s.Id, result.StatusString(), result.Status))

	return result
}

type TaskResult struct {
	Status   int
	Error    error
	Output   bytes.Buffer
	Task     *ScheduledTask
	StartAt  time.Time
	FinishAt time.Time
}

func (t *TaskResult) StatusString() string {
	switch t.Status {
	case Pending:
		return "pending"
	case Succeed:
		return "succeed"
	case Failed:
		return "failed"
	}
	return "unknown"
}

func PrepareScheduledTasks(tasks ScheduledTasks, logger *slog.Logger, user, workingDir string) {
	for _, task := range tasks {
		task.Logger = logger.With(log.TaskKey, task.Id)
		if task.User == "" {
			task.User = user
		}

		if task.Directory == "" {
			task.Directory = workingDir
		}
	}
}
