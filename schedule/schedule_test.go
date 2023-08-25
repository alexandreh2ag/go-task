package schedule

import (
	"alexandreh2ag/go-task/context"
	"alexandreh2ag/go-task/types"
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"os/exec"
	"testing"
	"time"
)

func TestGetCurrentTime(t *testing.T) {

	arizonaTZ, _ := time.LoadLocation("US/Arizona")

	tests := []struct {
		name     string
		now      time.Time
		timezone string
		want     time.Time
		wantErr  bool
	}{
		{
			name:     "SuccessWithoutTimezone",
			timezone: "",
			now:      time.Date(2023, time.January, 25, 15, 4, 13, 0, time.UTC),
			want:     time.Date(2023, time.January, 25, 15, 4, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "SuccessWithoutTimezone",
			timezone: "US/Arizona",
			now:      time.Date(2023, time.January, 25, 15, 4, 13, 0, time.UTC),
			want:     time.Date(2023, time.January, 25, 8, 4, 0, 0, arizonaTZ),
			wantErr:  false,
		},
		{
			name:     "ErrorWithWrongTimezone",
			timezone: "US/Wrong",
			now:      time.Date(2023, time.January, 25, 15, 4, 13, 0, time.UTC),
			want:     time.Date(2023, time.January, 25, 15, 4, 13, 0, time.UTC),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCurrentTime(tt.now, tt.timezone)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("GetCurrentTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	ctx := context.TestContext(io.Discard)

	type args struct {
		scheduledTasks types.ScheduledTasks
		ref            time.Time
	}
	tests := []struct {
		name string
		args args
		want types.ScheduledTasks
	}{
		{
			name: "SuccessWithScheduledTasksEmpty",
			args: args{
				scheduledTasks: types.ScheduledTasks{},
				ref:            time.Date(2023, time.January, 25, 15, 4, 0, 0, time.UTC),
			},
			want: types.ScheduledTasks{},
		},
		{
			name: "SuccessWithScheduledTasks",
			args: args{
				scheduledTasks: types.ScheduledTasks{
					&types.ScheduledTask{Id: "test", Command: "echo test", CronExpr: "* * * * *", Logger: ctx.Logger},
					&types.ScheduledTask{Id: "test2", Command: "echo test", CronExpr: "0 0 * * *", Logger: ctx.Logger},
					&types.ScheduledTask{Id: "test3", Command: "wrong", CronExpr: "* * * * *", Logger: ctx.Logger},
					&types.ScheduledTask{Id: "test4", Command: "wrong", CronExpr: "wrong", Logger: ctx.Logger},
				},
				ref: time.Date(2023, time.January, 25, 15, 4, 0, 0, time.UTC),
			},
			want: types.ScheduledTasks{
				&types.ScheduledTask{
					Id:       "test",
					Command:  "echo test",
					CronExpr: "* * * * *",
					Logger:   ctx.Logger,
					TaskResult: &types.TaskResult{
						Status:   types.Succeed,
						Output:   *bytes.NewBuffer([]byte("test\n")),
						StartAt:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
						FinishAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
				},
				&types.ScheduledTask{
					Id:       "test2",
					Command:  "echo test",
					CronExpr: "0 0 * * *",
					Logger:   ctx.Logger,
				},
				&types.ScheduledTask{
					Id:       "test3",
					Command:  "wrong",
					CronExpr: "* * * * *",
					Logger:   ctx.Logger,
					TaskResult: &types.TaskResult{
						Status:   types.Failed,
						Error:    &exec.Error{Name: "wrong", Err: errors.New("executable file not found in $PATH")},
						Output:   *bytes.NewBuffer(nil),
						StartAt:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
						FinishAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
				},
				&types.ScheduledTask{
					Id:       "test4",
					Command:  "wrong",
					CronExpr: "wrong",
					Logger:   ctx.Logger,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx.Config.Scheduled = tt.args.scheduledTasks
			Run(ctx, tt.args.ref)

			for _, task := range tt.args.scheduledTasks {
				if task.TaskResult != nil {
					assert.WithinDuration(t, task.TaskResult.StartAt, task.TaskResult.FinishAt, time.Second)
					assert.NotNil(t, task.TaskResult.Task)
					task.TaskResult.Task = nil
					task.TaskResult.StartAt = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
					task.TaskResult.FinishAt = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
				}
			}

			assert.Equal(t, tt.want, tt.args.scheduledTasks)
		})
	}
}
