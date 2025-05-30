package util

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/mattn/go-zglob"
)

const (
	// FileMode is default file mode
	FileMode = 0644
	// DirFileMode is default directory file mode
	DirFileMode = 0755
)

// ReadFile reads given file and return it as a byte slice:
// - file: the file to read
// Return:
// - content as a slice of bytes
// - an error if something went wrong
func ReadFile(file string) ([]byte, error) {
	path := ExpandUserHome(file)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file '%s': %v", file, err)
	}
	return bytes, nil
}

// FileExists tells if given file exists:
// - file: the name of the file to test
// Return: a boolean that tells if file exists
func FileExists(file string) bool {
	file = ExpandUserHome(file)
	if stat, err := os.Stat(file); err == nil && !stat.IsDir() {
		return true
	}
	return false
}

// DirExists tells if directory exists:
// - dir: directory to test
// Return: a boolean that tells if directory exists
func DirExists(dir string) bool {
	dir = ExpandUserHome(dir)
	if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
		return true
	}
	return false
}

// FileIsExecutable tells if given file is executable by user:
// - file: file to test
// Return: a boolean that tells if file is executable by user
func FileIsExecutable(file string) bool {
	if stat, err := os.Stat(file); err == nil && stat.Mode()&0111 != 0 {
		return true
	}
	return false
}

// CopyFile copies source file to destination, preserving mode:
// - source: the source file
// - dest: the destination file
// Return: error if something went wrong
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
	defer func() {
		_ = from.Close()
	}()
	to, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("creating destination file '%s': %v", dest, err)
	}
	defer func() {
		_ = to.Close()
	}()
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

// CopyFilesToDir copies files in root directory to destination directory:
// - dir: root directory
// - files: globs of source files
// - toDir: destination directory
// - flatten: tells if files should be flatten in destination directory
// Return: an error if something went wrong
func CopyFilesToDir(dir string, files []string, toDir string, flatten bool) error {
	if stat, err := os.Stat(toDir); err != nil || !stat.IsDir() {
		return fmt.Errorf("destination directory doesn't exist")
	}
	for _, file := range files {
		source := file
		if !filepath.IsAbs(file) {
			source = filepath.Join(dir, file)
		}
		var dest string
		if flatten || filepath.IsAbs(file) {
			base := filepath.Base(file)
			dest = filepath.Join(toDir, base)
		} else {
			dest = filepath.Join(toDir, file)
			destDir := filepath.Dir(dest)
			if !DirExists(destDir) {
				err := os.MkdirAll(destDir, DirFileMode)
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

// MoveFilesToDir moves files in source directory to destination:
// - dir: root directory of source files
// - files: globs of files to move
// - toDir: destination directory
// - flatten: tells if files should be flatten in destination directory
// Return: an error if something went wrong
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
				err := os.MkdirAll(destDir, DirFileMode)
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

// ExpandUserHome expand path starting with "~/":
// - path: the path to expand
// Return: expanded path
func ExpandUserHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		user, _ := user.Current()
		home := user.HomeDir
		path = filepath.Join(home, path[2:])
	}
	return path
}

// PathToUnix turns a path to Unix format (with "/"):
// - path: path to turn to unix format
// Return: converted path
func PathToUnix(path string) string {
	// replace path separator \ with /
	path = strings.ReplaceAll(path, "\\", "/")
	// replace c: with /c
	r := regexp.MustCompile("^[A-Za-z]:.*$")
	if r.MatchString(path) {
		path = "/" + path[0:1] + path[2:]
	}
	return path
}

// PathToWindows turns a path to Windows format (with "\"):
// - path: path to turn to windows format
// Return: converted path
func PathToWindows(path string) string {
	// replace path separator / with \
	path = strings.ReplaceAll(path, "/", "\\")
	// replace /c/ with c:/
	r := regexp.MustCompile(`^\\[A-Za-z]\\.*$`)
	if r.MatchString(path) {
		path = path[1:2] + ":" + path[2:]
	}
	return path
}

// FindFiles finds files in the context:
// - dir: the search root directory (current dir if empty)
// - includes: the list of globs to include
// - excludes: the list of globs to exclude
// - folder: tells if we should include folders
// Return the list of files as a slice of strings
func FindFiles(dir string, includes, excludes []string, folder bool) ([]string, error) {
	var err error
	included := joinPath(dir, includes)
	excluded := joinPath(dir, excludes)
	included, err = filterFolders(included, folder)
	if err != nil {
		return nil, err
	}
	files := filterExcluded(included, excluded)
	files, err = makeRelative(dir, files)
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

// FindInPath search given executable in PATH:
// - executable: executable to search.
// Return: list of directories containing executable
func FindInPath(executable string) []string {
	path := os.Getenv("PATH")
	dirs := strings.Split(path, string(os.PathListSeparator))
	var paths []string
	for _, dir := range dirs {
		file := filepath.Join(dir, executable)
		if FileIsExecutable(file) {
			paths = append(paths, file)
		}
	}
	return paths
}

func joinPath(dir string, paths []string) []string {
	var joined []string
	for _, path := range paths {
		if !filepath.IsAbs(path) {
			path = filepath.Join(dir, path)
		}
		joined = append(joined, path)
	}
	return joined
}

func filterFolders(included []string, folder bool) ([]string, error) {
	var candidates []string
	for _, include := range included {
		list, _ := zglob.Glob(include)
		for _, file := range list {
			stat, err := os.Stat(file)
			if err != nil {
				return nil, fmt.Errorf("stating file: %v", err)
			}
			if stat.Mode().IsRegular() ||
				(stat.Mode().IsDir() && folder) {
				candidates = append(candidates, file)
			}
		}
	}
	return candidates, nil
}

func filterExcluded(candidates []string, excluded []string) []string {
	var files []string
	if excluded != nil {
		for index, file := range candidates {
			for _, exclude := range excluded {
				match, err := zglob.Match(exclude, file)
				if match || err != nil {
					candidates[index] = ""
				}
			}
		}
		for _, file := range candidates {
			if file != "" {
				files = append(files, file)
			}
		}
	} else {
		files = candidates
	}
	return files
}

func makeRelative(dir string, files []string) ([]string, error) {
	var err error
	for index, file := range files {
		if dir != "" {
			files[index], err = filepath.Rel(dir, file)
			if err != nil {
				return nil, err
			}
		}
	}
	return files, nil
}
