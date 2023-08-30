package schedule

import (
	"alexandreh2ag/go-task/cli/flags"
	"alexandreh2ag/go-task/context"
	"alexandreh2ag/go-task/types"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestGetScheduleRunCmd_SuccessWithEmptyScheduledTasks(t *testing.T) {
	ctx := context.TestContext(io.Discard)

	cmd := GetScheduleRunCmd(ctx)
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestGetScheduleRunCmd_SuccessWithScheduledTasks(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	ctx.Config.Scheduled = types.ScheduledTasks{
		&types.ScheduledTask{Id: "test", Command: "echo", CronExpr: "0 0 * * *"},
	}
	cmd := GetScheduleRunCmd(ctx)
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestGetScheduleRunCmd_SuccessWithTimezoneOpt(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	ctx.Config.Scheduled = types.ScheduledTasks{
		&types.ScheduledTask{Id: "test", Command: "echo", CronExpr: "0 0 * * *"},
	}
	cmd := GetScheduleRunCmd(ctx)

	cmd.SetArgs([]string{"--" + TimeZone, "Europe/Paris"})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestGetScheduleRunCmd_SuccessWithUserOpt(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	task := &types.ScheduledTask{Id: "test", Command: "echo", CronExpr: "0 0 * * *"}
	ctx.Config.Scheduled = types.ScheduledTasks{
		task,
	}
	cmd := GetScheduleRunCmd(ctx)

	cmd.SetArgs([]string{"--" + flags.User, "foo"})
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "foo", task.User)
}

func TestGetScheduleRunCmd_SuccessWithDirectoryOpt(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	task := &types.ScheduledTask{Id: "test", Command: "echo", CronExpr: "0 0 * * *"}
	ctx.Config.Scheduled = types.ScheduledTasks{
		task,
	}
	cmd := GetScheduleRunCmd(ctx)

	cmd.SetArgs([]string{"--" + flags.WorkingDir, "/app/test"})
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "/app/test", task.Directory)
}

func TestGetScheduleRunCmd_FailedWithTimezoneOpt(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	ctx.Config.Scheduled = types.ScheduledTasks{
		&types.ScheduledTask{Id: "test", Command: "echo", CronExpr: "0 0 * * *"},
	}
	cmd := GetScheduleRunCmd(ctx)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"--" + TimeZone, "Europe/Wrong"})
	err := cmd.Execute()
	assert.Error(t, err)
}
