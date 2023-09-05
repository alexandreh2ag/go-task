package schedule

import (
	"errors"
	"github.com/alexandreh2ag/go-task/cli/flags"
	"github.com/alexandreh2ag/go-task/context"
	"github.com/alexandreh2ag/go-task/schedule"
	"github.com/alexandreh2ag/go-task/types"
	"github.com/spf13/cobra"
	"time"
)

const (
	Tick = "tick"
)

func GetScheduleStartCmd(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start process to run workers based on each cron expr",
		RunE:  GetScheduleStartRunFn(ctx),
	}

	flags.AddFlagUser(cmd)
	flags.AddFlagWorkingDir(cmd)
	flags.AddFlagTimezone(cmd)
	flags.AddFlagNoResultPrint(cmd)
	flags.AddFlagResultPath(cmd)
	cmd.Flags().Duration(
		Tick,
		5*time.Minute,
		"Define duration between each tick",
	)

	return cmd
}

func GetScheduleStartRunFn(ctx *context.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {

		user, _ := cmd.Flags().GetString(flags.User)
		workingDir, _ := cmd.Flags().GetString(flags.WorkingDir)
		timezone, _ := cmd.Flags().GetString(flags.TimeZone)
		noResultPrint, _ := cmd.Flags().GetBool(flags.NoResultPrint)
		resultPath, _ := cmd.Flags().GetString(flags.ResultPath)
		tick, _ := cmd.Flags().GetDuration(Tick)

		if tick.Minutes() == 0 {
			return errors.New("tick duration must be higher than 0")
		}

		if tick%time.Minute != 0 {
			return errors.New("tick duration must be only in minutes")
		}

		types.PrepareScheduledTasks(ctx.Config.Scheduled, ctx.Logger, user, workingDir)

		return schedule.Start(ctx, int(tick.Minutes()), timezone, noResultPrint, resultPath)
	}
}
