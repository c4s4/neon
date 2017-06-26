package builtin

import (
	"testing"
)

func TestJsonEncode(t *testing.T) {
	if JsonEncode([]string{"foo", "bar"}) != `["foo", "bar"]` {
		t.Errorf("Error encoding json")
	}
}
