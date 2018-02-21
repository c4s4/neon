package task

import (
	"fmt"
	"neon/build"
	"os"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "write",
		Func: write,
		Args: reflect.TypeOf(writeArgs{}),
		Help: `Write text into given file.

Arguments:

- write: file to write into (string, file).
- text: text to write into the file (string, optional).
- append: tells if we should append content to file, default to false (boolean,
  optional).

Examples:

    # write 'Hello World!' in file greetings.txt
    - write: 'greetings.txt'
      text:  'Hello World!'`,
	})
}

type writeArgs struct {
	Write  string `file`
	Text   string `optional`
	Append bool   `optional`
}

func write(context *build.Context, args interface{}) error {
	params := args.(writeArgs)
	var mode int
	if params.Append {
		mode = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	} else {
		mode = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}
	file, err := os.OpenFile(params.Write, mode, FileMode)
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
