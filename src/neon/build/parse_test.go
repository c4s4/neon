package build

import (
	"neon/util"
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
