package builtin

import (
	"neon/util"
	"testing"
)

func TestOspath(t *testing.T) {
	if util.Windows() {
		Assert(osPath(`foo/bar`), `foo\bar`, t)
		Assert(osPath(`foo\bar`), `foo\bar`, t)
	} else {
		Assert(osPath(`foo/bar`), `foo/bar`, t)
		Assert(osPath(`foo\bar`), `foo/bar`, t)
	}
}
