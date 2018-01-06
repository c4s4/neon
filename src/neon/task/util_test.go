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
