package builtin

import (
	"testing"
)

func TestThrow(t *testing.T) {
	err := throw("test")
	if err == nil {
		t.Errorf("No error returned")
	}
}
