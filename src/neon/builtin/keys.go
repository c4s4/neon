package builtin

import (
	"neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "keys",
		Func: Keys,
		Help: `Return keys of gien map.

Arguments:

- The map to get keys for.

Returns:

- A list of keys.

Examples:

    # get keys of a map
    keys({"foo": 1, "bar": 2})
    # returns: ["foo", "bar"]`,
	})
}

func Keys(themap map[interface{}]interface{}) []interface{} {
	keys := make([]interface{}, 0, len(themap))
	for key := range themap {
		keys = append(keys, key)
	}
	return keys
}
