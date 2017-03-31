package util

import (
	"reflect"
	"testing"
)

func TestNewVersion(t *testing.T) {
	expected := Version{
		Name:    "1.2.3",
		Numbers: []int{1, 2, 3},
	}
	actual, err := NewVersion("1.2.3")
	if err != nil {
		t.Error("Error getting version", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Bad version", actual)
	}
}

func TestLess(t *testing.T) {
	version, _ := NewVersion("1")
	other, _ := NewVersion("2")
	if !version.Less(other) {
		t.Error("Bad comparison")
	}
	version, other = other, version
	if version.Less(other) {
		t.Error("Bad comparison")
	}
	version, _ = NewVersion("1.1")
	other, _ = NewVersion("1.2")
	if !version.Less(other) {
		t.Error("Bad comparison")
	}
	version, other = other, version
	if version.Less(other) {
		t.Error("Bad comparison")
	}
}
