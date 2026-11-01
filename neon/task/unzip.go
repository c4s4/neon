package task

import (
	z "archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/c4s4/neon/neon/build"
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
      todir: 'build'

Notes:

- Existing files are overwritten. Directories and symbolic links are created.
- Extraction fails if an entry would be written outside the destination
  directory, including through a symbolic link.`,
	})
}

type unzipArgs struct {
	Unzip string `neon:"file"`
	Todir string `neon:"file"`
}

func unzip(context *build.Context, args interface{}) error {
	params := args.(unzipArgs)
	context.MessageArgs("Unzipping archive '%s' to directory '%s'...", params.Unzip, params.Todir)
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
	defer func() {
		_ = reader.Close()
	}()
	// all files are written through root, which forbids writing outside of
	// destination directory, even following symbolic links
	if err := os.MkdirAll(dir, DirFileMode); err != nil {
		return fmt.Errorf("creating destination directory '%s': %v", dir, err)
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return fmt.Errorf("opening destination directory '%s': %v", dir, err)
	}
	defer func() {
		_ = root.Close()
	}()
	for _, file := range reader.File {
		name, err := LocalName(file.Name)
		if err != nil {
			return err
		}
		if file.FileInfo().IsDir() {
			if err := root.MkdirAll(name, DirFileMode); err != nil {
				return fmt.Errorf("creating directory '%s': %v", name, err)
			}
			continue
		}
		if err := unzipEntry(root, file, name); err != nil {
			return err
		}
	}
	return nil
}

func unzipEntry(root *os.Root, file *z.File, name string) error {
	readCloser, err := file.Open()
	if err != nil {
		return fmt.Errorf("opening zip file %s: %v", file.Name, err)
	}
	defer func() {
		_ = readCloser.Close()
	}()
	if file.Mode()&os.ModeSymlink != 0 {
		target, err := io.ReadAll(readCloser)
		if err != nil {
			return fmt.Errorf("reading symbolic link '%s': %v", name, err)
		}
		if err := createLink(root, name, func() error { return root.Symlink(string(target), name) }); err != nil {
			return fmt.Errorf("creating symbolic link '%s': %v", name, err)
		}
		return nil
	}
	if err := root.MkdirAll(filepath.Dir(name), DirFileMode); err != nil {
		return fmt.Errorf("creating destination directory for '%s': %v", name, err)
	}
	dest, err := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode().Perm())
	if err != nil {
		return fmt.Errorf("creating destination file '%s': %v", name, err)
	}
	if _, err := io.Copy(dest, readCloser); err != nil {
		_ = dest.Close()
		return fmt.Errorf("copying to destination file '%s': %v", name, err)
	}
	if err := dest.Close(); err != nil {
		return fmt.Errorf("closing destination file '%s': %v", name, err)
	}
	return nil
}
