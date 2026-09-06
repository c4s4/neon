package task

import (
	"testing"
)

func TestSanitizeName(t *testing.T) {
	if SanitizeName("/foo/bar") != "foo/bar" {
		t.Errorf("bad sanitization: %s", SanitizeName("/foo/bar"))
	}
	if SanitizeName("../foo/bar") != "foo/bar" {
		t.Errorf("bad sanitization: %s", SanitizeName("../foo/bar"))
	}
	if SanitizeName(`foo\bar`) != `foo/bar` {
		t.Errorf("bad sanitization: %s", SanitizeName(`foo\bar`))
	}
	if SanitizeName("./.hidden") != ".hidden" {
		t.Errorf("bad sanitization: %s", SanitizeName("./.hidden"))
	}
	if SanitizeName(".hidden") != ".hidden" {
		t.Errorf("bad sanitization: %s", SanitizeName(".hidden"))
	}
	if SanitizeName("./foo/.bar/baz") != "foo/.bar/baz" {
		t.Errorf("bad sanitization: %s", SanitizeName("./foo/.bar/baz"))
	}
}

func TestRemoveStep(t *testing.T) {
	if RemoveStep("in step 1: message") != "message" {
		t.Errorf("bad step removing: '%s'", RemoveStep("in step 1: message"))
	}
}
