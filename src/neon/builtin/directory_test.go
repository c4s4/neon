package builtin

import (
	"testing"
)

func TestDirectory(t *testing.T) {
	Assert(directory("/foo/bar/spam.txt"), "/foo/bar", t)
}
