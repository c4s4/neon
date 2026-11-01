package task

import (
	t "archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	p "path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/c4s4/neon/neon/build"
	"github.com/c4s4/neon/neon/util"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "tar",
		Func: tar,
		Args: reflect.TypeOf(tarArgs{}),
		Help: `Create a tar archive.

Arguments:

- tar: globs of files to tar (strings, file, wrap).
- dir: root directory for glob, defaults to '.' (string, optional, file).
- exclude: globs of files to exclude (strings, optional, file, wrap).
- tofile: name of the tar file to create (string, file).
- prefix: prefix directory in the archive (string, optional).

Examples:

    # tar files in build directory in file named build.tar.gz
    - tar:    'build/**/*'
      tofile: 'build.tar.gz'

Notes:

- If archive filename ends with gz (with names such as 'foo.tar.gz' or
  'foo.tgz') the tar archive is also gzip compressed.`,
	})
}

type tarArgs struct {
	Tar     []string `neon:"file,wrap"`
	Dir     string   `neon:"optional,file"`
	Exclude []string `neon:"optional,file,wrap"`
	Tofile  string   `neon:"file"`
	Prefix  string   `neon:"optional"`
}

func tar(context *build.Context, args interface{}) error {
	params := args.(tarArgs)
	files, err := util.FindFiles(params.Dir, params.Tar, params.Exclude, false)
	if err != nil {
		return fmt.Errorf("getting source files for tar task: %v", err)
	}
	if len(files) > 0 {
		context.MessageArgs("Tarring %d file(s) into '%s'", len(files), params.Tofile)
		err = writeTar(params.Dir, files, params.Prefix, params.Tofile)
		if err != nil {
			return fmt.Errorf("tarring files: %v", err)
		}
	}
	return nil
}

func writeTar(dir string, files []string, prefix, to string) (err error) {
	stream, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("creating tar archive: %v", err)
	}
	defer func() {
		if e := stream.Close(); e != nil && err == nil {
			err = fmt.Errorf("closing tar archive: %v", e)
		}
	}()
	var gzipWriter *gzip.Writer
	var fileWriter io.Writer = stream
	if strings.HasSuffix(to, "gz") {
		gzipWriter = gzip.NewWriter(stream)
		fileWriter = gzipWriter
	}
	writer := t.NewWriter(fileWriter)
	for _, name := range files {
		var file string
		if dir != "" {
			file = filepath.Join(dir, name)
		} else {
			file = name
		}
		if err := writeFileToTar(writer, file, name, prefix); err != nil {
			return fmt.Errorf("writing stream to tar archive: %v", err)
		}
	}
	// close writers in order to flush archive trailers
	if err := writer.Close(); err != nil {
		return fmt.Errorf("writing tar archive: %v", err)
	}
	if gzipWriter != nil {
		if err := gzipWriter.Close(); err != nil {
			return fmt.Errorf("compressing tar archive: %v", err)
		}
	}
	return nil
}

func writeFileToTar(writer *t.Writer, file, name, prefix string) error {
	stream, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() {
		_ = stream.Close()
	}()
	info, err := stream.Stat()
	if err != nil {
		return err
	}
	header, err := t.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = p.Join(prefix, SanitizeName(name))
	if err = writer.WriteHeader(header); err != nil {
		return err
	}
	_, err = io.Copy(writer, stream)
	return err
}
