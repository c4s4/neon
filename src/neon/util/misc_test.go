package util

import (
	"reflect"
	"testing"
)

func TestToString(t *testing.T) {
	actual, err := ToString("string")
	if err != nil || actual != "string" {
		t.Errorf("ToString test failure")
	}
	_, err = ToString(25)
	if err == nil || err.Error() != "25 is not a string" {
		t.Errorf("ToString test failure")
	}
	var str interface{} = "string"
	actual, err = ToString(str)
	if err != nil || actual != "string" {
		t.Errorf("ToString test failure")
	}
}

func TestToSliceString(t *testing.T) {
	expected := []string{"foo", "bar"}
	var object interface{} = []string{"foo", "bar"}
	actual, err := ToSliceString(object)
	if err != nil || !reflect.DeepEqual(actual, expected) {
		t.Errorf("ToSliceString test failure")
	}
}

func TestToMapStringString(t *testing.T) {
	expected := map[string]string{"foo": "bar"}
	var object interface{} = map[string]string{"foo": "bar"}
	actual, err := ToMapStringString(object)
	if err != nil || !reflect.DeepEqual(actual, expected) {
		t.Errorf("ToMapStringString test failure")
	}
}

func TestToMapStringInterface(t *testing.T) {
	expected := map[string]interface{}{"foo": "bar"}
	var object interface{} = map[string]interface{}{"foo": "bar"}
	actual, err := ToMapStringInterface(object)
	if err != nil || !reflect.DeepEqual(actual, expected) {
		t.Errorf("ToMapStringInterface test failure")
	}
}

func TestIsMap(t *testing.T) {
	var object interface{} = map[string]string{"foo": "bar"}
	if !IsMap(object) {
		t.Errorf("IsMAp test failure")
	}
	object = "foo"
	if IsMap(object) {
		t.Errorf("IsMAp test failure")
	}
}

func TestIsSlice(t *testing.T) {
	var object interface{} = []string{"foo", "bar"}
	if !IsSlice(object) {
		t.Errorf("IsSlice test failure")
	}
	object = "foo"
	if IsMap(object) {
		t.Errorf("IsSlice test failure")
	}
}

func TestMaxLineLength(t *testing.T) {
	lines := []string{"12345", "1234567890"}
	if MaxLineLength(lines) != 10 {
		t.Errorf("MaxLineLength test failure")
	}
}

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
