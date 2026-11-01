package builtin

import (
	"testing"

	"github.com/c4s4/neon/neon/util"
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
	if exists("/path/that/doesnt/exist") {
		t.Errorf("Error builtin exists")
	}
}
