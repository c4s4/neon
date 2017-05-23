package builtin

import (
	"testing"
	"os/user"
)

func TestExpand(t *testing.T) {
	user, _ := user.Current()
	home := user.HomeDir
	if Expand("~/foo") != home + "/foo" {
		t.Errorf("Error builtin expand")
	}
}
