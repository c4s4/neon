package task

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"neon/build"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "untar",
		Func: Untar,
		Args: reflect.TypeOf(UntarArgs{}),
		Help: `Expand a tar file in a directory.

Arguments:

- untar: the tar file to expand (string, file).
- todir: the destination directory (string, file).

Examples:

    # untar foo.tar to build directory
    - untar: 'foo.tar'
      todir: 'build'

Notes:

- If archive filename ends with .gz (with a name such as foo.tar.gz or foo.tgz)
  the tar archive is uncompressed with gzip.`,
	})
}

type UntarArgs struct {
	Untar string `file`
	Todir string `file`
}

func Untar(context *build.Context, args interface{}) error {
	params := args.(UntarArgs)
	context.Message("Untarring archive '%s' to directory '%s'...", params.Untar, params.Todir)
	err := UntarFile(params.Untar, params.Todir)
	if err != nil {
		return fmt.Errorf("expanding archive: %v", err)
	}
	return nil
}

func UntarFile(file, dir string) error {
	reader, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("opening source tar file %s: %v", file, err)
	}
	defer reader.Close()
	var tarReader *tar.Reader
	if strings.HasSuffix(file, ".gz") || strings.HasSuffix(file, ".tgz") {
		gzipReader, err := gzip.NewReader(reader)
		defer gzipReader.Close()
		if err != nil {
			return fmt.Errorf("unzipping tar file: %v", err)
		}
		tarReader = tar.NewReader(gzipReader)
	} else {
		tarReader = tar.NewReader(reader)
	}
	for {
		header, err := tarReader.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}
		target := filepath.Join(dir, header.Name)
		if header.Typeflag == tar.TypeReg {
			destination := filepath.Dir(target)
			if _, err := os.Stat(destination); err != nil {
				if err := os.MkdirAll(destination, 0755); err != nil {
					return fmt.Errorf("creating destination director %s: %v", target, err)
				}
			}
			dest, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("creating destination file %s: %v", target, err)
			}
			if _, err := io.Copy(dest, tarReader); err != nil {
				return fmt.Errorf("copying to destination file %s: %v", target, err)
			}
			dest.Close()
		}
	}
}
