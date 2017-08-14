package builtin

import (
	"neon/util"
	"testing"
	"strings"
)

func TestAbsolute(t *testing.T) {
	actual := util.PathToUnix(Absolute("foo/../bar/spam.txt"))
	if !strings.HasSuffix(actual, "bar/spam.txt") {
		t.Errorf("TestAbsolute failed")
	}
}
