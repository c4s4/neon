package builtin

import (
	"testing"
)

func TestReplace(t *testing.T) {
	Assert(replace("foo bar", "bar", "spam"), "foo spam", t)
	Assert(replace("foo bar bar", "bar", "spam"), "foo spam spam", t)
	Assert(replace("foo bar", "spam", "eggs"), "foo bar", t)
}
