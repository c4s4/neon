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
	expected := "42"
	actual, err := Serialize(42)
	if err != nil || actual != expected {
		t.Error("Error serializing int", err)
	}
}

func TestSerializeFloat(t *testing.T) {
	expected := "4.2"
	actual, err := Serialize(4.2)
	if err != nil || actual != expected {
		t.Error("Error serializing float", err)
	}
}

func TestSerializeList(t *testing.T) {
	expected := "[1, 2, 3]"
	actual, err := Serialize([]int{1, 2, 3})
	if err != nil || actual != expected {
		t.Error("Error serializing slice", err)
	}
}
