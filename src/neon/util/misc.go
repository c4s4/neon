package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	DIR_FILE_MODE = 0755
)

func ReadFile(file string) ([]byte, error) {
	path := ExpandUserHome(file)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file '%s': %v", file, err)
	}
	return bytes, nil
}

func FileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	} else {
		return false
	}
}

func DirExists(dir string) bool {
	if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
		return true
	} else {
		return false
	}
}

func CopyFile(source, dest string) error {
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
	err = to.Chmod(info.Mode())
	if err != nil {
		return fmt.Errorf("changing mode of destination file '%s': %v", dest, err)
	}
	return nil
}

func CopyFilesToDir(dir string, files []string, toDir string, flatten bool) error {
	if stat, err := os.Stat(toDir); err != nil || !stat.IsDir() {
		return fmt.Errorf("destination directory doesn't exist")
	}
	for _, file := range files {
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

func ExpandUserHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		user, _ := user.Current()
		home := user.HomeDir
		path = filepath.Join(home, path[2:])
	}
	return path
}

func FilePath(root, path string) string {
	path = ExpandUserHome(path)
	if filepath.IsAbs(path) {
		return path
	} else {
		return filepath.Join(root, path)
	}
}

func ToList(object interface{}) ([]interface{}, error) {
	slice := reflect.ValueOf(object)
	if slice.Kind() == reflect.Slice {
		result := make([]interface{}, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			result[i] = slice.Index(i).Interface()
		}
		return result, nil
	} else {
		return nil, fmt.Errorf("must be a list")
	}
}

func ValueToInterface(value reflect.Value) interface{} {
	if value.IsValid() {
		kind := value.Kind()
		if kind == reflect.Slice {
			result := make([]interface{}, value.Len())
			for i := 0; i < value.Len(); i++ {
				result[i] = ValueToInterface(value.Index(i))
			}
			return result
		} else if kind == reflect.Map {
			result := make(map[interface{}]interface{})
			for _, key := range value.MapKeys() {
				keyInterface := ValueToInterface(key)
				valueInterface := ValueToInterface(value.MapIndex(key))
				result[keyInterface] = valueInterface
			}
			return result
		} else {
			return value.Interface()
		}
	} else {
		return nil
	}
}

func MaxLength(lines []string) int {
	length := 0
	for _, line := range lines {
		if utf8.RuneCountInString(line) > length {
			length = utf8.RuneCountInString(line)
		}
	}
	return length
}

func Singleton(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	go func() {
		for {
			listener.Accept()
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return nil
}
