package builtin

import (
	"github.com/c4s4/neon/neon/build"
	"strings"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "lowercase",
		Func: lowercase,
		Help: `Put a string in lower case.

Arguments:

- The string to put in lower case.

Returns:

- The string in lower case.

Examples:

    # set string in lower case
    lowercase("FooBAR")
    # returns: "foobar"`,
	})
}

func lowercase(message string) string {
	return strings.ToLower(message)
}
