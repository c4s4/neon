package task

import (
	"neon/build"
	"neon/util"
)

func init() {
	build.TaskMap["try"] = build.TaskDescriptor{
		Constructor: Try,
		Help: `Try/catch construct.

Arguments:
- try: steps to execute.
- catch: executed if an error occurs.

Examples:
# execute a system command and continue even if it fails
- try:
  - "command-that-doesnt-exist"
  catch:
  - print: "command failed!"

Notes:
- The error message for the failure is stored in 'error' variable as text.`,
	}
}

func Try(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"try", "catch"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	trySteps, err := ParseSteps(target, args, "try")
	if err != nil {
		return nil, err
	}
	catchSteps, err := ParseSteps(target, args, "catch")
	if err != nil {
		return nil, err
	}
	return func() error {
		target.Build.Context.SetProperty("error", "")
		err := RunSteps(target.Build, trySteps)
		if err != nil {
			target.Build.Context.SetProperty("error", err.Error())
			err = RunSteps(target.Build, catchSteps)
			if err != nil {
				return err
			}
		}
		return nil
	}, nil
}
