package builtin

import (
	"github.com/c4s4/neon/util"
	"strings"
	"testing"
)

func TestAbsolute(t *testing.T) {
	actual := util.PathToUnix(absolute("foo/../bar/spam.txt"))
	if !strings.HasSuffix(actual, "bar/spam.txt") {
		t.Errorf("TestAbsolute failed")
	}
}
