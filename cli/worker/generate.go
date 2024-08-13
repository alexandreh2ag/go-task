package worker

import (
	"fmt"
	"github.com/alexandreh2ag/go-task/cli/flags"
	"github.com/alexandreh2ag/go-task/context"
	"github.com/alexandreh2ag/go-task/generate"
	"github.com/alexandreh2ag/go-task/types"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

const (
	Format     = "format"
	OutputPath = "output"
	EnvVars    = "env"
)

func GetWorkerGenerateCmd(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate config file for worker",
		RunE:  GetWorkerGenerateRunFn(ctx),
	}

	outputPath, _ := os.Getwd()

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

	cmd.Flags().StringSliceP(
		EnvVars,
		"e",
		nil,
		"Injected env vars. Format: KEY1=value1,KEY2=value2",
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

		argsEnvVarsString, _ := cmd.Flags().GetStringSlice(EnvVars)
		envVars := FormatEnvVars(argsEnvVarsString)

		if groupName == "" || outputPath == "" {
			return fmt.Errorf("missing mandatory arguments (--%s, --%s)", OutputPath, flags.GroupName)
		}

		types.PrepareWorkerTasks(ctx.Config.Workers, groupName, user, workingDir, envVars)
		ctx.Logger.Info(fmt.Sprintf("Generate format type %s", generate.FormatSupervisor))

		return generate.Generate(ctx, outputPath, format, groupName)
	}
}

func FormatEnvVars(slice []string) map[string]string {
	result := map[string]string{}
	for _, envVar := range slice {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}
