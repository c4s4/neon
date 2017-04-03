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
- from: the name of the variable with the text to write (optional).
- text: the text to write into the file (optional).
- append: tells if we should append content to file (defaults to false).

Examples:

    # write 'Hello World!' in file greetings.txt
    - write: "greetings.txt"
      text: "Hello World!"`,
	}
}

func Write(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"write", "from", "text", "append"}
	if err := CheckFields(args, fields, fields[:1]); err != nil {
		return nil, err
	}
	file, err := args.GetString("write")
	if err != nil {
		return nil, fmt.Errorf("argument of task write must be a string")
	}
	var from string
	if args.HasField("from") {
		from, err = args.GetString("from")
		if err != nil {
			return nil, fmt.Errorf("argument from of task write must be a string")
		}
	}
	var text string
	if args.HasField("text") {
		text, err = args.GetString("text")
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
	if from != "" && text != "" {
		return nil, fmt.Errorf("you can't set both from and test arguments for task write")
	}
	return func() error {
		eval, err := target.Build.Context.ReplaceProperties(file)
		if err != nil {
			return fmt.Errorf("processing write argument: %v", err)
		}
		if from != "" {
			object, err := target.Build.Context.GetProperty(from)
			if err != nil {
				return fmt.Errorf("getting variable '%s': %v", from, err)
			}
			var ok bool
			text, ok = object.(string)
			if !ok {
				return fmt.Errorf("variable in argument from of task write must be of type string")
			}
		}
		var mode int
		if append {
			mode = os.O_CREATE | os.O_WRONLY | os.O_APPEND
		} else {
			mode = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
		}
		file, err := os.OpenFile(eval, mode, FILE_MODE)
		if err != nil {
			return fmt.Errorf("opening file '%s': %v", eval, err)
		}
		defer file.Close()
		_, err = file.WriteString(text)
		if err != nil {
			return fmt.Errorf("writing content to file '%s': %v", eval, err)
		}
		err = file.Sync()
		if err != nil {
			return fmt.Errorf("syncing content to file '%s': %v", eval, err)
		}
		return nil
	}, nil
}
