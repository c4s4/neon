package builtin

import (
	"neon/build"

	"github.com/mcuadros/go-version"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "sortversions",
		Func: sortVersions,
		Help: `Sort a list of versions.

Arguments:

- The list of versions to sort.

Returns:

- nothing but slice of versions is sorted

Examples:

    # sort version ["1.10", "1.1", "1.2"]
    sortversions(["1.10", "1.1", "1.2"])
    # returns nothing but slice of versions is sorted`,
	})
}

func sortVersions(versions []string) {
	version.Sort(versions)
}
