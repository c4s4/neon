//go:build !windows
// +build !windows

package builtin

import (
	"strings"
	"testing"
)

func TestEnv(t *testing.T) {
	path := env("PATH")
	Assert(strings.Contains(path, "/bin"), true, t)
}
