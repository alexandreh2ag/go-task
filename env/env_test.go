package env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_GetEnvVars(t *testing.T) {
	key := "GTASK_TESTING_GETENVVARS"
	value := "foo"
	_ = os.Setenv(key, value)
	tests := []struct {
		name      string
		key       string
		extraVars map[string]string
		want      string
	}{
		{
			name: "SuccessVarNotExist",
			key:  key + "WRONG",
			want: "",
		},
		{
			name: "SuccessOSVar",
			key:  key,
			want: "foo",
		},
		{
			name:      "SuccessGtaskVar",
			key:       "mykey",
			extraVars: map[string]string{"mykey": "bar"},
			want:      "bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetEnvVars(tt.extraVars)(tt.key), "getEnvVars(%v)", tt.extraVars)
		})
	}
}

func Test_ToUpperKeys(t *testing.T) {
	tests := []struct {
		name string
		args map[string]string
		want map[string]string
	}{
		{
			name: "NoChange",
			args: map[string]string{
				"KEY1": "foo",
				"KEY2": "bar",
			},
			want: map[string]string{
				"KEY1": "foo",
				"KEY2": "bar",
			},
		},
		{
			name: "WithChange",
			args: map[string]string{
				"Key1": "foo",
				"key2": "bar",
			},
			want: map[string]string{
				"KEY1": "foo",
				"KEY2": "bar",
			},
		},
	}
	for _, tt := range tests {
		got := ToUpperKeys(tt.args)
		assert.Equal(t, tt.want, got)
	}

}

func Test_EvalAll(t *testing.T) {

	_ = os.Setenv("OS_VAR", "FROM_OS")
	tests := []struct {
		name string
		args map[string]string
		want map[string]string
	}{
		{
			name: "NoChange",
			args: map[string]string{
				"KEY1": "foo",
				"KEY2": "bar",
			},
			want: map[string]string{
				"KEY1": "foo",
				"KEY2": "bar",
			},
		},
		{
			name: "WithChange",
			args: map[string]string{
				"KEY1": "${OS_VAR}",
				"KEY2": "${KEY3}",
				"KEY3": "FOO",
			},
			want: map[string]string{
				"KEY1": "FROM_OS",
				"KEY2": "FOO",
				"KEY3": "FOO",
			},
		},
	}
	for _, tt := range tests {
		got := EvalAll(tt.args)
		assert.Equal(t, tt.want, got)
	}
	_ = os.Unsetenv("OS_VAR")

}

func Test_GetEnvs(t *testing.T) {
	_ = os.Setenv("TEST_VAR", "NEED_ME")
	envVars := GetEnvs()
	assert.Equal(t, envVars["TEST_VAR"], "NEED_ME")
	_ = os.Unsetenv("TEST_VAR")
}
