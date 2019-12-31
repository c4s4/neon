package task

import (
	z "archive/zip"
	"compress/flate"
	"fmt"
	"github.com/c4s4/neon/neon/build"
	"github.com/c4s4/neon/neon/util"
	"io"
	"os"
	"path/filepath"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "zip",
		Func: zip,
		Args: reflect.TypeOf(zipArgs{}),
		Help: `Create a Zip archive.

Arguments:

- zip: globs of files to zip (strings, file, wrap).
- dir: root directory for globs, defaults to '.' (string, optional, file).
- exclude: globs of files to exclude (strings, optional, file, wrap).
- tofile: name of the Zip file to create (string, file).
- prefix: prefix directory in the archive (string, optional).

Examples:

    # zip files of build directory in file named build.zip
    - zip:    'build/**/*'
      tofile: 'build.zip'`,
	})
}

type zipArgs struct {
	Zip     []string `neon:"file,wrap"`
	Dir     string   `neon:"optional,file"`
	Exclude []string `neon:"optional,file,wrap"`
	Tofile  string   `neon:"file"`
	Prefix  string   `neon:"optional"`
}

func zip(context *build.Context, args interface{}) error {
	params := args.(zipArgs)
	files, err := util.FindFiles(params.Dir, params.Zip, params.Exclude, false)
	if err != nil {
		return fmt.Errorf("getting source files for zip task: %v", err)
	}
	if len(files) > 0 {
		context.Message("Zipping %d file(s) in '%s'", len(files), params.Tofile)
		err = writeZip(params.Dir, files, params.Prefix, params.Tofile)
		if err != nil {
			return fmt.Errorf("zipping files: %v", err)
		}
	}
	return nil
}

func writeZip(dir string, files []string, prefix, to string) error {
	archive, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("creating zip archive: %v", err)
	}
	defer archive.Close()
	zipper := z.NewWriter(archive)
	zipper.RegisterCompressor(z.Deflate, func(out io.Writer) (io.WriteCloser, error) {
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

func writeFileToZip(zipper *z.Writer, path, name, prefix string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := z.FileInfoHeader(info)
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
