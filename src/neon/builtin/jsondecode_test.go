package builtin

import (
	"testing"
)

func TestJsonDecode(t *testing.T) {
	if jsonEncode(jsonDecode(`["foo", "bar"]`)) != `["foo", "bar"]` {
		t.Errorf("Error decoding json")
	}
}
