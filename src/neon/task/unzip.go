package task

import (
	"archive/zip"
	"fmt"
	"io"
	"neon/build"
	"neon/util"
	"os"
	"path/filepath"
)

func init() {
	build.TaskMap["unzip"] = build.TaskDescriptor{
		Constructor: Unzip,
		Help: `Expand a zip file in a directory.

Arguments:

- unzip: the zip file to expand.
- todir: the destination directory.

Examples:

    # unzip foo.zip to build directory
    - untar: "foo.zip"
      todir: "build"`,
	}
}

func Unzip(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"unzip", "todir"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	file, err := args.GetString("unzip")
	if err != nil {
		return nil, fmt.Errorf("argument unzip must be a string")
	}
	todir, err := args.GetString("todir")
	if err != nil {
		return nil, fmt.Errorf("argument todir of task unzip must be a string")
	}
	return func(context *build.Context) error {
		// evaluate arguments
		var _err error
		_file, _err := context.EvaluateString(file)
		if _err != nil {
			return fmt.Errorf("evaluating source zip file: %v", _err)
		}
		_todir, _err := context.EvaluateString(todir)
		if _err != nil {
			return fmt.Errorf("evaluating destination directory: %v", _err)
		}
		_file = util.ExpandUserHome(_file)
		_todir = util.ExpandUserHome(_todir)
		context.Message("Unzipping archive '%s' to directory '%s'...", _file, _todir)
		_err = UnzipFile(_file, _todir)
		if _err != nil {
			return fmt.Errorf("expanding archive: %v", _err)
		}
		return nil
	}, nil
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
					return fmt.Errorf("creating destination director %s: %v", target, err)
				}
			}
			dest, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return fmt.Errorf("creating destination file %s: %v", target, err)
			}
			defer dest.Close()
			_, err = io.Copy(dest, readCloser)
			if err != nil {
				return fmt.Errorf("copying to destination file %s: %v", target, err)
			}
		}
	}
	return nil
}
