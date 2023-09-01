package flags

import (
	"github.com/spf13/cobra"
	"os"
	osUser "os/user"
)

const (
	WorkingDir    = "working-dir"
	User          = "user"
	GroupName     = "group-name"
	TimeZone      = "timezone"
	ResultPath    = "result-path"
	NoResultPrint = "no-result-print"
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

func AddFlagTimezone(cmd *cobra.Command) {
	cmd.Flags().StringP(
		TimeZone,
		"t",
		"",
		"Define timezone for used to calcul cron expression",
	)
}

func AddFlagNoResultPrint(cmd *cobra.Command) {
	cmd.Flags().Bool(
		NoResultPrint,
		false,
		"Flag to not print tasks results",
	)
}

func AddFlagResultPath(cmd *cobra.Command) {
	cmd.Flags().String(
		ResultPath,
		"",
		"Define path to save tasks results (default: no logs file)",
	)
}
