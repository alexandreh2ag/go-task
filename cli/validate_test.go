package cli

import (
	"bytes"
	"fmt"
	"github.com/alexandreh2ag/go-task/context"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetValidateRunFn_Success(t *testing.T) {
	ctx := context.TestContext(nil)
	cmd := GetValidateCmd(ctx)
	fsFake := afero.NewMemMapFs()
	viper.Reset()
	viper.SetFs(fsFake)
	path := "/app"
	_ = fsFake.Mkdir(path, 0775)
	_ = afero.WriteFile(
		fsFake,
		fmt.Sprintf("%s/tasks.yml", path),
		[]byte("workers:\n- {id: 'test',command: 'fake'}\nscheduled:\n- {id: 'test',command: 'fake',expr: '0 0 * * *'}\n"),
		0644,
	)

	viper.Set(Config, fmt.Sprintf("%s/tasks.yml", path))
	_ = cmd.Execute()
	err := GetValidateRunFn(ctx)(cmd, []string{})
	assert.NoError(t, err)
}

func Test_GetValidateRunFn_ErrorValidate(t *testing.T) {
	b := bytes.NewBufferString("")
	ctx := context.TestContext(b)
	cmd := GetValidateCmd(ctx)
	fsFake := afero.NewMemMapFs()
	viper.Reset()
	viper.SetFs(fsFake)
	path := "/app"
	_ = fsFake.Mkdir(path, 0775)
	_ = afero.WriteFile(
		fsFake,
		fmt.Sprintf("%s/tasks.yml", path),
		[]byte("workers:\n- {id: 'test'}\n"),
		0644,
	)

	viper.Set(Config, fmt.Sprintf("%s/tasks.yml", path))

	_ = cmd.Execute()
	err := GetValidateRunFn(ctx)(cmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, b.String(), "Key: 'Config.Workers[0].Command' Error:Field validation for 'Command' failed on the 'required' tag")
}
