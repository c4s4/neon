package builtin

import (
	"unicode/utf8"

	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "length",
		Func: length,
		Help: `Return length of given string.

Arguments:

- The string to get length for.

Returns:

- Length of the given string, in characters.

Examples:

    # get length of the string
    l = length("Hello World!")
    # returns: 12`,
	})
}

func length(str string) int {
	return utf8.RuneCountInString(str)
}
