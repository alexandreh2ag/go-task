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
	TimeZone      = "timezone"
	ResultPath    = "result-path"
	NoResultPrint = "no-result-print"
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

	cmd.Flags().Bool(
		NoResultPrint,
		false,
		"Flag to not print tasks results",
	)

	cmd.Flags().String(
		ResultPath,
		"",
		"Define path to save tasks results (default: no logs file)",
	)

	return cmd
}

func GetScheduleRunRunFn(ctx *context.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		user, _ := cmd.Flags().GetString(flags.User)
		workingDir, _ := cmd.Flags().GetString(flags.WorkingDir)
		timezone, _ := cmd.Flags().GetString(TimeZone)
		noResultPrint, _ := cmd.Flags().GetBool(NoResultPrint)
		resultPath, _ := cmd.Flags().GetString(ResultPath)

		types.PrepareScheduledTasks(ctx.Config.Scheduled, ctx.Logger, user, workingDir)
		refTime, err := schedule.GetCurrentTime(time.Now(), timezone)
		if err != nil {
			return err
		}
		schedule.Run(ctx, refTime, noResultPrint, resultPath)

		return nil
	}
}
