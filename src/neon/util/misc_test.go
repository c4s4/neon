package util

import (
	"testing"
)

func TestRemoveBlankLines(t *testing.T) {
	expected := "foo\nbar"
	actual := RemoveBlankLines("foo\n\nbar")
	if expected != actual {
		t.Errorf("Error removing blank lines: '%s' != '%s'", expected, actual)
	}
	expected = "foo\nbar"
	actual = RemoveBlankLines("foo\n  \nbar")
	if expected != actual {
		t.Errorf("Error removing blank lines: '%s' != '%s'", expected, actual)
	}
	expected = "foo\nbar"
	actual = RemoveBlankLines("foo\n  \n\t\nbar")
	if expected != actual {
		t.Errorf("Error removing blank lines: '%s' != '%s'", expected, actual)
	}
}

func TestToString(t *testing.T) {
	actual, err := ToString("string")
	if err != nil || actual != "string" {
		t.Errorf("ToString test failure")
	}
	_, err = ToString(25)
	if err == nil || err.Error() != "25 is not a string" {
		t.Errorf("ToString test failure")
	}
}
