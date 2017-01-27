package task

import (
	"archive/zip"
	"fmt"
	"io"
	"neon/build"
	"neon/util"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func init() {
	build.TaskMap["zip"] = build.TaskDescriptor{
		Constructor: Zip,
		Help: `Create a Zip archive.

Arguments:
- zip: the list of globs of files to zip (as a string or list of strings).
- dir: the root directory for glob (as a string, optional).
- exclude: globs of files to exclude (as a string or list of strings,
  optional).
- to: the name of the Zip file to create as a string.
- prefix: prefix directory in the archive.

Examples:
# zip files in build directory in file named build.zip
- zip: "build/**/*"
  to: "build.zip"`,
	}
}

func Zip(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"zip", "to", "dir", "exclude", "prefix"}
	if err := CheckFields(args, fields, fields[:2]); err != nil {
		return nil, err
	}
	includes, err := args.GetListStringsOrString("zip")
	if err != nil {
		return nil, fmt.Errorf("argument zip must be a string or list of strings")
	}
	var to string
	if args.HasField("to") {
		to, err = args.GetString("to")
		if err != nil {
			return nil, fmt.Errorf("argument to of task zip must be a string")
		}
	}
	var dir string
	if args.HasField("dir") {
		dir, err = args.GetString("dir")
		if err != nil {
			return nil, fmt.Errorf("argument dir of task zip must be a string")
		}
	}
	var excludes []string
	if args.HasField("exclude") {
		excludes, err = args.GetListStringsOrString("exclude")
		if err != nil {
			return nil, fmt.Errorf("argument exclude must be string or list of strings")
		}
	}
	var prefix string
	if args.HasField("prefix") {
		prefix, err = args.GetString("prefix")
		if err != nil {
			return nil, fmt.Errorf("argument prefix of task zip must be a string")
		}
	}
	return func() error {
		// evaluate arguments
		for index, pattern := range includes {
			eval, err := target.Build.Context.ReplaceProperties(pattern)
			if err != nil {
				return fmt.Errorf("evaluating pattern: %v", err)
			}
			includes[index] = eval
		}
		eval, err := target.Build.Context.ReplaceProperties(to)
		if err != nil {
			return fmt.Errorf("evaluating destination file: %v", err)
		}
		to = eval
		eval, err = target.Build.Context.ReplaceProperties(dir)
		if err != nil {
			return fmt.Errorf("evaluating source directory: %v", err)
		}
		dir = eval
		eval, err = target.Build.Context.ReplaceProperties(prefix)
		if err != nil {
			return fmt.Errorf("evaluating destination file: %v", err)
		}
		prefix = eval
		// find source files
		files, err := target.Build.Context.FindFiles(dir, includes, excludes)
		if err != nil {
			return fmt.Errorf("getting source files for zip task: %v", err)
		}
		if len(files) > 0 {
			target.Build.Info("Zipping %d file(s)", len(files))
			// zip files
			err = WriteZip(files, prefix, to)
			if err != nil {
				return fmt.Errorf("zipping files: %v", err)
			}
		}
		return nil
	}, nil
}

func WriteZip(files []string, prefix, to string) error {
	archive, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("creating zip archive: %v", err)
	}
	defer archive.Close()
	zipper := zip.NewWriter(archive)
	defer zipper.Close()
	for _, file := range files {
		err := writeFileToZip(zipper, file, prefix)
		if err != nil {
			return fmt.Errorf("writing file to zip archive: %v", err)
		}
	}
	return nil
}

func writeFileToZip(zipper *zip.Writer, filename, prefix string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	name := sanitizedName(filename)
	if prefix != "" {
		name = prefix + "/" + name
	}
	header.Name = name
	writer, err := zipper.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}

func sanitizedName(filename string) string {
	if len(filename) > 1 && filename[1] == ':' &&
		runtime.GOOS == "windows" {
		filename = filename[2:]
	}
	filename = filepath.ToSlash(filename)
	filename = strings.TrimLeft(filename, "/.")
	return strings.Replace(filename, "../", "", -1)
}
