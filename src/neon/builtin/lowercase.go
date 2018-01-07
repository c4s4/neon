package builtin

import (
	"neon/build"
	"strings"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "lowercase",
		Func: Lowercase,
		Help: `Put a string in lower case.

Arguments:

- The string to put in lower case.

Returns:

- The string in lower case.

Examples:

    // set string in lower case
    lowercase("FooBAR")
    // returns: "foobar"`,
	})
}

func Lowercase(message string) string {
	return strings.ToLower(message)
}
