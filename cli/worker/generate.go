package worker

import (
	"alexandreh2ag/go-task/cli/flags"
	"alexandreh2ag/go-task/context"
	"fmt"
	"github.com/spf13/cobra"
)

const (
	Format           = "format"
	FormatSupervisor = "supervisor"
)

func GetWorkerGenerateCmd(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate config file for worker",
		RunE:  GetWorkerGenerateRunFn(ctx),
	}

	flags.AddFlagGroupName(cmd)
	flags.AddFlagUser(cmd)
	flags.AddFlagWorkingDir(cmd)
	cmd.Flags().StringP(
		Format,
		"f",
		FormatSupervisor,
		"Choose format",
	)

	return cmd
}

func GetWorkerGenerateRunFn(ctx *context.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString(Format)
		switch format {
		case FormatSupervisor:
			ctx.Logger.Info(fmt.Sprintf("Generate format type %s", FormatSupervisor))
		}

		return nil
	}
}
