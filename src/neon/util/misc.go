package util

import (
	"os"
)

func FileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	} else {
		return false
	}
}
