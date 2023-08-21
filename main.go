package main

import (
	"alexandreh2ag/go-task/cli"
	"alexandreh2ag/go-task/context"
)

func main() {
	ctx := context.DefaultContext()
	rootCmd := cli.GetRootCmd(ctx)
	rootCmd.AddCommand(
		cli.GetWorkerCmd(ctx),
		cli.GetScheduleCmd(ctx),
		cli.GetVersionCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
