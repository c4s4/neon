package builtin

import (
	"os/user"
	"testing"

	"github.com/c4s4/neon/neon/util"
)

func TestExpand(t *testing.T) {
	user, _ := user.Current()
	home := user.HomeDir
	Assert(util.PathToUnix(expand("~/foo")), util.PathToUnix(home+"/foo"), t)
}
