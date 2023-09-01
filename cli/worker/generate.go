package worker

import (
	"alexandreh2ag/go-task/cli/flags"
	"alexandreh2ag/go-task/context"
	"alexandreh2ag/go-task/generate"
	"alexandreh2ag/go-task/types"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func GetWorkerGenerateCmd(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate config file for worker",
		RunE:  GetWorkerGenerateRunFn(ctx),
	}

	outputPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	flags.AddFlagGroupName(cmd)
	flags.AddFlagUser(cmd)
	flags.AddFlagWorkingDir(cmd)
	cmd.Flags().StringP(
		flags.Format,
		"f",
		flags.FormatSupervisor,
		"Choose format",
	)
	cmd.Flags().StringP(
		flags.OutputPath,
		"o",
		fmt.Sprintf("%s/workers.conf", outputPath),
		"Choose output path",
	)

	return cmd
}

func GetWorkerGenerateRunFn(ctx *context.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString(flags.Format)
		user, _ := cmd.Flags().GetString(flags.User)
		workingDir, _ := cmd.Flags().GetString(flags.WorkingDir)
		outputPath, _ := cmd.Flags().GetString(flags.OutputPath)
		groupName, _ := cmd.Flags().GetString(flags.GroupName)

		types.PrepareWorkerTasks(ctx.Config.Workers, user, workingDir)
		ctx.Logger.Info(fmt.Sprintf("Generate format type %s", flags.FormatSupervisor))

		return generate.Generate(ctx, outputPath, format, groupName)
	}
}
