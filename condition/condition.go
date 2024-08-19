package condition

import "github.com/hashicorp/go-bexpr"

func EvalExpression(expression string, envVars map[string]string) (bool, error) {
	if expression == "" {
		return true, nil
	}
	eval, err := bexpr.CreateEvaluator(expression)
	if err != nil {
		return false, err
	}
	return eval.Evaluate(envVars)
}
