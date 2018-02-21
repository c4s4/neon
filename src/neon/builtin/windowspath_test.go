package builtin

import (
	"testing"
)

func TestUnix2Windows(t *testing.T) {
	Assert(windowsPath("foo"), `foo`, t)
	Assert(windowsPath("foo/bar"), `foo\bar`, t)
	Assert(windowsPath("/foo/bar"), `\foo\bar`, t)
	Assert(windowsPath("/c/foo/bar"), `c:\foo\bar`, t)
}
