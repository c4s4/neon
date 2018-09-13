package task

import (
	"fmt"
	"io/ioutil"
	"neon/build"
	"neon/util"
	"os"
	"reflect"
	t "time"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "touch",
		Func: touch,
		Args: reflect.TypeOf(touchArgs{}),
		Help: `Touch a file (create it or change its time).

Arguments:

- touch: files to touch (strings, file, wrap).

Examples:

    # create file in build directory
    - touch: ['#{BUILD_DIR}/foo', '#{BUILD_DIR}/bar']

Notes:

- If the file already exists it changes it modification time.
- If the file doesn't exist, it creates an empty file.`,
	})
}

type touchArgs struct {
	Touch []string `neon:"file,wrap"`
}

func touch(context *build.Context, args interface{}) error {
	params := args.(touchArgs)
	context.Message("Touching %d file(s)", len(params.Touch))
	for _, file := range params.Touch {
		if util.FileExists(file) {
			time := t.Now()
			err := os.Chtimes(file, time, time)
			if err != nil {
				return fmt.Errorf("changing times of file '%s': %v", file, err)
			}
		} else {
			err := ioutil.WriteFile(file, []byte{}, FileMode)
			if err != nil {
				return fmt.Errorf("creating file '%s': %v", file, err)
			}
		}
	}
	return nil
}
