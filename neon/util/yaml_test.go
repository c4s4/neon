package util

import (
	"reflect"
	"testing"
)

func TestDecodeYAML(t *testing.T) {
	source := `
defaults: &defaults
  a: 1
  b: two
merged:
  <<: *defaults
  b: three
list: [1, 2.5, ~, true]
yes: yes
date: 2026-11-01
nested:
  key: value
`
	value, err := DecodeYAML([]byte(source))
	if err != nil {
		t.Fatalf("decoding YAML: %v", err)
	}
	expected := map[interface{}]interface{}{
		"defaults": map[interface{}]interface{}{"a": 1, "b": "two"},
		"merged":   map[interface{}]interface{}{"a": 1, "b": "three"},
		"list":     []interface{}{1, 2.5, nil, true},
		"yes":      "yes",
		"date":     "2026-11-01",
		"nested":   map[interface{}]interface{}{"key": "value"},
	}
	if !reflect.DeepEqual(value, expected) {
		t.Errorf("bad decoded value: %#v", value)
	}
}

func TestDecodeYAMLObject(t *testing.T) {
	for _, source := range []string{"", "  \n", "# comment\n"} {
		object, err := DecodeYAMLObject([]byte(source))
		if err != nil || object == nil || len(object) != 0 {
			t.Errorf("empty source '%s' should give empty object: %v (%v)", source, object, err)
		}
	}
	if _, err := DecodeYAMLObject([]byte("[1, 2]")); err == nil {
		t.Errorf("list should not be accepted as object")
	}
	if _, err := DecodeYAMLObject([]byte("1: foo")); err == nil {
		t.Errorf("non string keys should not be accepted")
	}
	object, err := DecodeYAMLObject([]byte("foo: {bar: 1}"))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := object["foo"].(map[interface{}]interface{}); !ok {
		t.Errorf("nested maps should be map[interface{}]interface{}: %#v", object["foo"])
	}
}
