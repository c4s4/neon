package builtin

import (
	"neon/util"
	"testing"
)

func TestDirectory(t *testing.T) {
	util.Assert(Directory("/foo/bar/spam.txt"), "/foo/bar", t)
}
