package builtin

import (
	"github.com/c4s4/neon/neon/build"
	"strings"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "uppercase",
		Func: uppercase,
		Help: `Put a string in upper case.

Arguments:

- The string to put in upper case.

Returns:

- The string in uppercase.

Examples:

    # set string in upper case
    uppercase("FooBAR")
    # returns: "FOOBAR"`,
	})
}

func uppercase(message string) string {
	return strings.ToUpper(message)
}
