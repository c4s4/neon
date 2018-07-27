package builtin

import (
	"neon/build"
	"testing"
)

func TestLower(t *testing.T) {
	build.NeonVersion = "1.2.3"
	if lower("0") {
		t.Errorf("lower test failure")
	}
	if lower("1.2.3") {
		t.Errorf("lower test failure")
	}
	if !lower("1.2.4") {
		t.Errorf("lower test failure")
	}
	if !lower("2") {
		t.Errorf("lower test failure")
	}
}

func TestLowerPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	lower("x.y.z")
}
