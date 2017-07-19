package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"bufio"
	"os"
	"regexp"
)

func init() {
	build.TaskMap["prompt"] = build.TaskDescriptor{
		Constructor: Prompt,
		Help: `Prompt the user for the value of a given property matching a pattern.

Arguments:

- message: message to print at prompt. Should include a description
  of the expected pattern.
- property: the name of the property to set.
- default: default value if user doesn't type anything. Written
  into square brakets after prompt message. Optional.
- pattern: a regular expression for prompted value. If this pattern is not
  matched, this task will prompt again. Optional, if no pattern is
  given, any value is accepted.
- error: the error message to print when pattern is not matched.

		Example
		 # returns: typed message
		 - prompt:		 "Enter your age"
			 to:		   "x"
			 defaut:  "18"
			 pattern:  "^\d+\s$"
			 error:		"Age must be a positive integer"
		`,
	}
}

func Prompt(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"prompt", "to", "default", "pattern", "error"}
	if err := CheckFields(args, fields, fields); err != nil {
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
	defaut, err := args.GetString("default")
	if err != nil {
		return nil, fmt.Errorf("argument default of task prompt must be a string")
	}
	pattern, err := args.GetString("pattern")
	if err != nil {
		return nil, fmt.Errorf("argument pattern of task prompt must be a string")
	}
	erro, err := args.GetString("error")
	if err != nil {
		return nil, fmt.Errorf("argument error of task prompt must be a string")
	}
	return func() error {
		_message, _err := target.Build.Context.EvaluateString(message)
		if _err != nil {
			return fmt.Errorf("processing prompt argument: %v", _err)
		}
		
		_eval, _err := target.Build.Context.EvaluateString(to)
		if _err != nil {
			return fmt.Errorf("evaluating destination variable: %v", _err)
		}
		_to := _eval
		_eval, _err = target.Build.Context.EvaluateString(defaut)
		if _err != nil {
			return fmt.Errorf("evaluating default value: %v", _err)
		}
		_defaut := _eval
		_eval, _err = target.Build.Context.EvaluateString(pattern)
		if _err != nil {
			return fmt.Errorf("evaluating input regular expression: %v", _err)
		}
		_pattern := _eval
		_eval, _err = target.Build.Context.EvaluateString(erro)
		if _err != nil {
			return fmt.Errorf("evaluating error message: %v", _err)
		}
		_error := _eval
		
		s := _message
		if defaut != "" {
		  s += " ["+_defaut+"]"
		}
		s += ": "
		fmt.Println(s)
		reader := bufio.NewReader(os.Stdin)
		value, _err := reader.ReadString('\n')
		if _err != nil {
			return fmt.Errorf("reading stdin '%s': %v", value, _err)
		}
		if value == "\n" {
		  target.Build.Context.SetProperty(to, string(_defaut))
		}
		if pattern != "" && !regexp.MustCompile(_pattern).MatchString(value) {
		  fmt.Println(_error)
		  target.Build.Context.SetProperty(_to, string(_defaut))
		  return fmt.Errorf("User Input does not match %s, err : %v", _pattern, _error)
		}
 
		target.Build.Context.SetProperty(to, string(value))
		return nil
	}, nil
}
