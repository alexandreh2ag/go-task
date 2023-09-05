package cli

import (
	"github.com/alexandreh2ag/go-task/cli/worker"
	"github.com/alexandreh2ag/go-task/context"
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
