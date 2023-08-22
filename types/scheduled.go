package types

type ScheduledTasks = []ScheduledTask

type ScheduledTask struct {
	Id        string `mapstructure:"id" validate:"required,alphanumunicode"`
	CronExpr  string `mapstructure:"expr" validate:"required,cron"`
	Command   string `mapstructure:"command" validate:"required,alphanumunicode"`
	User      string `mapstructure:"user" validate:"omitempty,required,alphanum"`
	Directory string `mapstructure:"directory" validate:"omitempty,required,dirpath"`
}
