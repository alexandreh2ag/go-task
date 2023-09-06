package generate

import (
	"bytes"
	"errors"
	"github.com/alexandreh2ag/go-task/context"
	mockOs "github.com/alexandreh2ag/go-task/mocks/os"
	mockAfero "github.com/alexandreh2ag/go-task/mocks/spf13"
	"github.com/alexandreh2ag/go-task/types"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"testing"
)

func TestCheckDir_dirOK(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	outputPath := "/tmp/subdir/output.txt"
	_ = afero.WriteFile(ctx.Fs, outputPath, []byte{}, 0644)
	err := checkDir(ctx, outputPath)
	assert.Equal(t, err, nil)
}

func TestCheckDir_dirNotOK(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	err := checkDir(ctx, "/tmp/anotherdir/output.txt")

	assert.NotEqual(t, err, nil)
}

func TestGenerate_invalidExtension(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	workers := types.WorkerTasks{
		{Id: "test", Command: "fake", User: "test", Directory: "/tmp/dir"},
		{Id: "test2", Command: "ping", User: "test2", Directory: "/tmp/dir"},
	}
	ctx.Config.Workers = workers
	outputPath := "/tmp/subdir/output.txt"
	_ = afero.WriteFile(ctx.Fs, outputPath, []byte{}, 0644)

	err := Generate(ctx, outputPath, "abcd", "myname")

	assert.NotEqual(t, err, nil)
	assert.Contains(t, err.Error(), "Error with unsupported format")
}

func TestGenerate_invalidDir(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	workers := types.WorkerTasks{
		{Id: "test", Command: "fake", User: "test", Directory: "/tmp/dir"},
		{Id: "test2", Command: "ping", User: "test2", Directory: "/tmp/dir"},
	}
	ctx.Config.Workers = workers
	err := Generate(ctx, "/tmp/anotherdir/output.txt", FormatSupervisor, "myname")

	assert.NotEqual(t, err, nil)
	assert.Contains(t, err.Error(), "Error with outputh dir")
}

func TestGenerate_invalidOutputFile(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	workers := types.WorkerTasks{
		{Id: "test", Command: "fake", User: "test", Directory: "/tmp/dir"},
		{Id: "test2", Command: "ping", User: "test2", Directory: "/tmp/dir"},
	}
	ctx.Config.Workers = workers

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fsMock := mockAfero.NewMockFs(ctrl)
	ctx.Fs = fsMock

	outputDir := "/tmp/subdir"
	outputPath := outputDir + "/output.txt"

	dirLogMock := mockOs.NewMockFileInfo(ctrl)
	fsMock.EXPECT().Stat(gomock.Eq(outputDir)).Times(1).Return(dirLogMock, nil)
	fsMock.EXPECT().Create(gomock.Eq(outputPath)).Times(1).Return(nil, errors.New("fail"))

	err := Generate(ctx, outputPath, FormatSupervisor, "myname")

	assert.NotEqual(t, err, nil)
	assert.Contains(t, err.Error(), "Error with output file")
}

func TestGenerate_OK(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	workers := types.WorkerTasks{
		{Id: "test", Command: "fake", User: "test", Directory: "/tmp/dir"},
		{Id: "test2", Command: "ping", User: "test2", Directory: "/tmp/dir"},
	}
	ctx.Config.Workers = workers

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

	err := Generate(ctx, outputPath, FormatSupervisor, "myname")

	assert.Equal(t, nil, err)
}

func TestTemplateSupervisorFile_OK(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	expectedOutput := "[group:test-group]\n" +
		"programs=test,test2\n\n\n" +
		"[program:test]\n" +
		"directory = /tmp/dir\n" +
		"autorestart = true\n" +
		"autostart = true\n" +
		"user = test\n" +
		"command = fake\n\n" +
		"[program:test2]\n" +
		"directory = /tmp/dir\n" +
		"autorestart = true\n" +
		"autostart = true\n" +
		"user = test2\n" +
		"command = ping\n"

	workers := types.WorkerTasks{
		{Id: "test", Command: "fake", User: "test", Directory: "/tmp/dir"},
		{Id: "test2", Command: "ping", User: "test2", Directory: "/tmp/dir"},
	}
	groupName := "test-group"

	ctx.Config.Workers = workers

	buffer := bytes.NewBufferString("")

	err := templateSupervisorFile(ctx, buffer, groupName)
	assert.Equal(t, err, nil)
	assert.Equal(t, expectedOutput, buffer.String())
}

func TestGenerateProgramList(t *testing.T) {
	workers := types.WorkerTasks{
		{Id: "test", Command: "fake"},
		{Id: "test2", Command: "fake"},
	}
	output := generateProgramList(workers)
	assert.Equal(t, output, "test,test2")
}

func TestDeleteFile_OK(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	outputDir := "/tmp/subdir"
	outputPath := outputDir + "/output.txt"
	_ = afero.WriteFile(ctx.Fs, outputPath, []byte{}, 0644)

	err := deleteFile(ctx, outputPath)

	fileExist, _ := afero.Exists(ctx.Fs, outputPath)
	assert.Equal(t, err, nil)
	assert.Equal(t, false, fileExist)
}

func TestDeleteFile_FileDoesNotExist(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	outputDir := "/tmp/subdir"
	outputPath := outputDir + "/output.txt"

	err := deleteFile(ctx, outputPath)

	fileExist, _ := afero.Exists(ctx.Fs, outputPath)
	assert.Equal(t, err, nil)
	assert.Equal(t, false, fileExist)
}

func TestDeleteFile_Error(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	outputDir := "/tmp/subdir"
	outputPath := outputDir + "/output.txt"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fsMock := mockAfero.NewMockFs(ctrl)
	ctx.Fs = fsMock
	dirLogMock := mockOs.NewMockFileInfo(ctrl)
	fsMock.EXPECT().Stat(gomock.Eq(outputDir)).Times(1).Return(dirLogMock, nil)
	fsMock.EXPECT().Remove(gomock.Eq(outputPath)).Times(1).Return(errors.New("fail"))

	err := Generate(ctx, outputPath, FormatSupervisor, "myname")

	assert.NotEqual(t, err, nil)
	assert.Contains(t, err.Error(), "Error when deleting output file")
}
