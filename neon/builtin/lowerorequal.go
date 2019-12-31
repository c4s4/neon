package builtin

import (
	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "lowerorequal",
		Func: lowerOrEqual,
		Help: `Check that NeON version is lower or equal that given version.

Arguments:

- The version to check against.

Returns:

- A boolean telling if NeON version is lower or equal than given version.

Examples:

    # check that NeON version is lower or equal than 0.12.0
    lowerorequal("0.12.0")
    # return true if version is lower or equal than 0.12.0, false otherwise`,
	})
}

func lowerOrEqual(version string) bool {
	n, err := build.NewVersion(build.NeonVersion)
	if err != nil {
		panic("Neon version could not be parsed")
	}
	v, err := build.NewVersion(version)
	if err != nil {
		panic("Version could not be parsed")
	}
	return n.Compare(v) <= 0
}
