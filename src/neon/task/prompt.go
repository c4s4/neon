package task

import (
	"bufio"
	"fmt"
	"neon/build"
	"os"
	"reflect"
	"regexp"
	"strings"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "prompt",
		Func: prompt,
		Args: reflect.TypeOf(promptArgs{}),
		Help: `Prompt the user for the value of a given property matching a pattern.

Arguments:

- prompt: message to print at prompt that include a description of expected
  pattern (string).
- to: name of the property to set (string).
- default: default value if user doesn't type anything, written into square
  brackets after prompt message (string, optional).
- pattern: a regular expression for prompted value. If this pattern is not
  matched, this task will prompt again. If no pattern is given, any value is
  accepted (string, optional).
- error: error message to print when pattern is not matched (string, optional).

Examples:

    # prompt for age that is a positive number
    - prompt:  'Enter your age'
      to:      'age'
      default: '18'
      pattern: '^\d+$'
      error:   'Age must be a positive integer'`,
	})
}

type promptArgs struct {
	Prompt  string
	To      string
	Default string `optional`
	Pattern string `optional`
	Error   string `optional`
}

func prompt(context *build.Context, args interface{}) error {
	params := args.(promptArgs)
	message := params.Prompt
	if params.Default != "" {
		message += " [" + params.Default + "]"
	}
	message += ": "
	done := false
	for !done {
		fmt.Print(message)
		value, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading user input: %v", err)
		}
		value = strings.TrimSpace(value)
		if value == "" && params.Default != "" {
			value = params.Default
		}
		if params.Pattern != "" && !regexp.MustCompile(params.Pattern).MatchString(value) {
			if params.Error != "" {
				context.Message(params.Error)
			} else {
				context.Message("value '%s' doesn't match pattern '%s'", value, params.Pattern)
			}
		} else {
			done = true
			context.SetProperty(params.To, string(value))
		}
	}
	return nil
}
