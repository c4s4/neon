package builtin

import (
	"testing"
)

func TestKeys(t *testing.T) {
	keys := Keys(map[interface{}]interface{}{"foo": 1, "bar": 2})
	if len(keys) != 2 {
		t.Errorf("Error builtin keys")
	}
	if keys[0] != "foo" && keys[0] != "bar" {
		t.Errorf("Error builtin keys")
	}
}
