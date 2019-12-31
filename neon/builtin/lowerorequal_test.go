package builtin

import (
	"github.com/c4s4/neon/neon/build"
	"testing"
)

func TestLowerOrEqual(t *testing.T) {
	build.NeonVersion = "1.2.3"
	if lowerOrEqual("0") {
		t.Errorf("lowerorequal test failure")
	}
	if lowerOrEqual("1.2.2") {
		t.Errorf("lowerorequal test failure")
	}
	if !lowerOrEqual("1.2.3") {
		t.Errorf("lowerorequal test failure")
	}
	if !lowerOrEqual("1.2.4") {
		t.Errorf("lowerorequal test failure")
	}
	if !lowerOrEqual("2") {
		t.Errorf("lowerorequal test failure")
	}
}

func TestLowerOrEqualPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	lowerOrEqual("x.y.z")
}
