package builtin

import (
	"testing"

	"github.com/c4s4/neon/neon/util"
)

func TestTitle(t *testing.T) {
	text := title("Test", "#")
	Assert(text[:7], "## Test", t)
	Assert(len(text), util.TerminalWidth(), t)
	Assert(text[len(text)-2:], "##", t)
}
