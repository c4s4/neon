package task

import (
	"path/filepath"
	"runtime"
	"strings"
)

const (
	FILE_MODE     = 0644
	DIR_FILE_MODE = 0755
)

func SanitizeName(filename string) string {
	if len(filename) > 1 && filename[1] == ':' &&
		runtime.GOOS == "windows" {
		filename = filename[2:]
	}
	filename = filepath.ToSlash(filename)
	filename = strings.TrimLeft(filename, "/.")
	return strings.Replace(filename, "../", "", -1)
}
