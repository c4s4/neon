package task

import (
	"fmt"
	"neon/build"
	"os"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc {
		Name: "write",
		Func: Write,
		Args: reflect.TypeOf(WriteArgs{}),
		Help: `Write text into a given file.

Arguments:

- write: the file to write into as a string.
- text: the text to write into the file.
- append: tells if we should append content to file (defaults to false).

Examples:

    # write 'Hello World!' in file greetings.txt
    - write: "greetings.txt"
      text: "Hello World!"`,
	})
}

type WriteArgs struct {
	Write  string `file`
	Text   string `optional`
	Append bool
}

func Write(context *build.Context, args interface{}) error {
	params := args.(WriteArgs)
	var mode int
	if params.Append {
		mode = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	} else {
		mode = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}
	file, err := os.OpenFile(params.Write, mode, FILE_MODE)
	if err != nil {
		return fmt.Errorf("opening file '%s': %v", params.Write, err)
	}
	defer file.Close()
	_, err = file.WriteString(params.Text)
	if err != nil {
		return fmt.Errorf("writing content to file '%s': %v", params.Write, err)
	}
	err = file.Sync()
	if err != nil {
		return fmt.Errorf("syncing content to file '%s': %v", params.Write, err)
	}
	return nil
}
