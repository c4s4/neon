package builtin

import (
	"testing"
)

func TestDirectory(t *testing.T) {
	Assert(Directory("/foo/bar/spam.txt"), "/foo/bar", t)
}
