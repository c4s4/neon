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
	return func() error {
		filename, err := target.Build.Context.ReplaceProperties(file)
		if err != nil {
			return fmt.Errorf("processing write argument: %v", err)
		}
		text, err := target.Build.Context.ReplaceProperties(source)
		if err != nil {
			return fmt.Errorf("processing text argument: %v", err)
		}
		var mode int
		if append {
			mode = os.O_CREATE | os.O_WRONLY | os.O_APPEND
		} else {
			mode = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
		}
		file, err := os.OpenFile(filename, mode, FILE_MODE)
		if err != nil {
			return fmt.Errorf("opening file '%s': %v", filename, err)
		}
		defer file.Close()
		_, err = file.WriteString(text)
		if err != nil {
			return fmt.Errorf("writing content to file '%s': %v", filename, err)
		}
		err = file.Sync()
		if err != nil {
			return fmt.Errorf("syncing content to file '%s': %v", filename, err)
		}
		return nil
	}, nil
}
