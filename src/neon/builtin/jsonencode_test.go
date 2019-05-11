package builtin

import (
	"testing"
)

func TestJsonEncode(t *testing.T) {
	if jsonEncode([]string{"foo", "bar"}) != `["foo", "bar"]` {
		t.Errorf("Error encoding json")
	}
}
