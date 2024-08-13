package types

import (
	"dario.cat/mergo"
	"fmt"
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
