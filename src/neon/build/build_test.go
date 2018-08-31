package build

import (
	"runtime"
	"testing"
)

func TestGetShell(t *testing.T) {
	build := &Build{
		Shell: map[string][]string{
			runtime.GOOS: {"foo"},
			"other":      {"bar"},
		},
	}
	shell, err := build.GetShell()
	if err != nil {
		t.Fail()
	}
	Assert(shell, []string{"foo"}, t)
}
