// +build ignore

package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
)

func init() {
	build.TaskMap["write"] = build.TaskDescriptor{
		Constructor: Write,
		Help: `Write text into a given file.

Arguments:

- write: the file to write into as a string.
- text: the text to write into the file.
- append: tells if we should append content to file (defaults to false).

Examples:

    # write 'Hello World!' in file greetings.txt
    - write: "greetings.txt"
      text: "Hello World!"`,
	}
}

func Write(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"write", "text", "append"}
	if err := CheckFields(args, fields, fields[:2]); err != nil {
		return nil, err
	}
	file, err := args.GetString("write")
	if err != nil {
		return nil, fmt.Errorf("argument of task write must be a string")
	}
	var source string
	if args.HasField("text") {
		source, err = args.GetString("text")
		if err != nil {
			return nil, fmt.Errorf("argument text of task write must be a string")
		}
	}
	append := false
	if args.HasField("append") {
		append, err = args.GetBoolean("append")
		if err != nil {
			return nil, fmt.Errorf("argument append of task write must be a boolean")
		}
	}
	return func(context *build.Context) error {
		_filename, _err := context.EvaluateString(file)
		if _err != nil {
			return fmt.Errorf("processing write argument: %v", _err)
		}
		_text, _err := context.EvaluateString(source)
		if _err != nil {
			return fmt.Errorf("processing text argument: %v", _err)
		}
		var _mode int
		if append {
			_mode = os.O_CREATE | os.O_WRONLY | os.O_APPEND
		} else {
			_mode = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
		}
		_file, _err := os.OpenFile(_filename, _mode, FILE_MODE)
		if _err != nil {
			return fmt.Errorf("opening file '%s': %v", _filename, _err)
		}
		defer _file.Close()
		_, _err = _file.WriteString(_text)
		if _err != nil {
			return fmt.Errorf("writing content to file '%s': %v", _filename, _err)
		}
		_err = _file.Sync()
		if _err != nil {
			return fmt.Errorf("syncing content to file '%s': %v", _filename, _err)
		}
		return nil
	}, nil
}
