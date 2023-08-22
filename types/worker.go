package types

type WorkerTasks = []WorkerTask

type WorkerTask struct {
	Id        string `mapstructure:"id" validate:"required,alphanumunicode"`
	Command   string `mapstructure:"command" validate:"required,alphanumunicode"`
	User      string `mapstructure:"user" validate:"omitempty,required,alphanum"`
	Directory string `mapstructure:"directory" validate:"omitempty,required,dirpath"`
}
