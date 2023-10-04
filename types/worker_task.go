package types

import (
	"fmt"
	"slices"
	"strings"
)

type WorkerTasks []*WorkerTask

func (ws WorkerTasks) GetUniqueExtraGroups() []string {
	uniqueExtraGroups := []string{}
	for _, task := range ws {
		for _, group := range task.ExtraGroups {
			if !slices.Contains(uniqueExtraGroups, group) {
				uniqueExtraGroups = append(uniqueExtraGroups, group)
			}
		}
	}
	return uniqueExtraGroups
}

func (ws WorkerTasks) GetTasksInGroup(selectedGroup string) WorkerTasks {
	workerTasks := WorkerTasks{}
	if strings.Compare("", selectedGroup) == 0 {
		return ws
	}

	for _, task := range ws {
		if slices.Contains(task.ExtraGroups, selectedGroup) {
			workerTasks = append(workerTasks, task)
		}
	}
	return workerTasks
}

func (ws WorkerTasks) GetAllPrefixedId() []string {
	ids := []string{}
	for _, task := range ws {
		ids = append(ids, task.PrefixedId())
	}
	return ids
}

type WorkerTask struct {
	Id          string `mapstructure:"id" validate:"required,excludesall=!@#$ "`
	Command     string `mapstructure:"command" validate:"required"`
	GroupName   string
	ExtraGroups []string `mapstructure:"extra_groups" validate:"omitempty,dive,excludesall=!@#$ "`
	User        string   `mapstructure:"user" validate:"omitempty,required,alphanum"`
	Directory   string   `mapstructure:"directory" validate:"omitempty,required,dirpath"`
}

func PrepareWorkerTasks(tasks WorkerTasks, groupName, user, workingDir string) {
	for _, task := range tasks {
		task.GroupName = groupName

		if task.User == "" {
			task.User = user
		}

		if task.Directory == "" {
			task.Directory = workingDir
		}
	}
}

func (w *WorkerTask) PrefixedId() string {
	return fmt.Sprintf("%s-%s", w.GroupName, w.Id)
}
