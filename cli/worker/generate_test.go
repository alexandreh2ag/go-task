package worker

import (
	"github.com/alexandreh2ag/go-task/cli/flags"
	"github.com/alexandreh2ag/go-task/context"
	mockOs "github.com/alexandreh2ag/go-task/mocks/os"
	mockAfero "github.com/alexandreh2ag/go-task/mocks/spf13"
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
