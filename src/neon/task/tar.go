package task

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"neon/build"
	"neon/util"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "tar",
		Func: Tar,
		Args: reflect.TypeOf(TarArgs{}),
		Help: `Create a tar archive.

Arguments:

- tar: the list of globs of files to tar (as a string or list of strings).
- dir: the root directory for glob (as a string, optional).
- exclude: globs of files to exclude (as a string or list of strings,
  optional).
- tofile: the name of the tar file to create as a string.
- prefix: prefix directory in the archive.

Examples:

    # tar files in build directory in file named build.tar.gz
    - tar: "build/**/*"
      tofile: "build.tar.gz"

Notes:

- If archive filename ends with gz (with a name such as foo.tar.gz or foo.tgz)
  the tar archive is compressed with gzip.`,
	})
}

type TarArgs struct {
	Tar     []string `file wrap`
	Dir     string   `optional file`
	Exclude []string `optional file wrap`
	Tofile  string   `file`
	Prefix  string   `optional`
}

func Tar(context *build.Context, args interface{}) error {
	params := args.(TarArgs)
	files, err := util.FindFiles(params.Dir, params.Tar, params.Exclude, false)
	if err != nil {
		return fmt.Errorf("getting source files for tar task: %v", err)
	}
	if len(files) > 0 {
		context.Message("Tarring %d file(s) into '%s'", len(files), params.Tofile)
		err = Writetar(params.Dir, files, params.Prefix, params.Tofile)
		if err != nil {
			return fmt.Errorf("tarring files: %v", err)
		}
	}
	return nil
}

func Writetar(dir string, files []string, prefix, to string) error {
	stream, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("creating tar archive: %v", err)
	}
	defer stream.Close()
	var fileWriter io.WriteCloser = stream
	if strings.HasSuffix(to, "gz") {
		fileWriter = gzip.NewWriter(stream)
		defer fileWriter.Close()
	}
	writer := tar.NewWriter(fileWriter)
	defer writer.Close()
	for _, name := range files {
		var file string
		if dir != "" {
			file = filepath.Join(dir, name)
		} else {
			file = name
		}
		err := writeFileToTar(writer, file, name, prefix)
		if err != nil {
			return fmt.Errorf("writing stream to tar archive: %v", err)
		}
	}
	return nil
}

func writeFileToTar(writer *tar.Writer, file, name, prefix string) error {
	stream, err := os.Open(file)
	if err != nil {
		return err
	}
	defer stream.Close()
	info, err := stream.Stat()
	if err != nil {
		return err
	}
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = path.Join(prefix, SanitizeName(name))
	if err = writer.WriteHeader(header); err != nil {
		return err
	}
	_, err = io.Copy(writer, stream)
	return err
}
