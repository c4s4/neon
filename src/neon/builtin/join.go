package builtin

import (
	"github.com/c4s4/neon/build"
	"github.com/c4s4/neon/util"
	"strings"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "join",
		Func: join,
		Help: `Join strings.

Arguments:

- The strings to join as a list of strings.
- The separator as a string.

Returns:

- Joined strings as a string.

Examples:

    # join "foo" and "bar" with a space
    join(["foo", "bar"], " ")
    # returns: "foo bar"`,
	})
}

func join(elements interface{}, separator string) string {
	elementsStrings, err := util.ToSliceString(elements)
	if err != nil {
		panic("Join first argument must ba a list of strings")
	}
	return strings.Join(elementsStrings, separator)
}
