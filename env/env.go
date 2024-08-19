package env

import (
	"golang.org/x/exp/maps"
	"os"
	"strings"
)

func GetEnvVars(extraVars map[string]string) func(string) string {
	return func(key string) string {

		val, ok := extraVars[key]
		if !ok {
			val = os.Getenv(key)
		}
		return os.ExpandEnv(val)
	}
}

func ToUpperKeys(envs map[string]string) map[string]string {
	// workaround with viper issue https://github.com/spf13/viper/issues/1014
	keys := maps.Keys(envs)
	for _, key := range keys {
		if key != strings.ToUpper(key) {
			envs[strings.ToUpper(key)] = envs[key]
			delete(envs, key)
		}
	}
	return envs
}

func EvalAll(envs map[string]string) map[string]string {
	evaluatedEnvs := map[string]string{}
	for key, value := range envs {
		evaluatedEnvs[key] = os.Expand(value, GetEnvVars(envs))
	}
	return evaluatedEnvs
}
