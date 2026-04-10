package types

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/alexandreh2ag/go-task/env"
)

type WorkerTasks = []*WorkerTask

type WorkerTaskTemplate struct {
	ExtraParams map[string]map[string]string `mapstructure:"extra_params" validate:"omitempty,dive,dive,required"`
}

type WorkerTask struct {
	Id         string `mapstructure:"id" validate:"required,excludesall=!@#$ "`
	Command    string `mapstructure:"command" validate:"required"`
	GroupName  string
	Expression string             `mapstructure:"if"`
	User       string             `mapstructure:"user" validate:"omitempty,required,alphanum"`
	Directory  string             `mapstructure:"directory" validate:"omitempty,required,dirpath"`
	Envs       map[string]string  `mapstructure:"environments"`
	Template   WorkerTaskTemplate `mapstructure:"template" validate:"omitempty"`
}

func PrepareWorkerTasks(tasks WorkerTasks, groupName, user, workingDir string, enVars map[string]string) {
	for _, task := range tasks {
		task.GroupName = groupName
		task.Envs = env.ToUpperKeys(task.Envs)
		if task.Template.ExtraParams == nil {
			task.Template.ExtraParams = map[string]map[string]string{}
		}
		_ = mergo.Merge(&task.Envs, enVars, mergo.WithOverride)
		if task.User == "" {
			task.User = user
		}

		if task.Directory == "" {
			task.Directory = workingDir
		}
		taskVars := map[string]string{
			GtaskGroupNameKey: task.GroupName,
			GtaskDirKey:       task.Directory,
			GtaskUserKey:      task.User,
			GtaskIDKey:        task.PrefixedName(),
		}
		_ = mergo.Merge(&task.Envs, taskVars, mergo.WithOverride)

		task.Envs = env.EvalAll(task.Envs)
	}

}

func (w *WorkerTask) PrefixedName() string {
	return fmt.Sprintf("%s-%s", w.GroupName, w.Id)
}
