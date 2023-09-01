package cli

import (
	"alexandreh2ag/go-task/cli/schedule"
	"alexandreh2ag/go-task/context"
	"github.com/spf13/cobra"
)

func GetScheduleCmd(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "schedule sub commands",
	}

	cmd.AddCommand(
		schedule.GetScheduleRunCmd(ctx),
		schedule.GetScheduleStartCmd(ctx),
	)

	return cmd
}
