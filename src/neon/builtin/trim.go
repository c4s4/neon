package builtin

import (
	"neon/build"
	"strings"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "trim",
		Func: trim,
		Help: `Trim spaces from given string.

Arguments:

- The string to trim.

Returns:

- Trimed string.

Examples:

    # trim string "\tfoo bar\n   "
    trim("\tfoo bar\n  ")
    # returns: "foo bar"`,
	})
}

func trim(str string) string {
	return strings.TrimSpace(str)
}
