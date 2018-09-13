package builtin

import (
	"testing"
)

func TestWinexe(t *testing.T) {
	Assert(toWindows("foo/command"), `foo\command.exe`, t)
	Assert(toWindows("bar/command.sh"), `bar\command.bat`, t)
}
