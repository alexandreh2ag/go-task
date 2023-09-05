package main

import (
	"github.com/alexandreh2ag/go-task/cli"
	"github.com/alexandreh2ag/go-task/context"
)

func main() {
	ctx := context.DefaultContext()
	rootCmd := cli.GetRootCmd(ctx)
	rootCmd.AddCommand(
		cli.GetWorkerCmd(ctx),
		cli.GetScheduleCmd(ctx),
		cli.GetValidateCmd(ctx),
		cli.GetVersionCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
