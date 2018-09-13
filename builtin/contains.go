package builtin

import (
	"github.com/c4s4/neon/build"
	"github.com/c4s4/neon/util"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "contains",
		Func: contains,
		Help: `Contains strings.

Arguments:

- List of strings to search into.
- Searched string.

Returns:

- A boolean telling if the string is contained in the list.

Examples:

    # Tell if the list contains "bar"
    contains(["foo", "bar"], "bar")
    # returns: true`,
	})
}

func contains(elements interface{}, value string) bool {
	slice, err := util.ToSliceString(elements)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(slice); i++ {
		if slice[i] == value {
			return true
		}
	}
	return false
}
