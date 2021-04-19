package main

import (
	"github.com/c4s4/neon/neon/build"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestParseConfiguration(t *testing.T) {
	file := "/tmp/neon.yml"
	source := `grey: true
theme: bee
time: true
repo: ~/.neon
colors:
  title: [FgBlue]
  ok:    [FgHiGreen, Bold]
  error: [FgHiRed, Bold]
links:
foo: bars`
	if err := ioutil.WriteFile(file, []byte(source), 0644); err != nil {
		t.Errorf("Error writing configuration file")
	}
	defer os.Remove(file)
	configuration, err := ParseConfiguration(file)
	if err != nil {
		t.Errorf("Error parsing configuration")
	}
	if configuration.Grey != true {
		t.Errorf("Error parsing configuration")
	}
	if configuration.Theme != "bee" {
		t.Errorf("Error parsing configuration")
	}
	if configuration.Time != true {
		t.Errorf("Error parsing configuration")
	}
	if configuration.Repo != "~/.neon" {
		t.Errorf("Error parsing configuration")
	}
}

func TestParseConfigurationError(t *testing.T) {
	file := "/tmp/neon.yml"
	source := "- foo\nbar:"
	if err := ioutil.WriteFile(file, []byte(source), 0644); err != nil {
		t.Errorf("Error writing configuration file")
	}
	defer os.Remove(file)
	_, err := ParseConfiguration(file)
	if err == nil || err.Error() != "yaml: line 1: did not find expected '-' indicator" {
		t.Errorf("Error parsing configuration: %v", err)
	}
}

func TestLoadConfiguration(t *testing.T) {
	file := "/tmp/neon.yml"
	source := `grey: true
theme: bee
time: true
repo: ~/.neon
colors:
  title: [FgBlue]
  ok:    [FgHiGreen, Bold]
  error: [FgHiRed, Bold]
links:
foo: bars`
	if err := ioutil.WriteFile(file, []byte(source), 0644); err != nil {
		t.Errorf("Error writing configuration file")
	}
	defer os.Remove(file)
	_, err := LoadConfiguration(file)
	if err != nil {
		t.Errorf("Error parsing configuration")
	}
	if build.Grey != true {
		t.Errorf("Error parsing configuration")
	}
}

func TestParseCommandLine(t *testing.T) {
	os.Args = []string{"cmd", "-file", "file", "-info", "-version", "-props", "{foo: bar, spam: eggs}",
		"-time", "-tasks", "-task", "task", "-targets", "-builtins", "-builtin", "builtin", "-tree",
		"-tasks-ref", "-builtins-ref", "-install", "install", "-repo", "repo", "-update", "-grey",
		"-template", "template", "-templates", "-themes", "-theme", "test", "-parents", "target1", "target2"}
	file, info, version, props, timeit, tasks, task, targs, builtins, builtin, tree, tasksRef, builtinsRef,
		install, repo, update, grey, template, templates, parents, theme, themes, targets := ParseCommandLine()
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
	Assert(tree, true, t)
	Assert(tasksRef, true, t)
	Assert(builtinsRef, true, t)
	Assert(install, "install", t)
	Assert(repo, "repo", t)
	Assert(update, true, t)
	Assert(grey, true, t)
	Assert(template, "template", t)
	Assert(templates, true, t)
	Assert(parents, true, t)
	Assert(theme, "test", t)
	Assert(themes, true, t)
	Assert(targets, []string{"target1", "target2"}, t)
}

func TestFindBuildFile(t *testing.T) {
	var configuration = &Configuration{}
	if os.Getenv("TRAVIS") == "true" {
		t.Skip("skip test on travis")
	}
	file, base, err := FindBuildFile("build.yml", "", configuration)
	if err != nil {
		t.Errorf("error finding build file: %v", err)
	}
	if !strings.HasSuffix(file, "neon/build.yml") {
		t.Errorf("expected 'test' but got '%s' instead", file)
	}
	if !strings.HasSuffix(base, "neon") {
		t.Errorf("expected 'test' but got '%s' instead", file)
	}
	_, _, err = FindBuildFile("toto.xyz", "", configuration)
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
