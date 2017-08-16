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
