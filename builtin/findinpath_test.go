// +build !windows

package builtin

import (
	"testing"
)

func TestFindInPath(t *testing.T) {
	Assert(findInPath("ls"), []string{"/bin/ls"}, t)
}
