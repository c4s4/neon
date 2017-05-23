package task

import (
	"neon/build"
	"neon/util"
)

func init() {
	build.TaskMap["try"] = build.TaskDescriptor{
		Constructor: Try,
		Help: `Try/catch/finally construct.

Arguments:

- try: steps to execute.
- catch: executed if an error occurs (optional).
- finally: executed in all cases (optional).

Examples:

    # execute a command and continue even if it fails
    - try:
      - "command-that-doesnt-exist"
	- print: "Continue even if command fails"
	# execute a command and print a message if it fails
	- try:
	  - "command-that-doesnt-exist"
	  catch:
	  - print: "There was an error!"
	# execute a command a print message in all cases
	- try:
	  - "command-that-doesnt-exist"
	  finally:
	  - print: "Print whatever happens"

Notes:

- The error message for the failure is stored in '_error' variable as text.`,
	}
}

func Try(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"try", "catch", "finally"}
	if err := CheckFields(args, fields, fields[:1]); err != nil {
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
	finallySteps, err := ParseSteps(target, args, "finally")
	if err != nil {
		return nil, err
	}
	return func() error {
		_depth := target.Build.Index.Len()
		target.Build.Context.SetProperty("_error", "")
		var _tryError error
		var _catchError error
		var _finallyError error
		_tryError = RunSteps(target.Build, trySteps)
		if _tryError != nil {
			for target.Build.Index.Len() > _depth {
				target.Build.Index.Shrink()
			}
			if len(catchSteps) > 0 || (len(catchSteps) == 0 && len(finallySteps) == 0) {
				target.Build.Context.SetProperty("_error", _tryError.Error())
				_tryError = nil
				_catchError = RunSteps(target.Build, catchSteps)
				if _catchError != nil {
					for target.Build.Index.Len() > _depth {
						target.Build.Index.Shrink()
					}
				}
			}
		}
		_finallyError = RunSteps(target.Build, finallySteps)
		if _finallyError != nil {
			for target.Build.Index.Len() > _depth {
				target.Build.Index.Shrink()
			}
			target.Build.Context.SetProperty("_error", _finallyError.Error())
			return _finallyError
		}
		if _catchError != nil {
			target.Build.Context.SetProperty("_error", _catchError.Error())
			return _catchError
		}
		if _tryError != nil {
			target.Build.Context.SetProperty("_error", _tryError.Error())
			return _tryError
		}
		return nil
	}, nil
}
