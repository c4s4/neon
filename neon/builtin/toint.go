package builtin

import (
	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "toint",
		Func: toint,
		Help: `Converts int64 value to int.

Arguments:

- The int64 value to convert.

Returns:

- Value converted to int.

Examples:

    # convert len([1, 2, 3]) to int
    toint(len([1, 2, 3]))
    # returns: 3`,
	})
}

func toint(i int64) int {
	return int(i)
}
