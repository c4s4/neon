package builtin

import (
	"neon/util"
	"testing"
)

func TestExists(t *testing.T) {
	if util.Windows() {
		if !exists("/Windows") {
			t.Errorf("Error builtin exists")
		}
	} else {
		if !exists("/tmp") {
			t.Errorf("Error builtin exists")
		}
	}
}
