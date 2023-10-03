package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateCronExpr(t *testing.T) {
	type args struct {
		Cron string `validate:"cron-expr"`
	}
	validate := validator.New()
	_ = validate.RegisterValidation(CronExprKey, ValidateCronExpr)

	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "successFullWildcard",
			args:    args{Cron: "* * * * *"},
			wantErr: assert.NoError,
		},
		{
			name:    "successRangeMinutes",
			args:    args{Cron: "5-20 * * * *"},
			wantErr: assert.NoError,
		},
		{
			name:    "successRangeHour",
			args:    args{Cron: "0 5-20 * * *"},
			wantErr: assert.NoError,
		},
		{
			name:    "successStepHour",
			args:    args{Cron: "0 */2 * * *"},
			wantErr: assert.NoError,
		},
		{
			name:    "successStepMinute",
			args:    args{Cron: "*/2 * * * *"},
			wantErr: assert.NoError,
		},
		{
			name:    "successOneMinuteHourDay",
			args:    args{Cron: "1 1 1 * *"},
			wantErr: assert.NoError,
		},
		{
			name:    "failStepDigitMissing",
			args:    args{Cron: "* */ * * *"},
			wantErr: assert.Error,
		},
		{
			name:    "failRangeDigitMissing",
			args:    args{Cron: "* 0- * * *"},
			wantErr: assert.Error,
		},
		{
			name:    "failEmpty",
			args:    args{Cron: ""},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.args)
			tt.wantErr(t, err, "ValidateCronExpr is not valid")
		})
	}
}
