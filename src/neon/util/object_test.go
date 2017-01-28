package util

import (
	"reflect"
	"testing"
)

func TestGetListStringsOrString(t *testing.T) {
	expected := []string{"bar"}
	object := Object{"foo": "bar"}
	actual, err := object.GetListStringsOrString("foo")
	if err != nil || !reflect.DeepEqual(expected, actual) {
		t.Error("Error getting string", err)
	}
}

func TestGetListStringsOrStringEmpty(t *testing.T) {
	object := Object{"foo": "bar"}
	actual, err := object.GetListStringsOrString("bar")
	if err != nil || !reflect.DeepEqual(actual, []string{}) {
		t.Error("No error getting string!", actual)
	}
}

func TestToMapStringString(t *testing.T) {
	object := Object{"foo": "1", "bar": "2"}
	actual, err := object.ToMapStringString()
	if err != nil {
		t.Error("Error getting the map string string")
	}
	if actual["foo"] != "1" {
		t.Error("Error getting the map string string")
	}
	if actual["bar"] != "2" {
		t.Error("Error getting the map string string")
	}
}
