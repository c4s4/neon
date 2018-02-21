package builtin

import (
	"testing"
)

func TestUnscapeUrl(t *testing.T) {
	if unescapeUrl("foo%20bar") != "foo bar" {
		t.Errorf("Error builtin unescapeulr")
	}
}
