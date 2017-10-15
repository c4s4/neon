package task

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"neon/build"
	"neon/util"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	build.TaskMap["untar"] = build.TaskDescriptor{
		Constructor: Untar,
		Help: `Expand a tar file in a directory.

Arguments:

- untar: the tar file to expand.
- todir: the destination directory.

Examples:

    # untar foo.tar to build directory
    - untar: "foo.tar"
      todir: "build"

Notes:

- If archive filename ends with gz (with a name such as foo.tar.gz or foo.tgz)
  the tar archive is uncompressed with gzip.`,
	}
}

func Untar(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"untar", "todir"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	file, err := args.GetString("untar")
	if err != nil {
		return nil, fmt.Errorf("argument untar must be a string")
	}
	todir, err := args.GetString("todir")
	if err != nil {
		return nil, fmt.Errorf("argument todir of task untar must be a string")
	}
	return func(context *build.Context) error {
		// evaluate arguments
		var _err error
		_file, _err := context.EvaluateString(file)
		if _err != nil {
			return fmt.Errorf("evaluating source tar file: %v", _err)
		}
		_todir, _err := context.EvaluateString(todir)
		if _err != nil {
			return fmt.Errorf("evaluating destination directory: %v", _err)
		}
		_file = util.ExpandUserHome(_file)
		_todir = util.ExpandUserHome(_todir)
		context.Message("Untarring archive '%s' to directory '%s'...", _file, _todir)
		_err = UntarFile(_file, _todir)
		if _err != nil {
			return fmt.Errorf("expanding archive: %v", _err)
		}
		return nil
	}, nil
}

// Untar given file to a directory
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
