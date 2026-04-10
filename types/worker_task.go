package types

import (
	"fmt"

	"dario.cat/mergo"
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
	ParentId   string
	Expression string             `mapstructure:"if"`
	User       string             `mapstructure:"user" validate:"omitempty,required,alphanum"`
	Directory  string             `mapstructure:"directory" validate:"omitempty,required,dirpath"`
	Envs       map[string]string  `mapstructure:"environments"`
	Instances  int                `mapstructure:"instances" validate:"omitempty,min=1"`
	Template   WorkerTaskTemplate `mapstructure:"template" validate:"omitempty"`
}

func prefixedName(groupName, id string) string {
	return fmt.Sprintf("%s-%s", groupName, id)
}

func (w *WorkerTask) PrefixedName() string {
	return prefixedName(w.GroupName, w.Id)
}

func (w *WorkerTask) PrefixedParentName() string {
	return prefixedName(w.GroupName, w.ParentId)
}

func ExpandWorkerTasks(tasks WorkerTasks) WorkerTasks {
	expanded := WorkerTasks{}
	for _, task := range tasks {
		task.ParentId = task.Id
		instances := task.Instances
		if instances <= 1 {
			expanded = append(expanded, task)
			continue
		}
		for i := 1; i <= instances; i++ {
			clone := *task
			clone.Id = fmt.Sprintf("%s_%d", task.Id, i)
			clone.Instances = 1
			if task.Envs != nil {
				clone.Envs = make(map[string]string, len(task.Envs))
				for k, v := range task.Envs {
					clone.Envs[k] = v
				}
			}
			expanded = append(expanded, &clone)
		}
	}
	return expanded
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
			GtaskGroupNameKey:  task.GroupName,
			GtaskDirKey:        task.Directory,
			GtaskUserKey:       task.User,
			GtaskIDKey:         task.PrefixedParentName(),
			GtaskInstanceIDKey: task.PrefixedName(),
		}
		_ = mergo.Merge(&task.Envs, taskVars, mergo.WithOverride)

		task.Envs = env.EvalAll(task.Envs)
	}

}
