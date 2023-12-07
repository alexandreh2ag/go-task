package schedule

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/alexandreh2ag/go-task/context"
	mockOs "github.com/alexandreh2ag/go-task/mocks/os"
	mockAfero "github.com/alexandreh2ag/go-task/mocks/spf13"
	"github.com/alexandreh2ag/go-task/types"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"log/slog"
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
		noResultPrint  bool
		resultPath     string
		taskFilter     []string
		force          bool
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(ctrl *gomock.Controller, fs *mockAfero.MockFs)
		want     []*types.TaskResult
	}{
		{
			name: "SuccessWithScheduledTasksEmpty",
			args: args{
				scheduledTasks: types.ScheduledTasks{},
				ref:            time.Date(2023, time.January, 25, 15, 4, 0, 0, time.UTC),
				noResultPrint:  true,
				resultPath:     "",
			},
			mockFunc: func(ctrl *gomock.Controller, fs *mockAfero.MockFs) {},
			want:     []*types.TaskResult{},
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
				ref:           time.Date(2023, time.January, 25, 15, 4, 0, 0, time.UTC),
				noResultPrint: true,
				resultPath:    "",
			},
			mockFunc: func(ctrl *gomock.Controller, fs *mockAfero.MockFs) {},
			want: []*types.TaskResult{
				{
					Status:   types.Succeed,
					Output:   *bytes.NewBuffer([]byte("test\n")),
					StartAt:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
					FinishAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Status:   types.Failed,
					Error:    &exec.Error{Name: "wrong", Err: errors.New("executable file not found in $PATH")},
					Output:   *bytes.NewBuffer(nil),
					StartAt:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
					FinishAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "SuccessWithScheduledTasksAndTaskFilter",
			args: args{
				scheduledTasks: types.ScheduledTasks{
					&types.ScheduledTask{Id: "test", Command: "echo test", CronExpr: "* * * * *", Logger: ctx.Logger},
					&types.ScheduledTask{Id: "test2", Command: "echo test", CronExpr: "* * * * *", Logger: ctx.Logger},
				},
				ref:           time.Date(2023, time.January, 25, 15, 4, 0, 0, time.UTC),
				taskFilter:    []string{"test"},
				noResultPrint: true,
				resultPath:    "",
			},
			mockFunc: func(ctrl *gomock.Controller, fs *mockAfero.MockFs) {},
			want: []*types.TaskResult{
				{
					Status:   types.Succeed,
					Output:   *bytes.NewBuffer([]byte("test\n")),
					StartAt:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
					FinishAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "SuccessWithScheduledTasksAndTaskFilterAndForceFlag",
			args: args{
				scheduledTasks: types.ScheduledTasks{
					&types.ScheduledTask{Id: "test", Command: "echo test", CronExpr: "0 0 * * *", Logger: ctx.Logger},
					&types.ScheduledTask{Id: "test2", Command: "echo test2", CronExpr: "0 0 * * *", Logger: ctx.Logger},
				},
				ref:           time.Date(2023, time.January, 25, 15, 4, 0, 0, time.UTC),
				taskFilter:    []string{"test"},
				force:         true,
				noResultPrint: true,
				resultPath:    "",
			},
			mockFunc: func(ctrl *gomock.Controller, fs *mockAfero.MockFs) {},
			want: []*types.TaskResult{
				{
					Status:   types.Succeed,
					Output:   *bytes.NewBuffer([]byte("test\n")),
					StartAt:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
					FinishAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "SuccessWithScheduledTasksAndForceFlag",
			args: args{
				scheduledTasks: types.ScheduledTasks{
					&types.ScheduledTask{Id: "test", Command: "echo test", CronExpr: "0 0 * * *", Logger: ctx.Logger},
					&types.ScheduledTask{Id: "test2", Command: "echo test2", CronExpr: "0 0 * * *", Logger: ctx.Logger},
				},
				ref:           time.Date(2023, time.January, 25, 15, 4, 0, 0, time.UTC),
				force:         true,
				noResultPrint: true,
				resultPath:    "",
			},
			mockFunc: func(ctrl *gomock.Controller, fs *mockAfero.MockFs) {},
			want: []*types.TaskResult{
				{
					Status:   types.Succeed,
					Output:   *bytes.NewBuffer([]byte("test\n")),
					StartAt:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
					FinishAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Status:   types.Succeed,
					Output:   *bytes.NewBuffer([]byte("test2\n")),
					StartAt:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
					FinishAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "SuccessWithScheduledTasksAndPrintResult",
			args: args{
				scheduledTasks: types.ScheduledTasks{
					&types.ScheduledTask{Id: "test", Command: "echo test", CronExpr: "* * * * *", Logger: ctx.Logger},
				},
				ref:           time.Date(2023, time.January, 25, 15, 4, 0, 0, time.UTC),
				noResultPrint: false,
				resultPath:    "",
			},
			mockFunc: func(ctrl *gomock.Controller, mockFs *mockAfero.MockFs) {},
			want: []*types.TaskResult{
				{
					Status:   types.Succeed,
					Output:   *bytes.NewBuffer([]byte("test\n")),
					StartAt:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
					FinishAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "SuccessWithScheduledTasksAndWriteToLog",
			args: args{
				scheduledTasks: types.ScheduledTasks{
					&types.ScheduledTask{Id: "test", Command: "echo test", CronExpr: "* * * * *", Logger: ctx.Logger},
				},
				ref:           time.Date(2023, time.January, 25, 15, 4, 0, 0, time.UTC),
				noResultPrint: true,
				resultPath:    "/var/log/gtask.log",
			},
			mockFunc: func(ctrl *gomock.Controller, mockFs *mockAfero.MockFs) {
				dirLogMock := mockOs.NewMockFileInfo(ctrl)
				fileLogMock := mockAfero.NewMockFile(ctrl)
				dirLogMock.EXPECT().IsDir().Times(1).Return(true)
				fileLogMock.EXPECT().WriteString(gomock.Any()).Times(1).Return(1, nil)

				mockFs.EXPECT().Stat(gomock.Eq("/var/log")).Times(1).Return(dirLogMock, nil)
				mockFs.EXPECT().OpenFile(gomock.Eq("/var/log/gtask.log"), gomock.Any(), gomock.Any()).Times(1).Return(fileLogMock, nil)
			},
			want: []*types.TaskResult{
				{
					Status:   types.Succeed,
					Output:   *bytes.NewBuffer([]byte("test\n")),
					StartAt:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
					FinishAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "ErrorWithScheduledTasksAndWriteToLogFailed",
			args: args{
				scheduledTasks: types.ScheduledTasks{
					&types.ScheduledTask{Id: "test", Command: "echo test", CronExpr: "* * * * *", Logger: ctx.Logger},
				},
				ref:           time.Date(2023, time.January, 25, 15, 4, 0, 0, time.UTC),
				noResultPrint: true,
				resultPath:    "/var/log/gtask.log",
			},
			mockFunc: func(ctrl *gomock.Controller, mockFs *mockAfero.MockFs) {
				dirLogMock := mockOs.NewMockFileInfo(ctrl)
				dirLogMock.EXPECT().IsDir().Times(1).Return(false)
				mockFs.EXPECT().Stat(gomock.Eq("/var/log")).Times(1).Return(dirLogMock, nil)
				mockFs.EXPECT().MkdirAll(gomock.Eq("/var/log"), gomock.Any()).Times(1).Return(errors.New("fail"))
			},
			want: []*types.TaskResult{
				{
					Status:   types.Succeed,
					Output:   *bytes.NewBuffer([]byte("test\n")),
					StartAt:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
					FinishAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fsMock := mockAfero.NewMockFs(ctrl)
			ctx.Config.Scheduled = tt.args.scheduledTasks
			tt.mockFunc(ctrl, fsMock)
			ctx.Fs = fsMock
			got := Run(ctx, tt.args.ref, tt.args.taskFilter, tt.args.force, tt.args.noResultPrint, tt.args.resultPath)

			for _, result := range got {
				assert.WithinDuration(t, result.StartAt, result.FinishAt, time.Second)
				assert.NotNil(t, result.Task)
				result.Task = nil
				result.StartAt = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
				result.FinishAt = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

			}

			assert.Equal(t, len(tt.want), len(got))
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestRun_SuccessWithOneTaskResult(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	scheduledTasks := types.ScheduledTasks{
		&types.ScheduledTask{Id: "test", Command: "echo test", CronExpr: "0 0 * * *", Logger: ctx.Logger, LatestTaskResult: &types.TaskResult{Status: types.Pending}},
		&types.ScheduledTask{Id: "test2", Command: "echo test", CronExpr: "0 0 * * *", Logger: ctx.Logger, LatestTaskResult: &types.TaskResult{Status: types.Succeed}},
	}
	want := types.ScheduledTasks{
		&types.ScheduledTask{Id: "test", Command: "echo test", CronExpr: "0 0 * * *", Logger: ctx.Logger, LatestTaskResult: &types.TaskResult{Status: types.Pending}},
		&types.ScheduledTask{Id: "test2", Command: "echo test", CronExpr: "0 0 * * *", Logger: ctx.Logger, LatestTaskResult: nil},
	}
	ref := time.Date(2023, time.January, 25, 15, 4, 0, 0, time.UTC)

	ctx.Config.Scheduled = scheduledTasks
	_ = Run(ctx, ref, []string{}, false, true, "")

	assert.ElementsMatch(t, want, scheduledTasks)
}

func TestFormatTaskResult(t *testing.T) {

	tests := []struct {
		name   string
		result *types.TaskResult
		want   string
	}{
		{
			name: "SuccessWithSucceedResult",
			result: &types.TaskResult{
				Status:   types.Succeed,
				Output:   *bytes.NewBuffer([]byte("my output\n")),
				Task:     &types.ScheduledTask{Id: "test"},
				StartAt:  time.Date(1970, time.January, 1, 0, 30, 0, 0, time.UTC),
				FinishAt: time.Date(1970, time.January, 1, 0, 35, 1, 0, time.UTC),
			},
			want: fmt.Sprintf(
				"%s\nTask test finish with status 'succeed'\nStart at 1970-01-01T00:30:00, finish at 1970-01-01T00:35:01 (5m1s)\noutput:\nmy output\n\n%s\n",
				BlocSeparator,
				BlocSeparator,
			),
		},
		{
			name: "SuccessWithFailedResultAndEmptyOutput",
			result: &types.TaskResult{
				Status:   types.Failed,
				Error:    errors.New("critical error"),
				Output:   *bytes.NewBuffer([]byte("")),
				Task:     &types.ScheduledTask{Id: "test"},
				StartAt:  time.Date(1970, time.January, 1, 0, 30, 0, 0, time.UTC),
				FinishAt: time.Date(1970, time.January, 1, 0, 35, 1, 0, time.UTC),
			},
			want: fmt.Sprintf(
				"%s\nTask test finish with status 'failed'\nStart at 1970-01-01T00:30:00, finish at 1970-01-01T00:35:01 (5m1s)\nDue to the following error: critical error\n%s\n",
				BlocSeparator,
				BlocSeparator,
			),
		},
		{
			name: "SuccessWithFailedResultAndOutput",
			result: &types.TaskResult{
				Status:   types.Failed,
				Error:    errors.New("critical error"),
				Output:   *bytes.NewBuffer([]byte("my output")),
				Task:     &types.ScheduledTask{Id: "test"},
				StartAt:  time.Date(1970, time.January, 1, 0, 30, 0, 0, time.UTC),
				FinishAt: time.Date(1970, time.January, 1, 0, 35, 1, 0, time.UTC),
			},
			want: fmt.Sprintf(
				"%s\nTask test finish with status 'failed'\nStart at 1970-01-01T00:30:00, finish at 1970-01-01T00:35:01 (5m1s)\noutput:\nmy output\nDue to the following error: critical error\n%s\n",
				BlocSeparator,
				BlocSeparator,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FormatTaskResult(tt.result), "FormatTaskResult(%v)", tt.result)
		})
	}
}

func TestWriteToLogFile(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	type args struct {
		resultPath string
		data       string
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(ctrl *gomock.Controller, mockFs *mockAfero.MockFs)
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "SuccessWhenCreateDir",
			args: args{
				resultPath: "/var/log/gtask.log",
				data:       "my data",
			},
			mockFunc: func(ctrl *gomock.Controller, mockFs *mockAfero.MockFs) {
				dirLogMock := mockOs.NewMockFileInfo(ctrl)
				fileLogMock := mockAfero.NewMockFile(ctrl)
				dirLogMock.EXPECT().IsDir().Times(1).Return(false)
				fileLogMock.EXPECT().WriteString(gomock.Any()).Times(1).Return(1, nil)

				mockFs.EXPECT().Stat(gomock.Eq("/var/log")).Times(1).Return(dirLogMock, nil)
				mockFs.EXPECT().MkdirAll(gomock.Eq("/var/log"), gomock.Any()).Times(1).Return(nil)
				mockFs.EXPECT().OpenFile(gomock.Eq("/var/log/gtask.log"), gomock.Any(), gomock.Any()).Times(1).Return(fileLogMock, nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "SuccessWhenFileExist",
			args: args{
				resultPath: "/var/log/gtask.log",
				data:       "my data",
			},
			mockFunc: func(ctrl *gomock.Controller, mockFs *mockAfero.MockFs) {
				dirLogMock := mockOs.NewMockFileInfo(ctrl)
				fileLogMock := mockAfero.NewMockFile(ctrl)
				dirLogMock.EXPECT().IsDir().Times(1).Return(true)
				fileLogMock.EXPECT().WriteString(gomock.Any()).Times(1).Return(1, nil)

				mockFs.EXPECT().Stat(gomock.Eq("/var/log")).Times(1).Return(dirLogMock, nil)
				mockFs.EXPECT().OpenFile(gomock.Eq("/var/log/gtask.log"), gomock.Any(), gomock.Any()).Times(1).Return(fileLogMock, nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "ErrorWhenDirCreateFailed",
			args: args{
				resultPath: "/var/log/gtask.log",
				data:       "my data",
			},
			mockFunc: func(ctrl *gomock.Controller, mockFs *mockAfero.MockFs) {
				dirLogMock := mockOs.NewMockFileInfo(ctrl)
				dirLogMock.EXPECT().IsDir().Times(1).Return(false)

				mockFs.EXPECT().Stat(gomock.Eq("/var/log")).Times(1).Return(dirLogMock, nil)
				mockFs.EXPECT().MkdirAll(gomock.Eq("/var/log"), gomock.Any()).Times(1).Return(errors.New("fail"))
			},
			wantErr: assert.Error,
		},
		{
			name: "ErrorWhenOpenFileFailed",
			args: args{
				resultPath: "/var/log/gtask.log",
				data:       "my data",
			},
			mockFunc: func(ctrl *gomock.Controller, mockFs *mockAfero.MockFs) {
				dirLogMock := mockOs.NewMockFileInfo(ctrl)
				dirLogMock.EXPECT().IsDir().Times(1).Return(true)

				mockFs.EXPECT().Stat(gomock.Eq("/var/log")).Times(1).Return(dirLogMock, nil)
				mockFs.EXPECT().OpenFile(gomock.Eq("/var/log/gtask.log"), gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.New("fail"))
			},
			wantErr: assert.Error,
		},
		{
			name: "ErrorWhenWriteFileFailed",
			args: args{
				resultPath: "/var/log/gtask.log",
				data:       "my data",
			},
			mockFunc: func(ctrl *gomock.Controller, mockFs *mockAfero.MockFs) {
				dirLogMock := mockOs.NewMockFileInfo(ctrl)
				fileLogMock := mockAfero.NewMockFile(ctrl)
				dirLogMock.EXPECT().IsDir().Times(1).Return(true)
				fileLogMock.EXPECT().WriteString(gomock.Any()).Times(1).Return(0, errors.New("fail"))

				mockFs.EXPECT().Stat(gomock.Eq("/var/log")).Times(1).Return(dirLogMock, nil)
				mockFs.EXPECT().OpenFile(gomock.Eq("/var/log/gtask.log"), gomock.Any(), gomock.Any()).Times(1).Return(fileLogMock, nil)
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fsMock := mockAfero.NewMockFs(ctrl)
			tt.mockFunc(ctrl, fsMock)
			ctx.Fs = fsMock
			err := WriteToLogFile(ctx, tt.args.resultPath, tt.args.data)
			tt.wantErr(t, err, fmt.Sprintf("WriteToLogFile(ctx, %v, %v)", tt.args.resultPath, tt.args.data))
		})
	}
}

func TestStart_SuccessFirstTick(t *testing.T) {
	tickUnit = time.Millisecond

	b := bytes.NewBufferString("")
	ctx := context.TestContext(b)
	ctx.LogLevel.Set(slog.LevelDebug)
	fakeClock := clockwork.NewFakeClockAt(time.Date(1970, time.January, 1, 0, 0, 59, 0, time.UTC))
	ctx.Clock = fakeClock

	go func() {
		err := Start(ctx, 1, "Europe/Paris", []string{}, true, "")
		assert.NoError(t, err)
	}()

	fakeClock.BlockUntil(1)
	fakeClock.Advance(1 * time.Second)
	ctx.Cancel()
	assert.Contains(t, b.String(), "msg=\"first tick\"")
	assert.NotContains(t, b.String(), "msg=tick")
}

func TestStart_SuccessFirstTickAndTick(t *testing.T) {
	tickUnit = time.Millisecond

	b := bytes.NewBufferString("")
	ctx := context.TestContext(b)
	ctx.LogLevel.Set(slog.LevelDebug)
	fakeClock := clockwork.NewFakeClockAt(time.Date(1970, time.January, 1, 0, 0, 59, 0, time.UTC))
	ctx.Clock = fakeClock

	go func() {
		err := Start(ctx, 1, "Europe/Paris", []string{}, true, "")
		assert.NoError(t, err)
	}()

	fakeClock.BlockUntil(1)
	fakeClock.Advance(1 * time.Second)
	time.Sleep(50 * time.Millisecond)
	fakeClock.Advance(1 * time.Second)
	ctx.Cancel()
	assert.Contains(t, b.String(), "msg=\"first tick\"")
	assert.Contains(t, b.String(), "msg=tick")
}

func TestStart_ErrorWhenNextTickAfterFailed(t *testing.T) {
	tickUnit = time.Millisecond

	b := bytes.NewBufferString("")
	ctx := context.TestContext(b)
	ctx.LogLevel.Set(slog.LevelDebug)
	fakeClock := clockwork.NewFakeClockAt(time.Date(1970, time.January, 1, 0, 0, 59, 0, time.UTC))
	ctx.Clock = fakeClock

	err := Start(ctx, 0, "Europe/Paris", []string{}, true, "")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "could not calculate next tick of expr */0 * * * *: tried so hard")
}

func TestStart_ErrorWhenGetCurrentTimeFailed(t *testing.T) {
	tickUnit = time.Millisecond

	b := bytes.NewBufferString("")
	ctx := context.TestContext(b)
	ctx.LogLevel.Set(slog.LevelDebug)
	fakeClock := clockwork.NewFakeClockAt(time.Date(1970, time.January, 1, 0, 0, 59, 0, time.UTC))
	ctx.Clock = fakeClock

	go func() {
		err := Start(ctx, 1, "Europe/Wrong", []string{}, true, "")
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "unknown time zone Europe/Wrong")
	}()

	fakeClock.BlockUntil(1)
	fakeClock.Advance(1 * time.Second)
	time.Sleep(50 * time.Millisecond)
}
