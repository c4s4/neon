package build

import (
	"neon/util"
	"os"
	"reflect"
	"testing"
)

func TestParseSingleton(t *testing.T) {
	object := map[string]interface{}{
		"singleton": 12345,
	}
	build := &Build{}
	ParseSingleton(object, build)
	if build.Singleton != "12345" {
		t.Errorf("Bad build doc: %s", build.infoDoc())
	}
}

func TestParseShell(t *testing.T) {
	// simple shell
	object := map[string]interface{}{
		"shell": []string{"sh", "-c"},
	}
	build := &Build{}
	ParseShell(object, build)
	expected := map[string][]string{
		"default": {"sh", "-c"},
	}
	if !reflect.DeepEqual(build.Shell, expected) {
		t.Errorf("Bad shell: %v", build.Shell)
	}
	// complex shell
	object = map[string]interface{}{
		"shell": map[string][]string{
			"foo": {"spam"},
			"bar": {"eggs"},
		},
	}
	build = &Build{}
	ParseShell(object, build)
	expected = map[string][]string{
		"foo": {"spam"},
		"bar": {"eggs"},
	}
	if !reflect.DeepEqual(build.Shell, expected) {
		t.Errorf("Bad shell: %v", build.Shell)
	}
	// default shell
	object = map[string]interface{}{}
	build = &Build{}
	ParseShell(object, build)
	expected = map[string][]string{
		"default": {"sh", "-c"},
		"windows": {"cmd", "/c"},
	}
	if !reflect.DeepEqual(build.Shell, expected) {
		t.Errorf("Bad shell: %v", build.Shell)
	}
}

func TestParseDefault(t *testing.T) {
	object := map[string]interface{}{
		"default": "test",
	}
	build := &Build{}
	ParseDefault(object, build)
	if !reflect.DeepEqual(build.Default, []string{"test"}) {
		t.Errorf("Bad build default: %v", build.Default)
	}
}

func TestParseDoc(t *testing.T) {
	object := map[string]interface{}{
		"doc": "test",
	}
	build := &Build{}
	ParseDoc(object, build)
	if build.Doc != "test" {
		t.Errorf("Bad build doc: %v", build.Doc)
	}
}

func TestParseRepository(t *testing.T) {
	// repo in build file
	object := map[string]interface{}{
		"repository": "repo",
	}
	build := &Build{}
	ParseRepository(object, build, "")
	if build.Repository != "repo" {
		t.Errorf("Bad build repo: %v", build.Repository)
	}
	// repo on command line
	object = map[string]interface{}{
		"repository": "repo",
	}
	build = &Build{}
	ParseRepository(object, build, "commandline")
	if build.Repository != "commandline" {
		t.Errorf("Bad build repo: %v", build.Repository)
	}
	// default repo
	object = map[string]interface{}{}
	build = &Build{}
	ParseRepository(object, build, "")
	if build.Repository != util.ExpandUserHome(DefaultRepository) {
		t.Errorf("Bad build repo: %v", build.Repository)
	}
}

func TestParseContext(t *testing.T) {
	object := map[string]interface{}{
		"context": []string{"foo", "bar"},
	}
	build := &Build{}
	ParseContext(object, build)
	expected := []string{"foo", "bar"}
	if !reflect.DeepEqual(build.Scripts, expected) {
		t.Errorf("Bad build context: %v", build.Scripts)
	}
}

func TestParseExtends(t *testing.T) {
	object := map[string]interface{}{
		"extends": []string{"foo", "bar"},
	}
	build := &Build{}
	ParseExtends(object, build)
	expected := []string{"foo", "bar"}
	if !reflect.DeepEqual(build.Extends, expected) {
		t.Errorf("Bad build extends: %v", build.Extends)
	}
}

func TestParseProperties(t *testing.T) {
	object := map[string]interface{}{
		"properties": map[string]interface{}{
			"foo": "spam",
			"bar": "eggs",
		},
	}
	build := &Build{}
	ParseProperties(object, build)
	expected := util.Object{
		"foo": "spam",
		"bar": "eggs",
	}
	if !reflect.DeepEqual(build.Properties, expected) {
		t.Errorf("Bad build properties: %v != %v", build.Properties, expected)
	}
}

func TestParseConfiguration(t *testing.T) {
	writeFile("/tmp", "config.yml", "foo: spam\nbar: eggs")
	defer os.RemoveAll("/tmp/config.yml")
	object := map[string]interface{}{
		"configuration": []string{"/tmp/config.yml"},
	}
	build := &Build{
		Properties: util.Object{},
	}
	ParseConfiguration(object, build)
	expected := util.Object{
		"foo": "spam",
		"bar": "eggs",
	}
	if !reflect.DeepEqual(build.Properties, expected) {
		t.Errorf("Bad build properties: %v", build.Properties)
	}
	if !reflect.DeepEqual(build.Config, []string{"/tmp/config.yml"}) {
		t.Errorf("Bad build config: %v", build.Config)
	}
}

func TestParseEnvironment(t *testing.T) {
	object := map[string]interface{}{
		"environment": map[string]string{
			"FOO": "SPAM",
			"BAR": "EGGS",
		},
	}
	build := &Build{}
	ParseEnvironment(object, build)
	expected := map[string]string{
		"FOO": "SPAM",
		"BAR": "EGGS",
	}
	if !reflect.DeepEqual(build.Environment, expected) {
		t.Errorf("Bad build environment: %v != %v", build.Environment, expected)
	}
}

func TestParseTargets(t *testing.T) {
	object := map[string]interface{}{
		"targets": util.Object{
			"test": util.Object{},
		},
	}
	build := &Build{}
	ParseTargets(object, build)
	if _, ok := build.Targets["test"]; !ok {
		t.Errorf("Bad build targets: missing test")
	}
}

func TestParseVersion(t *testing.T) {
	object := map[string]interface{}{
		"version": `greaterorequal("0.12")`,
	}
	build := &Build{}
	ParseVersion(object, build)
	if build.Version != `greaterorequal("0.12")` {
		t.Errorf("Bad build version: %v", build.Version)
	}
}
