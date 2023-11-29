package types

import "fmt"

type WorkerTasks = []*WorkerTask

type WorkerTask struct {
	Id        string `mapstructure:"id" validate:"required,excludesall=!@#$ "`
	Command   string `mapstructure:"command" validate:"required"`
	GroupName string
	User      string `mapstructure:"user" validate:"omitempty,required,alphanum"`
	Directory string `mapstructure:"directory" validate:"omitempty,required,dirpath"`
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

func (w *WorkerTask) PrefixedName() string {
	return fmt.Sprintf("%s-%s", w.GroupName, w.Id)
}
