package builtin

import (
	"neon/build"
	"strings"
)

func init() {
	build.BuiltinMap["join"] = build.BuiltinDescriptor{
		Function: Join,
		Help: `Join strings.

Arguments:

- The strings to join as a list of strings.
- The separator as a string.

Returns:

- Joined strings as a string.

Examples:

    // join "foo" and "bar" with a space
    join(["foo", "bar"], " ")`,
	}
}

func Join(elements []interface{}, separator string) string {
	var strs = make([]string, len(elements))
	for index, elt := range elements {
		str, ok := elt.(string)
		if !ok {
			panic("can only join strings")
		}
		strs[index] = str
	}
	return strings.Join(strs, separator)
}
