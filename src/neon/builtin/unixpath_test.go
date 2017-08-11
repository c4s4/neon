package builtin

import (
	"testing"
	"neon/util"
)

func TestWindows2Unix(t *testing.T) {
	util.Assert(UnixPath(`foo`), "foo", t)
	util.Assert(UnixPath(`foo\bar`), "foo/bar", t)
	util.Assert(UnixPath(`\foo\bar`), "/foo/bar", t)
	util.Assert(UnixPath(`c:\foo\bar`), "/c/foo/bar", t)
}
