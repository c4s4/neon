package builtin

import (
	"testing"
)

func TestYamlDecode(t *testing.T) {
	if yamlEncode(yamlDecode(`["foo", "bar"]`)) != `["foo", "bar"]` {
		t.Errorf("Error decoding yaml")
	}
}
