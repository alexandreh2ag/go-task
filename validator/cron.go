package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

const (
	cronRegexString = `(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)|((((\d+,)+\d+|(\d+-\d+)|\d+|\*|\*\/\d+) ?){5,7})`
	CronExprKey     = "cron-expr"
)

var (
	cronRegex = regexp.MustCompile(cronRegexString)
)

func ValidateCronExpr(fl validator.FieldLevel) bool {
	cronString := fl.Field().String()
	return cronRegex.MatchString(cronString)
}
