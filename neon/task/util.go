package task

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// FileMode is default file mode
	FileMode = 0644
	// DirFileMode is default file mode for directories
	DirFileMode = 0755
)

// SanitizeName cleans a path for archive:
// - filename: the path to sanitize.
// Return: sanitized path
func SanitizeName(filename string) string {
	if len(filename) > 1 && filename[1] == ':' &&
		runtime.GOOS == "windows" {
		filename = filename[2:]
	}
	filename = strings.ReplaceAll(filename, `\`, `/`)
	for strings.HasPrefix(filename, "./") {
		filename = filename[2:]
	}
	filename = strings.TrimLeft(filename, "/")
	return strings.ReplaceAll(filename, "../", "")
}

// SafeJoin joins an archive entry name to the destination directory, making
// sure the resulting path stays inside this directory (preventing Zip Slip):
// - dir: the destination directory.
// - name: the name of the archive entry.
// Return: the joined path and an error if the entry would be written outside
// of the destination directory.
func SafeJoin(dir, name string) (string, error) {
	local, err := LocalName(name)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, local), nil
}

// LocalName checks that an archive entry name is local (not absolute and
// without ".." that would go out of the destination directory):
// - name: the name of the archive entry (with "/" separators).
// Return: the name with OS separators and an error if it is not local.
func LocalName(name string) (string, error) {
	local := filepath.FromSlash(name)
	if !filepath.IsLocal(local) {
		return "", fmt.Errorf("illegal path in archive: '%s'", name)
	}
	return filepath.Clean(local), nil
}

// RemoveStep removes first part of an error message that includes step:
// - message: error message (such as "in step 1: message")
// Return: error message without step (such as "message")
func RemoveStep(message string) string {
	position := strings.Index(message, ":")
	if position > -1 {
		return message[position+2:]
	}
	return message
}
