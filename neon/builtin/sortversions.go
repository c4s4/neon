package builtin

import (
	"github.com/c4s4/neon/neon/build"
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

Notes:

- A leading 'v' (as in git tags such as 'v1.2.3') is ignored.
- Suffixes are ordered: SNAPSHOT < ALPHA < BETA < RC < release.
- Strings that are not valid versions are placed first.

Examples:

    # sort version ["1.10", "1.1", "1.2"]
    sortversions(["1.10", "1.1", "1.2"])
    # returns nothing but slice of versions is sorted`,
	})
}

func sortVersions(versions []string) {
	build.SortVersions(versions)
}
