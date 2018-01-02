package task

import (
	"archive/zip"
	"fmt"
	"io"
	"neon/build"
	"os"
	"path/filepath"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "unzip",
		Func: Unzip,
		Args: reflect.TypeOf(UnzipArgs{}),
		Help: `Expand a zip file in a directory.

Arguments:

- unzip: the zip file to expand.
- todir: the destination directory.

Examples:

    # unzip foo.zip to build directory
    - unzip: "foo.zip"
      todir: "build"`,
	})
}

type UnzipArgs struct {
	Unzip string `file`
	Todir string `file`
}

func Unzip(context *build.Context, args interface{}) error {
	params := args.(UnzipArgs)
	context.Message("Unzipping archive '%s' to directory '%s'...", params.Unzip, params.Todir)
	err := UnzipFile(params.Unzip, params.Todir)
	if err != nil {
		return fmt.Errorf("expanding archive: %v", err)
	}
	return nil
}

// Unzip given file to a directory
func UnzipFile(file, dir string) error {
	reader, err := zip.OpenReader(file)
	if err != nil {
		return fmt.Errorf("opening source zip file %s: %v", file, err)
	}
	defer reader.Close()
	for _, file := range reader.File {
		readCloser, err := file.Open()
		if err != nil {
			return fmt.Errorf("opening zip file %s: %v", file.Name, err)
		}
		defer readCloser.Close()
		target := filepath.Join(dir, file.Name)
		if file.FileInfo().IsDir() {
			continue
		} else {
			destination := filepath.Dir(target)
			if _, err := os.Stat(destination); err != nil {
				if err := os.MkdirAll(destination, 0755); err != nil {
					return fmt.Errorf("creating destination directory '%s': %v", target, err)
				}
			}
			dest, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return fmt.Errorf("creating destination file '%s': %v", target, err)
			}
			defer dest.Close()
			_, err = io.Copy(dest, readCloser)
			if err != nil {
				return fmt.Errorf("copying to destination file '%s': %v", target, err)
			}
		}
	}
	return nil
}
