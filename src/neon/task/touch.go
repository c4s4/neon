package task

import (
	"fmt"
	"io/ioutil"
	"neon/util"
	"os"
	"time"
)

func init() {
	TasksMap["touch"] = Descriptor{
		Constructor: Touch,
		Help:        "Touch a file (create it or change its time)",
	}
}

func Touch(target *Target, args util.Object) (Task, error) {
	fields := []string{"touch"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	files, err := args.GetListStringsOrString("touch")
	if err != nil {
		return nil, fmt.Errorf("argument to task touch must be a string or list of strings")
	}
	return func() error {
		fmt.Printf("Touching %d file(s)\n", len(files))
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
