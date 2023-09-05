package cli

import (
	"fmt"
	"github.com/alexandreh2ag/go-task/config"
	"github.com/alexandreh2ag/go-task/context"
	"github.com/alexandreh2ag/go-task/types"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_initConfig_ConfigFromFlagPath(t *testing.T) {
	ctx := context.TestContext(nil)
	cmd := GetRootCmd(ctx)
	fsFake := afero.NewMemMapFs()
	viper.Reset()
	viper.SetFs(fsFake)
	path := "/app"
	_ = fsFake.Mkdir(path, 0775)
	_ = afero.WriteFile(fsFake, fmt.Sprintf("%s/tasks.yml", path), []byte("workers:\n- {id: 'test',command: 'fake'}\nscheduled:\n- {id: 'test',command: 'fake',expr: '0 0 * * *'}\n"), 0644)
	viper.Set(Config, fmt.Sprintf("%s/tasks.yml", path))
	initConfig(ctx, cmd)
	want := &config.Config{
		Workers: types.WorkerTasks{
			{Id: "test", Command: "fake"},
		},
		Scheduled: types.ScheduledTasks{
			{Id: "test", Command: "fake", CronExpr: "0 0 * * *"},
		},
	}
	assert.Equal(t, want, ctx.Config)
}

func Test_initConfig_ErrorWithoutConfig(t *testing.T) {
	ctx := context.TestContext(nil)
	cmd := GetRootCmd(ctx)
	fsFake := afero.NewMemMapFs()
	viper.Reset()
	viper.SetFs(fsFake)

	want := &config.Config{
		Workers:   types.WorkerTasks{},
		Scheduled: types.ScheduledTasks{},
	}
	defer func() {
		if r := recover(); r != nil {

			assert.Equal(t, want, ctx.Config)
		} else {
			t.Errorf("initConfig should have panicked")
		}
	}()
	initConfig(ctx, cmd)
}

func Test_initConfig_ErrorLoadConfig(t *testing.T) {
	ctx := context.TestContext(nil)
	cmd := GetRootCmd(ctx)
	fsFake := afero.NewMemMapFs()
	viper.Reset()
	viper.SetFs(fsFake)
	path := "/app"
	_ = fsFake.Mkdir(path, 0775)
	_ = afero.WriteFile(fsFake, fmt.Sprintf("%s/tasks.yml", path), []byte("workers{]"), 0644)
	viper.Set(Config, fmt.Sprintf("%s/tasks.yml", path))
	want := &config.Config{
		Workers:   types.WorkerTasks{},
		Scheduled: types.ScheduledTasks{},
	}
	defer func() {
		if r := recover(); r != nil {

			assert.Equal(t, want, ctx.Config)
		} else {
			t.Errorf("initConfig should have panicked")
		}
	}()

	initConfig(ctx, cmd)
}

func Test_initConfig_ErrorUnmarshalConfig(t *testing.T) {
	ctx := context.TestContext(nil)
	cmd := GetRootCmd(ctx)
	fsFake := afero.NewMemMapFs()
	viper.Reset()
	viper.SetFs(fsFake)
	path := "/app"
	_ = fsFake.Mkdir(path, 0775)
	_ = afero.WriteFile(fsFake, fmt.Sprintf("%s/tasks.yml", path), []byte("workers:\n- {id: ['value']}"), 0644)
	viper.Set(Config, fmt.Sprintf("%s/tasks.yml", path))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("initConfig should have panicked")
		}
	}()

	initConfig(ctx, cmd)
}

func Test_GetRootCmd_SuccessWithValidate(t *testing.T) {
	ctx := context.TestContext(nil)
	cmd := GetRootCmd(ctx)
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
	err := GetRootPreRunEFn(ctx, true)(cmd, []string{})
	assert.NoError(t, err)
}

func Test_GetRootCmd_SuccessWithoutValidate(t *testing.T) {
	ctx := context.TestContext(nil)
	cmd := GetRootCmd(ctx)
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
	cmd.SetArgs([]string{fmt.Sprintf("--%s", Config), fmt.Sprintf("%s/tasks.yml", path)})
	_ = cmd.Execute()
	err := GetRootPreRunEFn(ctx, false)(cmd, []string{})
	assert.NoError(t, err)
}

func Test_GetRootCmd_SuccessWithLogLevel(t *testing.T) {
	ctx := context.TestContext(nil)
	cmd := GetRootCmd(ctx)
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

	cmd.SetArgs([]string{fmt.Sprintf("--%s", Config), fmt.Sprintf("%s/tasks.yml", path), fmt.Sprintf("--%s", LogLevel), "WARN"})
	_ = cmd.Execute()
	err := GetRootPreRunEFn(ctx, false)(cmd, []string{})
	assert.NoError(t, err)
}

func Test_GetRootCmd_ErrorWithLogLevel(t *testing.T) {
	ctx := context.TestContext(nil)
	cmd := GetRootCmd(ctx)
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

	cmd.SetArgs([]string{fmt.Sprintf("--%s", Config), fmt.Sprintf("%s/tasks.yml", path), fmt.Sprintf("--%s", LogLevel), "WRONG"})
	_ = cmd.Execute()
	err := GetRootPreRunEFn(ctx, false)(cmd, []string{})
	assert.Error(t, err)
}

func Test_GetRootCmd_ErrorValidate(t *testing.T) {
	ctx := context.TestContext(nil)
	cmd := GetRootCmd(ctx)
	fsFake := afero.NewMemMapFs()
	viper.Reset()
	viper.SetFs(fsFake)
	path := "/app"
	_ = fsFake.Mkdir(path, 0775)
	_ = afero.WriteFile(
		fsFake,
		fmt.Sprintf("%s/tasks.yml", path),
		[]byte("workers:\n- {id: 'test',command: 'fake'}\n- {id: 'test',command: 'fake'}\nscheduled:\n- {id: 'test',expr: '0 0 * * *'}\n"),
		0644,
	)

	cmd.SetArgs([]string{fmt.Sprintf("--%s", Config), fmt.Sprintf("%s/tasks.yml", path)})
	_ = cmd.Execute()
	err := GetRootPreRunEFn(ctx, true)(cmd, []string{})
	assert.Error(t, err)
}
