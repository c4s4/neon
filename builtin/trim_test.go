package builtin

import (
	"testing"
)

func TestTrim(t *testing.T) {
	Assert(trim("\tfoo bar\n   "), "foo bar", t)
	Assert(trim("foo bar"), "foo bar", t)
}
