package builtin

import (
	"github.com/c4s4/neon/build"
	"testing"
)

func TestGreater(t *testing.T) {
	build.NeonVersion = "1.2.3"
	if !greater("0") {
		t.Errorf("greater test failure")
	}
	if !greater("1.2.2") {
		t.Errorf("greater test failure")
	}
	if greater("1.2.3") {
		t.Errorf("greater test failure")
	}
	if greater("2") {
		t.Errorf("greater test failure")
	}
}

func TestGreaterPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	greater("x.y.z")
}
