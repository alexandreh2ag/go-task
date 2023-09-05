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

const (
	Format     = "format"
	OutputPath = "output"
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
		Format,
		"f",
		generate.FormatSupervisor,
		"Choose format",
	)
	cmd.Flags().StringP(
		OutputPath,
		"o",
		fmt.Sprintf("%s/workers.conf", outputPath),
		"Choose output path",
	)

	return cmd
}

func GetWorkerGenerateRunFn(ctx *context.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString(Format)
		user, _ := cmd.Flags().GetString(flags.User)
		workingDir, _ := cmd.Flags().GetString(flags.WorkingDir)
		outputPath, _ := cmd.Flags().GetString(OutputPath)
		groupName, _ := cmd.Flags().GetString(flags.GroupName)

		if groupName == "" || outputPath == "" {
			return fmt.Errorf("missing mandatory arguments (--%s, --%s)", OutputPath, flags.GroupName)
		}

		types.PrepareWorkerTasks(ctx.Config.Workers, user, workingDir)
		ctx.Logger.Info(fmt.Sprintf("Generate format type %s", generate.FormatSupervisor))

		return generate.Generate(ctx, outputPath, format, groupName)
	}
}
