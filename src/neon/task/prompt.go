package task

import (
	"bufio"
	"fmt"
	"neon/build"
	"neon/util"
	"os"
	"regexp"
	"strings"
)

func init() {
	build.TaskMap["prompt"] = build.TaskDescriptor{
		Constructor: Prompt,
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
	}
}

func Prompt(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"prompt", "to", "default", "pattern", "error"}
	if err := CheckFields(args, fields, fields[:2]); err != nil {
		return nil, err
	}
	message, err := args.GetString("prompt")
	if err != nil {
		return nil, fmt.Errorf("argument of task prompt must be a string")
	}
	to, err := args.GetString("to")
	if err != nil {
		return nil, fmt.Errorf("argument to of task prompt must be a string")
	}
	var def string
	if args.HasField("default") {
		def, err = args.GetString("default")
		if err != nil {
			return nil, fmt.Errorf("argument default of task prompt must be a string")
		}
	}
	var pattern string
	if args.HasField("pattern") {
		pattern, err = args.GetString("pattern")
		if err != nil {
			return nil, fmt.Errorf("argument pattern of task prompt must be a string")
		}
	}
	var errorMessage string
	if args.HasField("error") {
		errorMessage, err = args.GetString("error")
		if err != nil {
			return nil, fmt.Errorf("argument error of task prompt must be a string")
		}
	}
	return func(context *build.Context) error {
		_message, _err := context.EvaluateString(message)
		if _err != nil {
			return fmt.Errorf("processing prompt argument: %v", _err)
		}
		_to, _err := context.EvaluateString(to)
		if _err != nil {
			return fmt.Errorf("evaluating destination variable: %v", _err)
		}
		_default, _err := context.EvaluateString(def)
		if _err != nil {
			return fmt.Errorf("evaluating default value: %v", _err)
		}
		_pattern, _err := context.EvaluateString(pattern)
		if _err != nil {
			return fmt.Errorf("evaluating input regular expression: %v", _err)
		}
		_errorMessage, _err := context.EvaluateString(errorMessage)
		if _err != nil {
			return fmt.Errorf("evaluating error message: %v", _err)
		}
		if _default != "" {
			_message += " [" + _default + "]"
		}
		_message += ": "
		done := false
		for !done {
			fmt.Print(_message)
			_value, _err := bufio.NewReader(os.Stdin).ReadString('\n')
			if _err != nil {
				return fmt.Errorf("reading user input: %v", _err)
			}
			_value = strings.TrimSpace(_value)
			if _value == "" && _default != "" {
				_value = _default
			}
			if pattern != "" && !regexp.MustCompile(_pattern).MatchString(_value) {
				if _errorMessage != "" {
					context.Message(_errorMessage)
				} else {
					context.Message("value '%s' doesn't match pattern '%s'", _value, _pattern)
				}
			} else {
				done = true
				context.SetProperty(_to, string(_value))
			}
		}
		return nil
	}, nil
}
