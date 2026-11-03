package builtin

import (
	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "greater",
		Func: greater,
		Help: `Check that NeON version is greater than given version.

Arguments:

- The version to check against.

Returns:

- A boolean telling if NeON version is greater than given version.

Examples:

    # check that NeON version is greater than 0.12.0
    greater("0.12.0")
    # returns true if version is greater than 0.12.0, false otherwise`,
	})
}

func greater(version string) bool {
	return build.CompareNeonVersion(version) > 0
}
