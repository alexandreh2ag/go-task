package cli

import (
	"errors"
	"fmt"
	"github.com/alexandreh2ag/go-task/context"
	gtaskValidator "github.com/alexandreh2ag/go-task/validator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

func GetValidateCmd(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "validate",
		Short:             "validate config",
		SilenceUsage:      true,
		SilenceErrors:     true,
		PersistentPreRunE: GetRootPreRunEFn(ctx, false),
		RunE:              GetValidateRunFn(ctx),
	}

	return cmd
}

func GetValidateRunFn(ctx *context.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		validate := gtaskValidator.New()
		err := validate.Struct(ctx.Config)
		if err != nil {

			var validationErrors validator.ValidationErrors
			switch {
			case errors.As(err, &validationErrors):
				for _, validationError := range validationErrors {
					ctx.Logger.Error(fmt.Sprintf("%v", validationError))
				}
			default:
				ctx.Logger.Error(fmt.Sprintf("%v", err))
			}
			return errors.New("configuration file is not valid")
		}
		ctx.Logger.Info("configuration file is valid")

		return nil
	}
}
