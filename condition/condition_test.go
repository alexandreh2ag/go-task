package condition

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_EvalExpression_NoError(t *testing.T) {
	tests := []struct {
		name string
		expr string
		env  map[string]string
		want bool
	}{
		{
			name: "Empty Expression",
			expr: "",
			env:  map[string]string{},
			want: true,
		},
		{
			name: "True",
			expr: "VAR == OK",
			env: map[string]string{
				"VAR": "OK",
			},
			want: true,
		},
		{
			name: "False",
			expr: "VAR != OK",
			env: map[string]string{
				"VAR": "OK",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		gotBool, gotErr := EvalExpression(tt.expr, tt.env)
		assert.Equal(t, tt.want, gotBool)
		assert.NoError(t, gotErr)
	}
}

func Test_EvalExpression_Error(t *testing.T) {

	expr := "VAR & OK"
	env := map[string]string{
		"VAR": "OK",
	}
	got, err := EvalExpression(expr, env)
	assert.False(t, got)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no match found, expected:")

}
