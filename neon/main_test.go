package main

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/c4s4/neon/neon/build"
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
	if err := os.WriteFile(file, []byte(source), 0644); err != nil {
		t.Errorf("Error writing configuration file")
	}
	defer func() {
		_ = os.Remove(file)
	}()
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
	if err := os.WriteFile(file, []byte(source), 0644); err != nil {
		t.Errorf("Error writing configuration file")
	}
	defer func() {
		_ = os.Remove(file)
	}()
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
	if err := os.WriteFile(file, []byte(source), 0644); err != nil {
		t.Errorf("Error writing configuration file")
	}
	defer func() {
		_ = os.Remove(file)
	}()
	_, err := LoadConfiguration(file)
	if err != nil {
		t.Errorf("Error parsing configuration")
	}
	if build.Gray != true {
		t.Errorf("Error parsing configuration")
	}
}

func TestParseCommandLine(t *testing.T) {
	os.Args = []string{"cmd", "-file", "file", "-info", "-version", "-props", "{foo: bar, spam: eggs}",
		"-time", "-tasks", "-task", "task", "-targets", "-builtins", "-builtin", "builtin", "-tree",
		"-tasks-ref", "-builtins-ref", "-install", "install", "-repo", "repo", "-update", "-batch", "-grey",
		"-template", "template", "-templates", "-themes", "-theme", "test", "-parents", "target1", "target2"}
	opts := ParseCommandLine()
	Assert(opts.File, "file", t)
	Assert(opts.Info, true, t)
	Assert(opts.Version, true, t)
	Assert(opts.Props, "{foo: bar, spam: eggs}", t)
	Assert(opts.Time, true, t)
	Assert(opts.Tasks, true, t)
	Assert(opts.Task, "task", t)
	Assert(opts.PrintTargets, true, t)
	Assert(opts.Builtins, true, t)
	Assert(opts.Builtin, "builtin", t)
	Assert(opts.Tree, true, t)
	Assert(opts.TasksRef, true, t)
	Assert(opts.BuiltinsRef, true, t)
	Assert(opts.Install, "install", t)
	Assert(opts.Repo, "repo", t)
	Assert(opts.Update, true, t)
	Assert(opts.Batch, true, t)
	Assert(opts.Grey, true, t)
	Assert(opts.Template, "template", t)
	Assert(opts.Templates, true, t)
	Assert(opts.Parents, true, t)
	Assert(opts.Theme, "test", t)
	Assert(opts.Themes, true, t)
	Assert(opts.Targets, []string{"target1", "target2"}, t)
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
