package task

import (
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc {
		Name: "try",
		Func: Try,
		Args: reflect.TypeOf(TryArgs{}),
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
	})
}

type TryArgs struct {
	Try     []build.Step `steps`
	Catch   []build.Step `steps`
	Finally []build.Step `steps`
}


func Try(context *build.Context, args interface{}) error {
	params := args.(TryArgs)
	depth := context.Index.Len()
	context.SetProperty("_error", "")
	var tryError error
	var catchError error
	var finallyError error
	tryError = context.Run(params.Try)
	if tryError != nil {
		for context.Index.Len() > depth {
			context.Index.Shrink()
		}
		if len(params.Catch) > 0 || (len(params.Catch) == 0 && len(params.Finally) == 0) {
			context.SetProperty("_error", tryError.Error())
			tryError = nil
			catchError = context.Run(params.Catch)
			if catchError != nil {
				for context.Index.Len() > depth {
					context.Index.Shrink()
				}
			}
		}
	}
	finallyError = context.Run(params.Finally)
	if finallyError != nil {
		for context.Index.Len() > depth {
			context.Index.Shrink()
		}
		context.SetProperty("_error", finallyError.Error())
		return finallyError
	}
	if catchError != nil {
		context.SetProperty("_error", catchError.Error())
		return catchError
	}
	if tryError != nil {
		context.SetProperty("_error", tryError.Error())
		return tryError
	}
	return nil
}
