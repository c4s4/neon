package builtin

import (
	"testing"
)

func TestYamlEncode(t *testing.T) {
	if yamlEncode([]string{"foo", "bar"}) != `["foo", "bar"]` {
		t.Errorf("Error encoding yaml")
	}
}
