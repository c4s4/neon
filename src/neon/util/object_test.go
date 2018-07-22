package util

import (
	"reflect"
	"testing"
)

func TestNewObject(t *testing.T) {
	var thing interface{} = 1
	_, err := NewObject(thing)
	if err == nil || err.Error() != "field must be a map with string keys" {
		t.Errorf("Bad object construction: '%s'", err.Error())
	}
}

func TestGetString(t *testing.T) {
	expected := "bar"
	object := Object{"foo": "bar"}
	actual, err := object.GetString("foo")
	if err != nil || actual != expected {
		t.Error("Error getting string", err)
	}
}

func TestGetBoolean(t *testing.T) {
	expected := true
	object := Object{"foo": true}
	actual, err := object.GetBoolean("foo")
	if err != nil || actual != expected {
		t.Error("Error getting boolean", err)
	}
}

func TestGetInteger(t *testing.T) {
	expected := 1
	object := Object{"foo": 1}
	actual, err := object.GetInteger("foo")
	if err != nil || actual != expected {
		t.Error("Error getting integer", err)
	}
}

func TestGetList(t *testing.T) {
	expected := []interface{}{1, 2, 3}
	object := Object{"foo": []int{1, 2, 3}}
	actual, err := object.GetList("foo")
	if err != nil || !reflect.DeepEqual(expected, actual) {
		t.Error("Error getting list", err)
	}
}

func TestGetListStrings(t *testing.T) {
	expected := []string{"spam", "eggs"}
	object := Object{"foo": []string{"spam", "eggs"}}
	actual, err := object.GetListStrings("foo")
	if err != nil || !reflect.DeepEqual(expected, actual) {
		t.Error("Error getting list of strings", err)
	}
}

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

func TestGetObject(t *testing.T) {
	expected := Object{"spam": "eggs"}
	object := Object{"foo": Object{"spam": "eggs"}}
	actual, err := object.GetObject("foo")
	if err != nil || !reflect.DeepEqual(expected, actual) {
		t.Error("Error getting object", err)
	}
}

func TestGetMapStringString(t *testing.T) {
	expected := map[string]string{"spam": "eggs"}
	object := Object{"foo": map[string]string{"spam": "eggs"}}
	actual, err := object.GetMapStringString("foo")
	if err != nil || !reflect.DeepEqual(expected, actual) {
		t.Error("Error getting map string string", err)
	}
}

func TestCheckFields(t *testing.T) {
	object := Object{"foo": "bar", "spam": "eggs"}
	err := object.CheckFields([]string{"foo"})
	if err == nil {
		t.Errorf("Error checking fields: %v", err)
	}
	err = object.CheckFields([]string{"foo", "spam"})
	if err != nil {
		t.Errorf("Error checking fields: %v", err)
	}
}

func TestCopy(t *testing.T) {
	object := Object{"foo": "bar", "spam": "eggs"}
	actual := object.Copy()
	if !reflect.DeepEqual(actual, object) {
		t.Errorf("Error copying objects")
	}
}

func TestFields(t *testing.T) {
	object := Object{"foo": "bar", "spam": "eggs"}
	fields := object.Fields()
	if !reflect.DeepEqual([]string{"foo", "spam"}, fields) {
		t.Errorf("Error getting fields")
	}
}

func TestHasField(t *testing.T) {
	object := Object{"foo": "bar", "spam": "eggs"}
	if !object.HasField("foo") {
		t.Errorf("Error checking field")
	}
	if object.HasField("google") {
		t.Errorf("Error checking field")
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
