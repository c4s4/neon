package builtin

import (
	"testing"
)

func TestWindows2Unix(t *testing.T) {
	Assert(unixPath(`foo`), "foo", t)
	Assert(unixPath(`foo\bar`), "foo/bar", t)
	Assert(unixPath(`\foo\bar`), "/foo/bar", t)
	Assert(unixPath(`c:\foo\bar`), "/c/foo/bar", t)
}
