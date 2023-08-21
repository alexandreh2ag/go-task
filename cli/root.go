package cli

import (
	appCtx "alexandreh2ag/go-task/context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log/slog"

	"path"
	"strings"

	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Config   = "config"
	LogLevel = "level"
	Name     = "gtask"
)

func GetRootCmd(ctx *appCtx.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:               Name,
		Short:             "go task: run & generate config for scheduled / cron task",
		PersistentPreRunE: GetRootPreRunEFn(ctx),
	}

	cmd.PersistentFlags().StringP(Config, "c", "", "Define config path")
	cmd.PersistentFlags().StringP(LogLevel, "l", "INFO", "Define log level")
	_ = viper.BindPFlag(Config, cmd.Flags().Lookup(Config))
	_ = viper.BindPFlag(LogLevel, cmd.Flags().Lookup(LogLevel))

	return cmd
}

func GetRootPreRunEFn(ctx *appCtx.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var err error
		//initConfig(ctx, cmd)

		validate := validator.New()
		err = validate.Struct(ctx.Config)
		if err != nil {
			return err
		}
		logLevelFlagStr, _ := cmd.Flags().GetString(LogLevel)
		level := slog.LevelInfo
		err = level.UnmarshalText([]byte(logLevelFlagStr))
		if err != nil {
			return err
		}
		ctx.LogLevel.Set(level)

		return nil
	}
}

func initConfig(ctx *appCtx.Context, cmd *cobra.Command) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(fmt.Errorf("unable to find current path, %v", err))
	}

	viper.SetConfigName("tasks")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix(strings.ToUpper(Name))

	if err = viper.BindPFlags(cmd.Flags()); err != nil {
		fmt.Println(err)
	}
	configPath := viper.GetString(Config)
	if configPath != "" {
		viper.SetConfigFile(configPath)
		configDir := path.Dir(configPath)
		if configDir != "." && configDir != dir {
			viper.AddConfigPath(configDir)
		}
	}

	// If a config file is found, read it in.
	if err = viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("Config: " + err.Error())
	}

	err = viper.Unmarshal(ctx.Config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into config struct, %v", err))
	}
}
