package worker

import (
	"github.com/alexandreh2ag/go-task/cli/flags"
	"github.com/alexandreh2ag/go-task/context"
	mockOs "github.com/alexandreh2ag/go-task/mocks/os"
	mockAfero "github.com/alexandreh2ag/go-task/mocks/spf13"
	"github.com/alexandreh2ag/go-task/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"strings"
	"testing"
)

func TestGetWorkerGenerateCmd_Success(t *testing.T) {
	ctx := context.TestContext(io.Discard)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fsMock := mockAfero.NewMockFs(ctrl)
	fileMock := mockAfero.NewMockFile(ctrl)
	ctx.Fs = fsMock

	outputDir := "/tmp/subdir"
	outputPath := outputDir + "/output.txt"

	workers := types.WorkerTasks{
		{Id: "test", Command: "fake", User: "test", Directory: "/tmp/dir"},
		{Id: "test2", Command: "ping", User: "test2", Directory: "/tmp/dir"},
	}
	ctx.Config.Workers = workers

	dirLogMock := mockOs.NewMockFileInfo(ctrl)
	fsMock.EXPECT().Stat(gomock.Eq(outputDir)).Times(1).Return(dirLogMock, nil)
	fsMock.EXPECT().Create(gomock.Eq(outputPath)).Times(1).Return(fileMock, nil)
	fileMock.EXPECT().Write(gomock.Any()).AnyTimes().Return(1, nil)

	cmd := GetWorkerGenerateCmd(ctx)
	cmd.SetArgs([]string{"--" + flags.GroupName, "test", "--" + OutputPath, outputPath})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestGetWorkerGenerateCmd_MissingArgs(t *testing.T) {
	ctx := context.TestContext(io.Discard)

	cmd := GetWorkerGenerateCmd(ctx)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	err := cmd.Execute()
	assert.NotEqual(t, err, nil)
	assert.Equal(t, true, strings.Contains(err.Error(), "missing mandatory arguments"))
}

func TestFormatEnvVars(t *testing.T) {
	tests := []struct {
		name       string
		argStrings []string
		want       map[string]string
	}{
		{
			name:       "Success",
			want:       map[string]string{"key1": "value1", "key2": "value2"},
			argStrings: []string{"key1=value1", "key2=value2"},
		},
		{
			name:       "SuccessEmpty",
			want:       map[string]string{},
			argStrings: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatEnvVars(tt.argStrings)
			assert.Equal(t, tt.want, got)
		})
	}
}
