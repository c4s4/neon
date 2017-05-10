package build

import (
	"testing"
)

func TestExpandNeonPath(t *testing.T) {
	expected := "~/.neon/go/latest/simple-command.yml"
	actual, err := ExpandNeonPath(`:go/simple-command.yml`)
	if err != nil || expected != actual {
		t.Error("Error expanding neon path")
	}
	expected = "~/.neon/go/1.2.3/simple-command.yml"
	actual, err = ExpandNeonPath(`:go/1.2.3/simple-command.yml`)
	if err != nil || expected != actual {
		t.Error("Error expanding neon path")
	}
	expected = "/foo/bar"
	actual, err = ExpandNeonPath("/foo/bar")
	if err != nil || expected != actual {
		t.Error("Error expanding neon path")
	}
	_, err = ExpandNeonPath(`:go`)
	if err == nil || err.Error() != `Bad Neon path ':go'` {
		t.Error("Error expanding neon path")
	}
}
