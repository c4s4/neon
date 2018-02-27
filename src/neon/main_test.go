package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestParseCommandLine(t *testing.T) {
	os.Args = []string{"cmd", "-file", "file", "-info", "-version", "-props", "{foo: bar, spam: eggs}",
		"-time", "-tasks", "-task", "task", "-targets", "-builtins", "-builtin", "builtin",
		"-refs", "-install", "install", "-repo", "repo", "-grey", "-template", "template",
		"-templates", "-parents", "target1", "target2"}
	file, info, version, props, timeit, tasks, task, targs, builtins, builtin, refs, install,
		repo, grey, template, templates, parents, targets := ParseCommandLine()
	Assert(file, "file", t)
	Assert(info, true, t)
	Assert(version, true, t)
	Assert(props, "{foo: bar, spam: eggs}", t)
	Assert(timeit, true, t)
	Assert(tasks, true, t)
	Assert(task, "task", t)
	Assert(targs, true, t)
	Assert(builtins, true, t)
	Assert(builtin, "builtin", t)
	Assert(refs, true, t)
	Assert(install, "install", t)
	Assert(repo, "repo", t)
	Assert(grey, true, t)
	Assert(template, "template", t)
	Assert(templates, true, t)
	Assert(parents, true, t)
	Assert(targets, []string{"target1", "target2"}, t)
}

func TestFindBuildFile(t *testing.T) {
	file, err := FindBuildFile("build.yml")
	if err != nil {
		t.Errorf("error finding build file: %v", err)
	}
	if !strings.HasSuffix(file, "neon/build.yml") {
		t.Errorf("expected 'test' but got '%s' instead", file)
	}
	_, err = FindBuildFile("toto.xyz")
	if err == nil {
		t.Errorf("error finding build file")
	}
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
