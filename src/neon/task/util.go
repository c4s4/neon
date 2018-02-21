package task

import (
	"runtime"
	"strings"
)

const (
	// FileMode is default file mode
	FileMode    = 0644
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
	filename = strings.Replace(filename, `\`, `/`, -1)
	filename = strings.TrimLeft(filename, "/.")
	return strings.Replace(filename, "../", "", -1)
}

// RemoveStep remove first part of an error message that include step:
// - message: error message (such as "in step 1: message")
// Return: error message without step (such as "message")
func RemoveStep(message string) string {
	position := strings.Index(message, ":")
	if position > -1 {
		return message[position+2:]
	}
	return message
}
