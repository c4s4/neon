package builtin

import (
	"testing"
)

func TestUnix2Windows(t *testing.T) {
	Assert(WindowsPath("foo"), `foo`, t)
	Assert(WindowsPath("foo/bar"), `foo\bar`, t)
	Assert(WindowsPath("/foo/bar"), `\foo\bar`, t)
	Assert(WindowsPath("/c/foo/bar"), `c:\foo\bar`, t)
}
