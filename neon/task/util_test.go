package task

import (
	"testing"
)

func TestSanitizeName(t *testing.T) {
	if SanitizeName("/foo/bar") != "foo/bar" {
		t.Errorf("bad sanitization: %s", SanitizeName("/foo/bar"))
	}
	if SanitizeName("../foo/bar") != "foo/bar" {
		t.Errorf("bad sanitization: %s", SanitizeName("/foo/bar"))
	}
	if SanitizeName(`foo\bar`) != `foo/bar` {
		t.Errorf("bad sanitization: %s", SanitizeName(`foo\bar`))
	}
}

func TestRemoveStep(t *testing.T) {
	if RemoveStep("in step 1: message") != "message" {
		t.Errorf("bad step removing: '%s'", RemoveStep("in step 1: message"))
	}
}
