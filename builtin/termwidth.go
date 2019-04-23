package builtin

import (
	"github.com/c4s4/neon/build"
	"github.com/c4s4/neon/util"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "termwidth",
		Func: termwidth,
		Help: `Return terminal width.

Arguments:

- None

Returns:

- Terminal width in characters.

Examples:

	# get terminal width
	width = termwidth()`,
	})
}

func termwidth() int {
	return util.TerminalWidth()
}
