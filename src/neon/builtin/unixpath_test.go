package builtin

import (
	"testing"
)

func TestWindows2Unix(t *testing.T) {
	Assert(UnixPath(`foo`), "foo", t)
	Assert(UnixPath(`foo\bar`), "foo/bar", t)
	Assert(UnixPath(`\foo\bar`), "/foo/bar", t)
	Assert(UnixPath(`c:\foo\bar`), "/c/foo/bar", t)
}
