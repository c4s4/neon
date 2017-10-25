package task

import (
	"fmt"
	"neon/build"
	"neon/util"
)

func init() {
	build.TaskMap["assert"] = build.TaskDescriptor{
		Constructor: Assert,
		Help: `Make an assertion and fail if assertion condition is false.

Arguments:

- assert: the assertion to perform (as a script expression).

Examples:

    # assert that foo == "bar", and fail otherwize
    - assert: 'foo == "bar"'`,
	}
}

func Assert(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"assert"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	assertion, err := args.GetString("assert")
	if err != nil {
		return nil, fmt.Errorf("evaluating assert construct: %v", err)
	}
	return func(context *build.Context) error {
		_result, _err := context.EvaluateExpression(assertion)
		if _err != nil {
			return fmt.Errorf("evaluating 'assert' condition: %v", _err)
		}
		_pass, _ok := _result.(bool)
		if !_ok {
			return fmt.Errorf("evaluating assert condition: must return a bool")
		}
		if !_pass {
			return fmt.Errorf("assertion '%s' failed", assertion)
		}
		return nil
	}, nil
}
