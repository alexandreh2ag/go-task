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
	"os"
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
		{
			Id:        "test",
			Command:   "fake",
			User:      "toto",
			Directory: "/tmp/dir",
		},
		{
			Id:        "test2",
			Command:   "fake",
			User:      "toto",
			Directory: "/tmp/dir",
		},
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
		{
			Id:        "test",
			Command:   "fake",
			User:      "toto",
			Directory: "/tmp/dir",
		},
		{
			Id:        "test2",
			Command:   "fake",
			User:      "toto",
			Directory: "/tmp/dir",
		},
	}
	ctx.Config.Workers = workers
	err := Generate(ctx, "/tmp/anotherdir/output.txt", FormatSupervisor, "myname")

	assert.NotEqual(t, err, nil)
	assert.Contains(t, err.Error(), "Error with outputh dir")
}

func TestGenerate_invalidOutputFile(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	workers := types.WorkerTasks{
		{
			Id:        "test",
			Command:   "fake",
			User:      "toto",
			Directory: "/tmp/dir",
		},
		{
			Id:        "test2",
			Command:   "fake",
			User:      "toto",
			Directory: "/tmp/dir",
		},
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
		{
			Id:        "test",
			Command:   "fake",
			User:      "toto",
			Directory: "/tmp/dir",
		},
		{
			Id:        "test2",
			Command:   "fake",
			User:      "toto",
			Directory: "/tmp/dir",
		},
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
	groupName := "test-group"
	workers := types.WorkerTasks{
		{
			Id:        "test",
			Command:   "fake",
			GroupName: groupName,
			User:      "toto",
			Directory: "/tmp/dir",
			Envs: map[string]string{
				"FOO": "BAR",
			},
		},
		{
			Id:        "test2",
			Command:   "fake",
			GroupName: groupName,
			User:      "toto",
			Directory: "/tmp/dir",
			Envs: map[string]string{
				"BAR": "FOO",
			},
			Expression: "BAR == FOO",
		},
		{
			Id:        "test3",
			Command:   "fake",
			GroupName: groupName,
			User:      "toto",
			Directory: "/tmp/dir",
			Envs: map[string]string{
				"BAR": "NOT_FOO",
			},
			Expression: "BAR == FOO",
		},
	}

	expectedOutput := "[group:test-group]\n" +
		"programs=test-group-test,test-group-test2\n\n\n" +
		"[program:test-group-test]\n" +
		"directory = /tmp/dir\n" +
		"autorestart = true\n" +
		"autostart = true\n" +
		"user = toto\n" +
		"command = fake\n" +
		"environment = FOO=\"BAR\"\n\n" +
		"[program:test-group-test2]\n" +
		"directory = /tmp/dir\n" +
		"autorestart = true\n" +
		"autostart = true\n" +
		"user = toto\n" +
		"command = fake\n" +
		"environment = BAR=\"FOO\"\n"

	ctx.Config.Workers = workers

	buffer := bytes.NewBufferString("")

	err := templateSupervisorFile(ctx, buffer, groupName)
	assert.Equal(t, err, nil)

	assert.Equal(t, expectedOutput, buffer.String())
}

func TestTemplateSupervisorFile_Eval_Fail(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	groupName := "test-group"
	workers := types.WorkerTasks{
		{
			Id:        "test",
			Command:   "fake",
			GroupName: groupName,
			User:      "toto",
			Directory: "/tmp/dir",
			Envs: map[string]string{
				"FOO": "BAR",
			},
			Expression: "BAR & FOO",
		},
	}

	ctx.Config.Workers = workers

	buffer := bytes.NewBufferString("")

	err := templateSupervisorFile(ctx, buffer, groupName)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "can't evaluate expression for task 'test'")
}

func TestGenerateProgramList(t *testing.T) {
	prefix := "pref"
	workers := types.WorkerTasks{
		{
			Id:        "test",
			Command:   "fake",
			User:      "toto",
			Directory: "/tmp/dir",
			GroupName: prefix,
		},
		{
			Id:        "test2",
			Command:   "fake",
			User:      "toto",
			Directory: "/tmp/dir",
			GroupName: prefix,
		},
	}

	output := generateProgramList(workers)
	assert.Equal(t, output, prefix+"-test,"+prefix+"-test2")
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

func TestGenerate_ErrorDeleteFile(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	outputDir := "/tmp/subdir"
	outputPath := outputDir + "/output.txt"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fsMock := mockAfero.NewMockFs(ctrl)
	ctx.Fs = fsMock
	outputDirMock := mockOs.NewMockFileInfo(ctrl)
	outputMock := mockOs.NewMockFileInfo(ctrl)
	fsMock.EXPECT().Stat(gomock.Eq(outputDir)).Times(1).Return(outputDirMock, nil)
	fsMock.EXPECT().Stat(gomock.Eq(outputPath)).Times(1).Return(outputMock, nil)
	fsMock.EXPECT().Remove(gomock.Eq(outputPath)).Times(1).Return(errors.New("fail"))

	err := Generate(ctx, outputPath, FormatSupervisor, "myname")

	assert.NotEqual(t, err, nil)
	assert.Contains(t, err.Error(), "Error when deleting output file")
}

func TestGenerate_NoErrorNoWorkerAndNoOutputFile(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	outputDir := "/tmp/subdir"
	outputPath := outputDir + "/output.txt"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fsMock := mockAfero.NewMockFs(ctrl)
	ctx.Fs = fsMock
	outputDirMock := mockOs.NewMockFileInfo(ctrl)
	outputMock := mockOs.NewMockFileInfo(ctrl)
	fsMock.EXPECT().Stat(gomock.Eq(outputDir)).Times(1).Return(outputDirMock, nil)
	fsMock.EXPECT().Stat(gomock.Eq(outputPath)).Times(1).Return(outputMock, errors.New("fail not found"))

	err := Generate(ctx, outputPath, FormatSupervisor, "myname")

	assert.NoError(t, err)
}

func TestGenerate_NoErrorDeleteFile(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	outputDir := "/tmp/subdir"
	outputPath := outputDir + "/output.txt"

	_ = afero.WriteFile(ctx.Fs, outputPath, []byte{}, 0644)

	err := Generate(ctx, outputPath, FormatSupervisor, "myname")

	fileExist, _ := afero.Exists(ctx.Fs, outputPath)
	assert.Equal(t, err, nil)
	assert.Equal(t, false, fileExist)
}

func TestGenerateEnvVars(t *testing.T) {
	groupName := "group"

	worker := types.WorkerTask{
		Id:        "test2",
		Command:   "fake",
		User:      "toto",
		GroupName: groupName,
		Directory: "/tmp/dir",
	}

	tests := []struct {
		name      string
		groupName string
		argsEnv   map[string]string
		worker    types.WorkerTask
		want      string
	}{
		{
			name:   "expanded vars",
			worker: worker,
			argsEnv: map[string]string{
				"B_SECOND": "foo",
				"A_FIRST":  "bar",
			},
			want: "A_FIRST=\"bar\",B_SECOND=\"foo\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			worker.Envs = tt.argsEnv
			output := generateEnvVars(worker)
			assert.Equal(t, output, tt.want)
			worker.Envs = nil
		})
	}

	_ = os.Unsetenv("MY_VAR")

}
