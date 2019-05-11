package builtin

import (
	"github.com/c4s4/neon/build"
	"strings"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "replace",
		Func: replace,
		Help: `Replace string with another.

Arguments:

- The strings where take place replacements.
- The substring to replace.
- The replacement substring.

Returns:

- Replaced string.

Examples:

    # replace "foo" with "bar" in string "spam foo eggs"
    replace("spam foo eggs", "foo", "bar")
    # returns: "spam bar eggs"`,
	})
}

func replace(str, from, to string) string {
	return strings.Replace(str, from, to, -1)
}
