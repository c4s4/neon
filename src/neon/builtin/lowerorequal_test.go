package builtin

import (
	"neon/build"
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
