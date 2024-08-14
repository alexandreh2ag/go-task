package schedule

import (
	"github.com/alexandreh2ag/go-task/cli/flags"
	"github.com/alexandreh2ag/go-task/context"
	"github.com/alexandreh2ag/go-task/schedule"
	"github.com/alexandreh2ag/go-task/types"
	"github.com/spf13/cobra"
	"strings"
)

func GetScheduleRunCmd(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Short:   "run workers based on each cron expr",
		Example: "run task1,task2",
		RunE:    GetScheduleRunRunFn(ctx),
		Args:    cobra.MatchAll(cobra.MaximumNArgs(1)),
	}

	flags.AddFlagWorkingDir(cmd)
	flags.AddFlagTimezone(cmd)
	flags.AddFlagNoResultPrint(cmd)
	flags.AddFlagResultPath(cmd)
	flags.AddFlagForce(cmd)
	flags.AddFlagEnvVars(cmd)

	return cmd
}

func GetScheduleRunRunFn(ctx *context.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		workingDir, _ := cmd.Flags().GetString(flags.WorkingDir)
		timezone, _ := cmd.Flags().GetString(flags.TimeZone)
		noResultPrint, _ := cmd.Flags().GetBool(flags.NoResultPrint)
		resultPath, _ := cmd.Flags().GetString(flags.ResultPath)
		force, _ := cmd.Flags().GetBool(flags.Force)
		envVars, _ := cmd.Flags().GetStringToString(flags.EnvVars)

		taskFilter := []string{}
		if len(args) == 1 {
			taskFilter = strings.Split(args[0], ",")
		}

		types.PrepareScheduledTasks(ctx.Config.Scheduled, ctx.Logger, workingDir, envVars)
		refTime, err := schedule.GetCurrentTime(ctx.Clock.Now(), timezone)
		if err != nil {
			return err
		}
		schedule.Run(ctx, refTime, taskFilter, force, noResultPrint, resultPath)

		return nil
	}
}
