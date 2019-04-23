package builtin

import (
	"github.com/c4s4/neon/build"
	"unicode/utf8"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "length",
		Func: length,
		Help: `Return length of given string.

Arguments:

- The string to get length for.

Returns:

- Length of the given string.

Examples:

    # get length of the string
    l = length("Hello World!")`,
	})
}

func length(str string) int {
	return utf8.RuneCountInString(str)
}
