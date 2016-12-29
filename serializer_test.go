package main

import (
	"testing"
)

func TestSerializeString(t *testing.T) {
	expected := `"test"`
	actual, err := Serialize("test")
	if err != nil || actual != expected {
		t.Error("Error serializing string", err)
	}
}

func TestSerializeInt(t *testing.T) {
	expected := `42`
	actual, err := Serialize(42)
	if err != nil || actual != expected {
		t.Error("Error serializing int", err)
	}
}

func TestSerializeFloat(t *testing.T) {
	expected := `4.2`
	actual, err := Serialize(4.2)
	if err != nil || actual != expected {
		t.Error("Error serializing float", err)
	}
}

func TestSerializeList(t *testing.T) {
	expected := `[1, 2, 3]`
	actual, err := Serialize([]int{1, 2, 3})
	if err != nil || actual != expected {
		t.Error("Error serializing slice", err)
	}
}

func TestSerializeCompositeList(t *testing.T) {
	expected := `[1, 2, 3, "spam"]`
	actual, err := Serialize([]interface{}{1, 2, 3, "spam"})
	if err != nil || actual != expected {
		t.Error("Error serializing composite slice", err)
	}
}

func TestSerializeMap(t *testing.T) {
	expected := `["bar": 2, "foo": 1]`
	actual, err := Serialize(map[string]int{"foo": 1, "bar": 2})
	if err != nil || actual != expected {
		t.Error("Error serializing map", err)
	}
}

func TestSerializeCompositeMap(t *testing.T) {
	expected := `["bar": 2, "foo": 1, 3: "spam"]`
	actual, err := Serialize(map[interface{}]interface{}{"foo": 1, "bar": 2, 3: "spam"})
	if err != nil || actual != expected {
		t.Error("Error serializing composite map", err)
	}
}
