package builtin

import (
	"neon/build"
	"strings"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc {
		Name: "split",
		Func: Split,
		Help: `Split strings.

Arguments:

- The strings to split.
- The separator for splitting.

Returns:

- A list of strings.

Examples:

    // split "foo bar" with space
    split("foo bar", " ")
    // returns: ["foo"," "bar"]`,
	})
}

func Split(str, sep string) []string {
	return strings.Split(str, sep)
}
