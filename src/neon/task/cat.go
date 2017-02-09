package task

import (
	"fmt"
	"io/ioutil"
	"neon/build"
	"neon/util"
)

func init() {
	build.TaskMap["cat"] = build.TaskDescriptor{
		Constructor: Cat,
		Help: `Print the content of e given file on the console.

Arguments:

- cat: the file to print on console as a string.

Examples:

    # print content of LICENSE file on the console
    - cat: "LICENSE"`,
	}
}

func Cat(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"cat"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	file, ok := args["cat"].(string)
	if !ok {
		return nil, fmt.Errorf("argument of task cat must be a string")
	}
	return func() error {
		eval, err := target.Build.Context.ReplaceProperties(file)
		if err != nil {
			return fmt.Errorf("processing cat argument: %v", err)
		}
		content, err := ioutil.ReadFile(eval)
		if err != nil {
			return fmt.Errorf("printing content of file '%s': %v", eval, err)
		}
		target.Build.Info(string(content))
		return nil
	}, nil
}
