package builtin

import (
	"neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "greater",
		Func: greater,
		Help: `Check that NeON version is greater that given version.

Arguments:

- The version to check against.

Returns:

- A boolean telling if NeON version is greater than given version.

Examples:

    # check that NeON version is greater than 0.12.0
    greater("0.12.0")
    # return true if version is greater than 0.12.0, false otherwise`,
	})
}

func greater(version string) bool {
	n, err := build.NewVersion(build.NeonVersion)
	if err != nil {
		panic("Neon version could not be parsed")
	}
	v, err := build.NewVersion(version)
	if err != nil {
		panic("Version could not be parsed")
	}
	return n.Compare(v) > 0
}
