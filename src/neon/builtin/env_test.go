// +build !windows

package builtin

import (
	"testing"
	"strings"
)

func TestEnv(t *testing.T) {
	path := env("PATH")
	Assert(strings.Contains(path, "/bin"), true, t)
}
