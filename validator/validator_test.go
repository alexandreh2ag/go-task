package validator

import (
	extValidator "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewValidator(t *testing.T) {

	validator := New()
	assert.IsType(t, &extValidator.Validate{}, validator)
}
