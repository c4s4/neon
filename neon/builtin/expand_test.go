package builtin

import (
	"github.com/c4s4/neon/neon/util"
	"os/user"
	"testing"
)

func TestExpand(t *testing.T) {
	user, _ := user.Current()
	home := user.HomeDir
	Assert(util.PathToUnix(expand("~/foo")), util.PathToUnix(home+"/foo"), t)
}
