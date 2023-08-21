package flags

import (
	"github.com/spf13/cobra"
	"os"
	osUser "os/user"
)

const (
	WorkingDir = "working-dir"
	User       = "user"
	GroupName  = "group-name"
)

func AddFlagWorkingDir(cmd *cobra.Command) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmd.Flags().StringP(
		WorkingDir,
		"w",
		path,
		"Define working directory",
	)
}

func AddFlagUser(cmd *cobra.Command) {
	user, err := osUser.Current()
	if err != nil {
		panic(err.Error())
	}
	cmd.Flags().StringP(
		User,
		"u",
		user.Username,
		"Define user used to run command",
	)
}

func AddFlagGroupName(cmd *cobra.Command) {
	cmd.Flags().StringP(
		GroupName,
		"g",
		"",
		"Define group name",
	)
}
