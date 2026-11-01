package builtin

import (
	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "lowerorequal",
		Func: lowerOrEqual,
		Help: `Check that NeON version is lower than or equal to given version.

Arguments:

- The version to check against.

Returns:

- A boolean telling if NeON version is lower than or equal to given version.

Examples:

    # check that NeON version is lower than or equal to 0.12.0
    lowerorequal("0.12.0")
    # returns true if version is lower than or equal to 0.12.0, false otherwise`,
	})
}

func lowerOrEqual(version string) bool {
	return build.CompareNeonVersion(version) <= 0
}
