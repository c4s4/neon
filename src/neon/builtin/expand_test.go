package builtin

import (
	"neon/util"
	"testing"
	"os/user"
)

func TestExpand(t *testing.T) {
	user, _ := user.Current()
	home := user.HomeDir
	util.Assert(util.PathToUnix(Expand("~/foo")), util.PathToUnix(home + "/foo"), t)
}
