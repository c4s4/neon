package builtin

import (
	"neon/util"
	"testing"
)

func TestExists(t *testing.T) {
	if util.Windows() {
		if !Exists("/Windows") {
			t.Errorf("Error builtin exists")
		}
	} else {
		if !Exists("/tmp") {
			t.Errorf("Error builtin exists")
		}
	}
}
