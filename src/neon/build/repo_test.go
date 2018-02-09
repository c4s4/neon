package build

import (
	"neon/util"
	"testing"
)

func TestPluginPath(t *testing.T) {
	build := &Build{
		Repository: "~/.neon",
	}
	path, err := build.ParentPath("foo/bar/spam.yml")
	Assert(err, nil, t)
	Assert(path, util.ExpandUserHome("~/.neon/foo/bar/spam.yml"), t)
}
