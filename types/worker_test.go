package types

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_WorkerTask_SuccessValidate(t *testing.T) {
	validate := validator.New()
	worker := WorkerTask{
		Id:      "test",
		Command: "fake",
	}
	err := validate.Struct(worker)

	assert.NoError(t, err)
}

func Test_WorkerTask_SuccessValidateWithOptionalData(t *testing.T) {
	validate := validator.New()
	worker := WorkerTask{
		Id:        "test",
		Command:   "fake",
		User:      "test",
		Directory: "/tmp/test/",
	}
	err := validate.Struct(worker)

	assert.NoError(t, err)
}

func Test_WorkerTask_ErrorValidate(t *testing.T) {
	validate := validator.New()
	worker := WorkerTask{
		Id: "test",
	}
	err := validate.Struct(worker)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'Command' failed on the 'required' tag")
}

func Test_WorkerTask_ErrorValidateComplex(t *testing.T) {
	validate := validator.New()
	worker := WorkerTask{
		Id:        "test",
		Command:   "fake",
		User:      "user/test",
		Directory: "wrong",
	}
	err := validate.Struct(worker)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'User' failed on the 'alphanum' tag")
	assert.Contains(t, err.Error(), "Field validation for 'Directory' failed on the 'dirpath' tag")
}
