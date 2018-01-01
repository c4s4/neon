package task

import (
	"bufio"
	"fmt"
	"neon/build"
	"os"
	"regexp"
	"strings"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc {
		Name: "prompt",
		Func: Prompt,
		Args: reflect.TypeOf(PromptArgs{}),
		Help: `Prompt the user for the value of a given property matching a pattern.

Arguments:

- prompt: message to print at prompt. Should include a description of the
  expected pattern.
- property: the name of the property to set.
- default: default value if user doesn't type anything. Written into square
  brackets after prompt message. Optional.
- pattern: a regular expression for prompted value. If this pattern is not
  matched, this task will prompt again. Optional, if no pattern is given, any
  value is accepted.
- error: the error message to print when pattern is not matched.

Examples:

    # returns: typed message
    - prompt:  "Enter your age"
      to:      "age"
      default: "18"
      pattern: "^\d+\s$"
      error:   "Age must be a positive integer"`,
	})
}

type PromptArgs struct {
	Prompt   string
	Property string
	Default  string `optional`
	Pattern  string `optional`
	Error    string `optional`
}

func Prompt(context *build.Context, args interface{}) error {
	params := args.(PromptArgs)
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
			context.SetProperty(params.Property, string(value))
		}
	}
	return nil
}
