package builtin

import (
	"testing"
	"neon/util"
)

func TestUnix2Windows(t *testing.T) {
	util.Assert(WindowsPath("foo"), `foo`, t)
	util.Assert(WindowsPath("foo/bar"), `foo\bar`, t)
	util.Assert(WindowsPath("/foo/bar"), `\foo\bar`, t)
	util.Assert(WindowsPath("/c/foo/bar"), `c:\foo\bar`, t)
}
