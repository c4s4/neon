package builtin

import (
	"strings"
	"testing"

	"github.com/c4s4/neon/neon/util"
)

func TestAbsolute(t *testing.T) {
	actual := util.PathToUnix(absolute("foo/../bar/spam.txt"))
	if !strings.HasSuffix(actual, "bar/spam.txt") {
		t.Errorf("TestAbsolute failed")
	}
}
