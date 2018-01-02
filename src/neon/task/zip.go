package task

import (
	"archive/zip"
	"compress/flate"
	"fmt"
	"io"
	"neon/build"
	"neon/util"
	"os"
	"path/filepath"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "zip",
		Func: Zip,
		Args: reflect.TypeOf(ZipArgs{}),
		Help: `Create a Zip archive.

Arguments:

- zip: the list of globs of files to zip (as a string or list of strings).
- dir: the root directory for glob (as a string, optional).
- exclude: globs of files to exclude (as a string or list of strings,
  optional).
- tofile: the name of the Zip file to create as a string.
- prefix: prefix directory in the archive.

Examples:

    # zip files in build directory in file named build.zip
    - zip: "build/**/*"
      tofile: "build.zip"`,
	})
}

type ZipArgs struct {
	Zip     []string `file wrap`
	Dir     string   `optional file`
	Exclude []string `optional file wrap`
	Tofile  string   `optional file`
	Prefix  string   `optional`
}

func Zip(context *build.Context, args interface{}) error {
	params := args.(ZipArgs)
	files, err := util.FindFiles(params.Dir, params.Zip, params.Exclude, false)
	if err != nil {
		return fmt.Errorf("getting source files for zip task: %v", err)
	}
	if len(files) > 0 {
		context.Message("Zipping %d file(s) in '%s'", len(files), params.Tofile)
		err = WriteZip(params.Dir, files, params.Prefix, params.Tofile)
		if err != nil {
			return fmt.Errorf("zipping files: %v", err)
		}
	}
	return nil
}

func WriteZip(dir string, files []string, prefix, to string) error {
	archive, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("creating zip archive: %v", err)
	}
	defer archive.Close()
	zipper := zip.NewWriter(archive)
	zipper.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})
	defer zipper.Close()
	for _, file := range files {
		var path string
		if dir != "" {
			path = filepath.Join(dir, file)
		} else {
			path = file
		}
		err := writeFileToZip(zipper, path, file, prefix)
		if err != nil {
			return fmt.Errorf("writing file to zip archive: %v", err)
		}
	}
	return nil
}

func writeFileToZip(zipper *zip.Writer, path, name, prefix string) error {
	file, err := os.Open(path)
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
	name = SanitizeName(name)
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
