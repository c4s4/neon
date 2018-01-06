package task

import (
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "try",
		Func: Try,
		Args: reflect.TypeOf(TryArgs{}),
		Help: `Try/catch/finally construct.

Arguments:

- try: steps to execute (steps).
- catch: executed if an error occurs (steps, optional).
- finally: executed in any case (steps, optional).

Examples:

    # execute a command and continue even if it fails
    - try:
      - 'command-that-doesnt-exist'
	- print: 'Continue even if command fails'
	# execute a command and print a message if it fails
	- try:
	  - 'command-that-doesnt-exist'
	  catch:
	  - print: 'There was an error!'
	# execute a command a print message in all cases
	- try:
	  - 'command-that-doesnt-exist'
	  finally:
	  - print: 'Print whatever happens'

Notes:

- The error message for the failure is stored in '_error' variable as text.`,
	})
}

type TryArgs struct {
	Try     build.Steps `steps`
	Catch   build.Steps `optional steps`
	Finally build.Steps `optional steps`
}

func Try(context *build.Context, args interface{}) error {
	params := args.(TryArgs)
	context.SetProperty("_error", "")
	var tryError error
	var catchError error
	var finallyError error
	tryError = params.Try.Run(context)
	if tryError != nil {
		if len(params.Catch) > 0 || (len(params.Catch) == 0 && len(params.Finally) == 0) {
			context.SetProperty("_error", tryError.Error())
			tryError = nil
			catchError = params.Catch.Run(context)
		}
	}
	finallyError = params.Finally.Run(context)
	if finallyError != nil {
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
