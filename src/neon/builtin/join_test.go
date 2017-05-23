package builtin

import (
	"testing"
)

func TestJoin(t *testing.T) {
	if Join([]interface{}{"foo", "bar"}, " ") != "foo bar" {
		t.Errorf("Error builtin join")
	}
}
