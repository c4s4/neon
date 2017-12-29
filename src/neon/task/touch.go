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
	build.TaskMap["touch"] = build.TaskDesc {
		Func: Touch,
		Args: reflect.TypeOf(TouchArgs{}),
		Help: `Touch a file (create it or change its time).

Arguments:

- touch: the file or files to create.

Examples:

    # create file in build directory
    - touch: "#{BUILD_DIR}/foo"

Notes:

- If the file already exists it changes it modification time.
- If the file doesn't exist, it creates an empty file.`,
	}
}

type TouchArgs struct {
	Touch string `file`
}

func Touch(context *build.Context, args interface{}) error {
	params := args.(TouchArgs)
	if util.FileExists(params.Touch) {
		time := time.Now()
		err := os.Chtimes(params.Touch, time, time)
		if err != nil {
			return fmt.Errorf("changing times of file '%s': %v", params.Touch, err)
		}
	} else {
		err := ioutil.WriteFile(params.Touch, []byte{}, FILE_MODE)
		if err != nil {
			return fmt.Errorf("creating file '%s': %v", params.Touch, err)
		}
	}
	return nil
}
