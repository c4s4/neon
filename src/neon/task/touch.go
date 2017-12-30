package task

import (
	"fmt"
	"io/ioutil"
	"neon/build"
	"neon/util"
	"os"
	"time"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc {
		Name: "touch",
		Func: Touch,
		Args: reflect.TypeOf(TouchArgs{}),
		Help: `Touch a file (create it or change its time).

Arguments:

- touch: the file or list of files to touch.

Examples:

    # create file in build directory
    - touch: ["#{BUILD_DIR}/foo", "#{BUILD_DIR}/bar"]

Notes:

- If the file already exists it changes it modification time.
- If the file doesn't exist, it creates an empty file.`,
	})
}

type TouchArgs struct {
	Touch []string `file wrap`
}

func Touch(context *build.Context, args interface{}) error {
	params := args.(TouchArgs)
	// DEBUG
	fmt.Printf("%#v\n", params)
	for _, file := range params.Touch {
		if util.FileExists(file) {
			time := time.Now()
			err := os.Chtimes(file, time, time)
			if err != nil {
				return fmt.Errorf("changing times of file '%s': %v", file, err)
			}
		} else {
			err := ioutil.WriteFile(file, []byte{}, FILE_MODE)
			if err != nil {
				return fmt.Errorf("creating file '%s': %v", file, err)
			}
		}
	}
	return nil
}
