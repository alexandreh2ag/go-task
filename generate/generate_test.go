package generate

import (
	"alexandreh2ag/go-task/context"
	mockOs "alexandreh2ag/go-task/mocks/os"
	mockAfero "alexandreh2ag/go-task/mocks/spf13"
	"errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"strings"
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
	outputPath := "/tmp/subdir/output.txt"
	_ = afero.WriteFile(ctx.Fs, outputPath, []byte{}, 0644)
	err := Generate(ctx, outputPath, "abcd", "myname")

	assert.NotEqual(t, err, nil)
	assert.Equal(t, true, strings.Contains(err.Error(), "Error with unsupported format"))
}

func TestGenerate_invalidDir(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	err := Generate(ctx, "/tmp/anotherdir/output.txt", FormatSupervisor, "myname")

	assert.NotEqual(t, err, nil)
	assert.Equal(t, true, strings.Contains(err.Error(), "Error with outputh dir"))
}

func TestGenerate_invalidOutputFile(t *testing.T) {
	ctx := context.TestContext(io.Discard)
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
	assert.Equal(t, true, strings.Contains(err.Error(), "Error with output file"))
}

func TestGenerate_OK(t *testing.T) {
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
	fileMock.EXPECT().Write(gomock.Any()).Times(5).Return(1, nil)

	err := Generate(ctx, outputPath, FormatSupervisor, "myname")

	assert.Equal(t, nil, err)
}

func TestTemplateSupervisorFile_OK(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fsMock := mockAfero.NewMockFs(ctrl)
	fileMock := mockAfero.NewMockFile(ctrl)
	ctx.Fs = fsMock

	fileMock.EXPECT().Write(gomock.Any()).Times(5).Return(1, nil)

	groupName := "test-group"

	err := templateSupervisorFile(ctx, fileMock, groupName)
	assert.Equal(t, err, nil)
}
