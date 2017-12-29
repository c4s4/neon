// +build ignore

package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
)

func init() {
	build.TaskMap["link"] = build.TaskDescriptor{
		Constructor: Link,
		Help: `Create a symbolic link.

Arguments:

- link: the source file.
- to: the destination of the link.

Examples:

    # create a link from file foo to bar
    - link: "foo"
      to: "bar"`,
	}
}

func Link(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"link", "to"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	s, err := args.GetString("link")
	if err != nil {
		return nil, fmt.Errorf("argument link must be a string")
	}
	d, err := args.GetString("to")
	if err != nil {
		return nil, fmt.Errorf("argument to of task link must be a string")
	}
	return func(context *build.Context) error {
		_source, _err := context.EvaluateString(s)
		if _err != nil {
			return fmt.Errorf("processing link argument: %v", _err)
		}
		_dest, _err := context.EvaluateString(d)
		if _err != nil {
			return fmt.Errorf("processing to argument of link task: %v", _err)
		}
		context.Message("Linking file '%s' to '%s'", _source, _dest)
		_err = os.Symlink(_source, _dest)
		if _err != nil {
			return fmt.Errorf("linking files: %v", _err)
		}
		return nil
	}, nil
}
