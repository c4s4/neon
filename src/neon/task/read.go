package task

import (
	"fmt"
	"io/ioutil"
	"neon/build"
	"neon/util"
)

func init() {
	build.TaskMap["read"] = build.TaskDescriptor{
		Constructor: Read,
		Help: `Read given file as text and put its content in a variable.

Arguments:

- read: the file to read as a string.
- to: the name of the variable to set with the content.

Examples:

    # put content of LICENSE file on license variable
    - read: "LICENSE"
      to: license`,
	}
}

func Read(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"read", "to"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	file, err := args.GetString("read")
	if err != nil {
		return nil, fmt.Errorf("argument of task read must be a string")
	}
	to, err := args.GetString("to")
	if err != nil {
		return nil, fmt.Errorf("argument to of task read must be a string")
	}
	return func() error {
		eval, err := target.Build.Context.EvaluateString(file)
		if err != nil {
			return fmt.Errorf("processing read argument: %v", err)
		}
		content, err := ioutil.ReadFile(eval)
		if err != nil {
			return fmt.Errorf("reading content of file '%s': %v", eval, err)
		}
		target.Build.Context.SetProperty(to, string(content))
		return nil
	}, nil
}
