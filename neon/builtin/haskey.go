package builtin

import (
	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "haskey",
		Func: haskey,
		Help: `Tells if a map contains given key.

Arguments:

- The map to test.
- The key to test.

Returns:

- A boolean telling if the map contains given key.

Examples:

    # Tell if map "map" contains key "key"
    haskey(map, "key")
    # returns: true or false`,
	})
}

func haskey(themap map[interface{}]interface{}, key interface{}) bool {
	_, ok := themap[key]
	return ok
}
