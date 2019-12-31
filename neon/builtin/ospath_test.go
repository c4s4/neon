package builtin

import (
	"github.com/c4s4/neon/neon/util"
	"testing"
)

func TestOspathUnix(t *testing.T) {
	goos := util.GOOS
	util.GOOS = "linux"
	Assert(osPath(`foo/bar`), `foo/bar`, t)
	Assert(osPath(`foo\bar`), `foo/bar`, t)
	util.GOOS = goos
}

func TestOspathWindows(t *testing.T) {
	goos := util.GOOS
	util.GOOS = "windows"
	Assert(osPath(`foo/bar`), `foo\bar`, t)
	Assert(osPath(`foo\bar`), `foo\bar`, t)
	util.GOOS = goos
}
