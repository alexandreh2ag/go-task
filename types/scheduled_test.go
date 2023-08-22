package types

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ScheduledTask_SuccessValidate(t *testing.T) {
	validate := validator.New()
	scheduled := ScheduledTask{
		Id:       "test",
		CronExpr: "* * * * *",
		Command:  "fake",
	}
	err := validate.Struct(scheduled)

	assert.NoError(t, err)
}

func Test_ScheduledTask_SuccessValidateWithOptionalData(t *testing.T) {
	validate := validator.New()
	scheduled := ScheduledTask{
		Id:        "test",
		CronExpr:  "* * * * *",
		Command:   "fake",
		User:      "test",
		Directory: "/tmp/test/",
	}
	err := validate.Struct(scheduled)

	assert.NoError(t, err)
}

func Test_ScheduledTask_ErrorValidate(t *testing.T) {
	validate := validator.New()
	scheduled := ScheduledTask{
		Id:       "test",
		CronExpr: "* * * * *",
	}
	err := validate.Struct(scheduled)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'Command' failed on the 'required' tag")
}

func Test_ScheduledTask_ErrorValidateComplex(t *testing.T) {
	validate := validator.New()
	scheduled := ScheduledTask{
		Id:        "test",
		CronExpr:  "wrong",
		Command:   "fake",
		User:      "user/test",
		Directory: "wrong",
	}
	err := validate.Struct(scheduled)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'CronExpr' failed on the 'cron' tag")
	assert.Contains(t, err.Error(), "Field validation for 'User' failed on the 'alphanum' tag")
	assert.Contains(t, err.Error(), "Field validation for 'Directory' failed on the 'dirpath' tag")
}
