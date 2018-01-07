package build

import (
	"neon/util"
	"reflect"
	"runtime"
	"testing"
)

func TestPluginPath(t *testing.T) {
	build := &Build{
		Repository: "~/.neon",
	}
	Assert(build.PluginPath("foo/bar"),
		util.ExpandUserHome("~/.neon/foo/bar"), t)
}

func TestPluginName(t *testing.T) {
	build := &Build{}
	Assert(build.PluginName("foo/bar/spam.yml"), "foo/bar", t)
}

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

// Assert make an assertion for testing purpose, failing test if different:
// - actual: actual value
// - expected: expected value
// - t: test
func Assert(actual, expected interface{}, t *testing.T) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual (\"%s\") != expected (\"%s\")", actual, expected)
	}
}
