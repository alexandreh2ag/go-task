package types

import (
	"dario.cat/mergo"
	"fmt"
	"golang.org/x/exp/maps"
	"strings"
)

type WorkerTasks = []*WorkerTask

type WorkerTask struct {
	Id        string `mapstructure:"id" validate:"required,excludesall=!@#$ "`
	Command   string `mapstructure:"command" validate:"required"`
	GroupName string
	User      string            `mapstructure:"user" validate:"omitempty,required,alphanum"`
	Directory string            `mapstructure:"directory" validate:"omitempty,required,dirpath"`
	Envs      map[string]string `mapstructure:"environments"`
}

func PrepareWorkerTasks(tasks WorkerTasks, groupName, user, workingDir string, enVars map[string]string) {
	for _, task := range tasks {
		task.GroupName = groupName

		// workaround with viper issue https://github.com/spf13/viper/issues/1014
		keys := maps.Keys(task.Envs)
		for _, key := range keys {
			if key != strings.ToUpper(key) {
				task.Envs[strings.ToUpper(key)] = task.Envs[key]
				delete(task.Envs, key)
			}

		}

		_ = mergo.Merge(&task.Envs, enVars, mergo.WithOverride)
		if task.User == "" {
			task.User = user
		}

		if task.Directory == "" {
			task.Directory = workingDir
		}

	}
}

func (w *WorkerTask) PrefixedName() string {
	return fmt.Sprintf("%s-%s", w.GroupName, w.Id)
}
