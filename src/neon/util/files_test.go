package util

import (
	"testing"
)

func TestWindowsToUnix(t *testing.T) {
	Assert(PathToUnix("foo"), "foo", t)
	Assert(PathToUnix("foo\\bar"), "foo/bar", t)
	Assert(PathToUnix("\\foo\\bar"), "/foo/bar", t)
	Assert(PathToUnix("C:\\foo\\bar"), "/C/foo/bar", t)
	Assert(PathToUnix("c:\\foo\\bar"), "/c/foo/bar", t)
}

func TestUnixToWindows(t *testing.T) {
	Assert(PathToWindows("foo"), "foo", t)
	Assert(PathToWindows("foo/bar"), "foo\\bar", t)
	Assert(PathToWindows("/foo/bar"), "\\foo\\bar", t)
	Assert(PathToWindows("/C/foo/bar"), "C:\\foo\\bar", t)
	Assert(PathToWindows("/c/foo/bar"), "c:\\foo\\bar", t)
}
