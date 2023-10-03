package validator

import "github.com/go-playground/validator/v10"

func New(options ...validator.Option) *validator.Validate {
	validate := validator.New()
	_ = validate.RegisterValidation(CronExprKey, ValidateCronExpr)
	return validate
}
