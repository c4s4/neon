package builtin

import (
	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "lower",
		Func: lower,
		Help: `Check that NeON version is lower than given version.

Arguments:

- The version to check against.

Returns:

- A boolean telling if NeON version is lower than given version.

Examples:

    # check that NeON version is lower than 0.12.0
    lower("0.12.0")
    # returns true if version is lower than 0.12.0, false otherwise`,
	})
}

func lower(version string) bool {
	return build.CompareNeonVersion(version) < 0
}
