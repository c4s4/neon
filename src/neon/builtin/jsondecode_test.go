package builtin

import (
	"testing"
)

func TestJsonDecode(t *testing.T) {
	if JsonEncode(JsonDecode(`["foo","bar"]`)) != `["foo","bar"]` {
		t.Errorf("Error decoding json")
	}
}
