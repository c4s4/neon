package builtin

import (
	"neon/build"
	"strings"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc {
		Name: "uppercase",
		Func: Uppercase,
		Help: `Put a string in upper case.

Arguments:

- The string to put in upper case.

Returns:

- The string in uppercase.

Examples:

    // set string in upper case
    uppercase("FooBAR")
    // returns: "FOOBAR"`,
	})
}

func Uppercase(message string) string {
	return strings.ToUpper(message)
}
