package main

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
