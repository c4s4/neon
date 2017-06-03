package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	DIR_FILE_MODE = 0755
)

// Read given file and return it as a byte slice
func ReadFile(file string) ([]byte, error) {
	path := ExpandUserHome(file)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file '%s': %v", file, err)
	}
	return bytes, nil
}

// Tells if file exists
func FileExists(file string) bool {
	file = ExpandUserHome(file)
	if stat, err := os.Stat(file); err == nil && !stat.IsDir() {
		return true
	} else {
		return false
	}
}

// Tells if directory exists
func DirExists(dir string) bool {
	dir = ExpandUserHome(dir)
	if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
		return true
	} else {
		return false
	}
}

// Copy source file to destination, preserving mode
func CopyFile(source, dest string) error {
	source = ExpandUserHome(source)
	dest = ExpandUserHome(dest)
	from, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("opening source file '%s': %v", source, err)
	}
	info, err := from.Stat()
	if err != nil {
		return fmt.Errorf("getting mode of source file '%s': %v", source, err)
	}
	defer from.Close()
	to, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("creating desctination file '%s': %v", dest, err)
	}
	defer to.Close()
	_, err = io.Copy(to, from)
	if err != nil {
		return fmt.Errorf("copying file: %v", err)
	}
	err = to.Sync()
	if err != nil {
		return fmt.Errorf("syncing destination file: %v", err)
	}
	if !Windows() {
		err = to.Chmod(info.Mode())
		if err != nil {
			return fmt.Errorf("changing mode of destination file '%s': %v", dest, err)
		}
	}
	return nil
}

// Copy files in root directory to destination directory. If flatten, all files
// are copied in destination directory, even if in source subdirectories
func CopyFilesToDir(dir string, files []string, toDir string, flatten bool) error {
	dir = ExpandUserHome(dir)
	toDir = ExpandUserHome(toDir)
	if stat, err := os.Stat(toDir); err != nil || !stat.IsDir() {
		return fmt.Errorf("destination directory doesn't exist")
	}
	for _, file := range files {
		file = ExpandUserHome(file)
		source := filepath.Join(dir, file)
		var dest string
		if flatten {
			base := filepath.Base(file)
			dest = filepath.Join(toDir, base)
		} else {
			dest = filepath.Join(toDir, file)
			destDir := filepath.Dir(dest)
			if !DirExists(destDir) {
				err := os.MkdirAll(destDir, DIR_FILE_MODE)
				if err != nil {
					return fmt.Errorf("creating directory for destination file: %v", err)
				}
			}
		}
		err := CopyFile(source, dest)
		if err != nil {
			return err
		}
	}
	return nil
}

// Move files in source directory to destination. If flatten is set to true, all
// files are moved in the root of destination directory.
func MoveFilesToDir(dir string, files []string, toDir string, flatten bool) error {
	dir = ExpandUserHome(dir)
	toDir = ExpandUserHome(toDir)
	if stat, err := os.Stat(toDir); err != nil || !stat.IsDir() {
		return fmt.Errorf("destination directory doesn't exist")
	}
	for _, file := range files {
		file = ExpandUserHome(file)
		source := filepath.Join(dir, file)
		var dest string
		if flatten {
			base := filepath.Base(file)
			dest = filepath.Join(toDir, base)
		} else {
			dest = filepath.Join(toDir, file)
			destDir := filepath.Dir(dest)
			if !DirExists(destDir) {
				err := os.MkdirAll(destDir, DIR_FILE_MODE)
				if err != nil {
					return fmt.Errorf("creating directory for destination file: %v", err)
				}
			}
		}
		err := os.Rename(source, dest)
		if err != nil {
			return err
		}
	}
	return nil
}

// Expand user home is path starts with "~/"
func ExpandUserHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		user, _ := user.Current()
		home := user.HomeDir
		path = filepath.Join(home, path[2:])
	}
	return path
}

// Expand user home in path and join it to root path if relative
func ExpandAndJoinToRoot(root, path string) string {
	path = ExpandUserHome(path)
	if filepath.IsAbs(path) {
		return path
	} else {
		return filepath.Join(root, path)
	}
}
