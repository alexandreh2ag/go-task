package types

type WorkerTasks = []*WorkerTask

type WorkerTask struct {
	Id        string `mapstructure:"id" validate:"required,excludesall=!@#$ "`
	Command   string `mapstructure:"command" validate:"required"`
	User      string `mapstructure:"user" validate:"omitempty,required,alphanum"`
	Directory string `mapstructure:"directory" validate:"omitempty,required,dirpath"`
}

func PrepareWorkerTasks(tasks WorkerTasks, user, workingDir string) {
	for _, task := range tasks {
		if task.User == "" {
			task.User = user
		}

		if task.Directory == "" {
			task.Directory = workingDir
		}
	}
}
