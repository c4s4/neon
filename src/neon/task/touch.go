package task

import (
	"fmt"
	"io/ioutil"
	"neon/build"
	"neon/util"
	"os"
	"time"
)

func init() {
	build.TaskMap["touch"] = build.TaskDescriptor{
		Constructor: Touch,
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

func Touch(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"touch"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	files, err := args.GetListStringsOrString("touch")
	if err != nil {
		return nil, fmt.Errorf("argument to task touch must be a string or list of strings")
	}
	return func() error {
		build.Info("Touching %d file(s)", len(files))
		for _, file := range files {
			path, err := target.Build.Context.ReplaceProperties(file)
			if err != nil {
				return fmt.Errorf("processing touch argument: %v", err)
			}
			if util.FileExists(path) {
				time := time.Now()
				err = os.Chtimes(path, time, time)
				if err != nil {
					return fmt.Errorf("changing times of file '%s': %v", path, err)
				}
			} else {
				err := ioutil.WriteFile(path, []byte{}, FILE_MODE)
				if err != nil {
					return fmt.Errorf("creating file '%s': %v", path, err)
				}
			}
		}
		return nil
	}, nil
}
