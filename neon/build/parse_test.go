package build

import (
	"github.com/c4s4/neon/neon/util"
	"os"
	"reflect"
	"testing"
)

func TestParseSingleton(t *testing.T) {
	object := map[string]interface{}{
		"singleton": 12345,
	}
	build := &Build{}
	if err := ParseSingleton(object, build); err != nil {
		t.Fatalf("parsing singleton: %v", err)
	}
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
	if err := ParseShell(object, build); err != nil {
		t.Fatalf("parse shell: %v", err)
	}
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
	if err := ParseShell(object, build); err != nil {
		t.Fatalf("parse shell: %v", err)
	}
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
	if err := ParseShell(object, build); err != nil {
		t.Fatalf("parse shell: %v", err)
	}
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
	if err := ParseDefault(object, build); err != nil {
		t.Fatalf("parse default: %v", err)
	}
	if !reflect.DeepEqual(build.Default, []string{"test"}) {
		t.Errorf("Bad build default: %v", build.Default)
	}
}

func TestParseDoc(t *testing.T) {
	object := map[string]interface{}{
		"doc": "test",
	}
	build := &Build{}
	if err := ParseDoc(object, build); err != nil {
		t.Fatalf("parse doc: %v", err)
	}
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
	if err := ParseRepository(object, build, ""); err != nil {
		t.Fatalf("parse repository: %v", err)
	}
	if build.Repository != "repo" {
		t.Errorf("Bad build repo: %v", build.Repository)
	}
	// repo on command line
	object = map[string]interface{}{
		"repository": "repo",
	}
	build = &Build{}
	if err := ParseRepository(object, build, "commandline"); err != nil {
		t.Fatalf("parse repository: %v", err)
	}
	if build.Repository != "commandline" {
		t.Errorf("Bad build repo: %v", build.Repository)
	}
	// default repo
	object = map[string]interface{}{}
	build = &Build{}
	if err := ParseRepository(object, build, ""); err != nil {
		t.Fatalf("parse repository: %v", err)
	}
	if build.Repository != util.ExpandUserHome(DefaultRepository) {
		t.Errorf("Bad build repo: %v", build.Repository)
	}
}

func TestParseContext(t *testing.T) {
	object := map[string]interface{}{
		"context": []string{"foo", "bar"},
	}
	build := &Build{}
	if err := ParseContext(object, build); err != nil {
		t.Fatalf("parse context: %v", err)
	}
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
	if err := ParseExtends(object, build); err != nil {
		t.Fatalf("parse extends: %v", err)
	}
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
	if err := ParseProperties(object, build); err != nil {
		t.Fatalf("parse properties: %v", err)
	}
	expected := util.Object{
		"foo": "spam",
		"bar": "eggs",
	}
	if !reflect.DeepEqual(build.Properties, expected) {
		t.Errorf("Bad build properties: %v != %v", build.Properties, expected)
	}
}

func TestParseConfiguration(t *testing.T) {
	if _, err := WriteFile("/tmp", "config.yml", "foo: spam\nbar: eggs"); err != nil {
		t.Fatalf("write file: %v", err)
	}
	defer os.RemoveAll("/tmp/config.yml")
	object := map[string]interface{}{
		"configuration": []string{"/tmp/config.yml"},
	}
	build := &Build{
		Properties: util.Object{},
	}
	if err := ParseConfiguration(object, build); err != nil {
		t.Fatalf("parse configuration: %v", err)
	}
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
	if err := ParseEnvironment(object, build); err != nil {
		t.Fatalf("parse environment: %v", err)
	}
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
	if err := ParseTargets(object, build); err != nil {
		t.Fatalf("parse targets: %v", err)
	}
	if _, ok := build.Targets["test"]; !ok {
		t.Errorf("Bad build targets: missing test")
	}
}

func TestParseVersion(t *testing.T) {
	object := map[string]interface{}{
		"version": `greaterorequal("0.12")`,
	}
	build := &Build{}
	if err := ParseVersion(object, build); err != nil {
		t.Fatalf("parse version: %v", err)
	}
	if build.Version != `greaterorequal("0.12")` {
		t.Errorf("Bad build version: %v", build.Version)
	}
}
