package builtin

import (
	"testing"
)

func TestJoin(t *testing.T) {
	if join([]string{"foo", "bar"}, " ") != "foo bar" {
		t.Errorf("Error builtin join")
	}
}
