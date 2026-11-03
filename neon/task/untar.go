package task

import (
	t "archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/c4s4/neon/neon/build"
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
  the tar archive is uncompressed with gzip.
- Existing files are overwritten. Directories, symbolic and hard links are
  created.
- Extraction fails if an entry would be written outside the destination
  directory, including through a symbolic link.`,
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
	// all files are written through root, which forbids writing outside of
	// destination directory, even following symbolic links
	if err := os.MkdirAll(dir, DirFileMode); err != nil {
		return fmt.Errorf("creating destination directory %s: %v", dir, err)
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return fmt.Errorf("opening destination directory %s: %v", dir, err)
	}
	defer func() {
		_ = root.Close()
	}()
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
		name, err := LocalName(header.Name)
		if err != nil {
			return err
		}
		if err := untarEntry(root, tarReader, header, name); err != nil {
			return err
		}
	}
}

func untarEntry(root *os.Root, reader io.Reader, header *t.Header, name string) error {
	switch header.Typeflag {
	case t.TypeDir:
		if err := root.MkdirAll(name, DirFileMode); err != nil {
			return fmt.Errorf("creating directory %s: %v", name, err)
		}
	case t.TypeReg:
		if err := root.MkdirAll(filepath.Dir(name), DirFileMode); err != nil {
			return fmt.Errorf("creating destination directory for %s: %v", name, err)
		}
		dest, err := root.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode).Perm())
		if err != nil {
			return fmt.Errorf("creating destination file %s: %v", name, err)
		}
		if _, err := io.Copy(dest, reader); err != nil {
			_ = dest.Close()
			return fmt.Errorf("copying to destination file %s: %v", name, err)
		}
		if err := dest.Close(); err != nil {
			return fmt.Errorf("closing destination file %s: %v", name, err)
		}
		if err := root.Chtimes(name, header.ModTime, header.ModTime); err != nil {
			return fmt.Errorf("setting modification time of %s: %v", name, err)
		}
	case t.TypeSymlink:
		if err := createLink(root, name, func() error { return root.Symlink(header.Linkname, name) }); err != nil {
			return fmt.Errorf("creating symbolic link %s: %v", name, err)
		}
	case t.TypeLink:
		target, err := LocalName(header.Linkname)
		if err != nil {
			return err
		}
		if err := createLink(root, name, func() error { return root.Link(target, name) }); err != nil {
			return fmt.Errorf("creating hard link %s: %v", name, err)
		}
	}
	return nil
}

// createLink creates a link, removing any existing file with the same name
func createLink(root *os.Root, name string, create func() error) error {
	if err := root.MkdirAll(filepath.Dir(name), DirFileMode); err != nil {
		return err
	}
	if err := root.Remove(name); err != nil && !os.IsNotExist(err) {
		return err
	}
	return create()
}
