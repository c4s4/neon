package build

import (
	"testing"
	"reflect"
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
