package schedule

import (
	"alexandreh2ag/go-task/cli/flags"
	"alexandreh2ag/go-task/context"
	"github.com/spf13/cobra"
)

func GetScheduleRunCmd(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run workers based on each cron expr",
		RunE:  GetScheduleRunRunFn(ctx),
	}

	flags.AddFlagUser(cmd)
	flags.AddFlagWorkingDir(cmd)

	return cmd
}

func GetScheduleRunRunFn(ctx *context.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}
}
