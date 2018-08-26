package builtin

import (
	"testing"
)

func TestEscapeUrl(t *testing.T) {
	if escapeURL("/foo bar") != "/foo%20bar" {
		t.Errorf("Error builtin escapeulr")
	}
}

func TestEscapeUrlPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	escapeURL("foo%ZZbar")
}
