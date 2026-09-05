package task

import (
	"reflect"
	"testing"
)

func TestCafeCommand(t *testing.T) {
	// test darwin
	wrapper, err := cafeCommand("darwin")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else {
		expected := []string{"caffeinate", "-i"}
		if !reflect.DeepEqual(wrapper, expected) {
			t.Errorf("actual (%v) != expected (%v)", wrapper, expected)
		}
	}
	// test linux
	wrapper, err = cafeCommand("linux")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else {
		expected := []string{"systemd-inhibit", "--what=idle", "--who=NeON",
			"--why=build", "--mode=block"}
		if !reflect.DeepEqual(wrapper, expected) {
			t.Errorf("actual (%v) != expected (%v)", wrapper, expected)
		}
	}
	// test unsupported operating systems
	if _, err := cafeCommand("windows"); err == nil {
		t.Error("expected an error for unsupported operating system")
	}
	if _, err := cafeCommand("freebsd"); err == nil {
		t.Error("expected an error for unsupported operating system")
	}
}
