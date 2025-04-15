package task

import (
	t "archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/c4s4/neon/neon/build"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "untar",
		Func: untar,
		Args: reflect.TypeOf(untarArgs{}),
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

type untarArgs struct {
	Untar string `neon:"file"`
	Todir string `neon:"file"`
}

func untar(context *build.Context, args interface{}) error {
	params := args.(untarArgs)
	context.MessageArgs("Untarring archive '%s' to directory '%s'...", params.Untar, params.Todir)
	err := untarFile(params.Untar, params.Todir)
	if err != nil {
		return fmt.Errorf("expanding archive: %v", err)
	}
	return nil
}

func untarFile(file, dir string) error {
	reader, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("opening source tar file %s: %v", file, err)
	}
	defer func() {
		_ = reader.Close()
	}()
	var tarReader *t.Reader
	if strings.HasSuffix(file, ".gz") || strings.HasSuffix(file, ".tgz") {
		gzipReader, err := gzip.NewReader(reader)
		if err != nil {
			return fmt.Errorf("unzipping tar file: %v", err)
		}
		defer func() {
			_ = gzipReader.Close()
		}()
		tarReader = t.NewReader(gzipReader)
	} else {
		tarReader = t.NewReader(reader)
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
		if header.Typeflag == t.TypeReg {
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
			if err := dest.Close(); err != nil {
				return fmt.Errorf("closing destination file %s: %v", target, err)
			}
		}
	}
}
