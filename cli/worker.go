package cli

import (
	"alexandreh2ag/go-task/cli/worker"
	"alexandreh2ag/go-task/context"
	"github.com/spf13/cobra"
)

func GetWorkerCmd(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "worker sub commands",
	}

	cmd.AddCommand(worker.GetWorkerGenerateCmd(ctx))

	return cmd
}
