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
	return func(context *build.Context) error {
		_eval, _err := context.EvaluateString(file)
		if _err != nil {
			return fmt.Errorf("processing cat argument: %v", _err)
		}
		_content, _err := ioutil.ReadFile(_eval)
		if _err != nil {
			return fmt.Errorf("printing _content of file '%s': %v", _eval, _err)
		}
		context.Message(string(_content))
		return nil
	}, nil
}
