package builtin

import (
	"testing"
)

func TestExists(t *testing.T) {
	if !Exists("/tmp") {
		t.Errorf("Error builtin exists")
	}
}
