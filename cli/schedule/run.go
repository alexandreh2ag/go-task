package schedule

import (
	"github.com/alexandreh2ag/go-task/cli/flags"
	"github.com/alexandreh2ag/go-task/context"
	"github.com/alexandreh2ag/go-task/schedule"
	"github.com/alexandreh2ag/go-task/types"
	"github.com/spf13/cobra"
)

func GetScheduleRunCmd(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run workers based on each cron expr",
		RunE:  GetScheduleRunRunFn(ctx),
	}

	flags.AddFlagWorkingDir(cmd)
	flags.AddFlagTimezone(cmd)
	flags.AddFlagNoResultPrint(cmd)
	flags.AddFlagResultPath(cmd)

	return cmd
}

func GetScheduleRunRunFn(ctx *context.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		workingDir, _ := cmd.Flags().GetString(flags.WorkingDir)
		timezone, _ := cmd.Flags().GetString(flags.TimeZone)
		noResultPrint, _ := cmd.Flags().GetBool(flags.NoResultPrint)
		resultPath, _ := cmd.Flags().GetString(flags.ResultPath)

		types.PrepareScheduledTasks(ctx.Config.Scheduled, ctx.Logger, workingDir)
		refTime, err := schedule.GetCurrentTime(ctx.Clock.Now(), timezone)
		if err != nil {
			return err
		}
		schedule.Run(ctx, refTime, noResultPrint, resultPath)

		return nil
	}
}
