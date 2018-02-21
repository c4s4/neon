package task

import (
	z "archive/zip"
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
		Func: unzip,
		Args: reflect.TypeOf(unzipArgs{}),
		Help: `Expand a zip file in a directory.

Arguments:

- unzip: the zip file to expand (string, file).
- todir: the destination directory (string, file).

Examples:

    # unzip foo.zip to build directory
    - unzip: 'foo.zip'
      todir: 'build'`,
	})
}

type unzipArgs struct {
	Unzip string `file`
	Todir string `file`
}

func unzip(context *build.Context, args interface{}) error {
	params := args.(unzipArgs)
	context.Message("Unzipping archive '%s' to directory '%s'...", params.Unzip, params.Todir)
	err := unzipFile(params.Unzip, params.Todir)
	if err != nil {
		return fmt.Errorf("expanding archive: %v", err)
	}
	return nil
}

func unzipFile(file, dir string) error {
	reader, err := z.OpenReader(file)
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
