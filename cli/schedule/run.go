package schedule

import (
	"alexandreh2ag/go-task/cli/flags"
	"alexandreh2ag/go-task/context"
	"alexandreh2ag/go-task/schedule"
	"alexandreh2ag/go-task/types"
	"github.com/spf13/cobra"
	"time"
)

const (
	TimeZone = "timezone"
)

func GetScheduleRunCmd(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run workers based on each cron expr",
		RunE:  GetScheduleRunRunFn(ctx),
	}

	flags.AddFlagUser(cmd)
	flags.AddFlagWorkingDir(cmd)
	cmd.Flags().StringP(
		TimeZone,
		"t",
		"",
		"Define timezone for used to calcul cron expression",
	)

	return cmd
}

func GetScheduleRunRunFn(ctx *context.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		user, _ := cmd.Flags().GetString(flags.User)
		workingDir, _ := cmd.Flags().GetString(flags.WorkingDir)
		timezone, _ := cmd.Flags().GetString(TimeZone)

		types.PrepareScheduledTasks(ctx.Config.Scheduled, ctx.Logger, user, workingDir)
		refTime, err := schedule.GetCurrentTime(time.Now(), timezone)
		if err != nil {
			return err
		}
		schedule.Run(ctx, refTime)

		return nil
	}
}
