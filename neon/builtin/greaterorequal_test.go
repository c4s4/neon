package builtin

import (
	"github.com/c4s4/neon/neon/build"
	"testing"
)

func TestGreaterOrEqual(t *testing.T) {
	build.NeonVersion = "1.2.3"
	if !greaterOrEqual("0") {
		t.Errorf("greaterorequal test failure")
	}
	if !greaterOrEqual("1.2.2") {
		t.Errorf("greaterorequal test failure")
	}
	if !greaterOrEqual("1.2.3") {
		t.Errorf("greaterorequal test failure")
	}
	if greaterOrEqual("2") {
		t.Errorf("greaterorequal test failure")
	}
}

func TestGreaterOrEqualPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	greaterOrEqual("x.y.z")
}
