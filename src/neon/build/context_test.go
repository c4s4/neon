package build

import (
	"testing"
)

func TestPropertyToStringString(t *testing.T) {
	expected := `"test"`
	actual, err := PropertyToString("test", true)
	if err != nil || actual != expected {
		t.Error("Error serializing string", err)
	}
}

func TestPropertyToStringInt(t *testing.T) {
	expected := `42`
	actual, err := PropertyToString(42, true)
	if err != nil || actual != expected {
		t.Error("Error serializing int", err)
	}
}

func TestPropertyToStringFloat(t *testing.T) {
	expected := `4.2`
	actual, err := PropertyToString(4.2, true)
	if err != nil || actual != expected {
		t.Error("Error serializing float", err)
	}
}

func TestPropertyToStringList(t *testing.T) {
	expected := `[1, 2, 3]`
	actual, err := PropertyToString([]int{1, 2, 3}, true)
	if err != nil || actual != expected {
		t.Error("Error serializing slice", err)
	}
}

func TestPropertyToStringCompositeList(t *testing.T) {
	expected := `[1, 2, 3, "spam"]`
	actual, err := PropertyToString([]interface{}{1, 2, 3, "spam"}, true)
	if err != nil || actual != expected {
		t.Error("Error serializing composite slice", err)
	}
}

func TestPropertyToStringMap(t *testing.T) {
	expected := `["bar": 2, "foo": 1]`
	actual, err := PropertyToString(map[string]int{"foo": 1, "bar": 2}, true)
	if err != nil || actual != expected {
		t.Error("Error serializing map", err)
	}
}

func TestPropertyToStringCompositeMap(t *testing.T) {
	expected := `["bar": 2, "foo": 1, 3: "spam"]`
	actual, err := PropertyToString(map[interface{}]interface{}{"foo": 1, "bar": 2, 3: "spam"}, true)
	if err != nil || actual != expected {
		t.Error("Error serializing composite map", err)
	}
}
