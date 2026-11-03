package builtin

import (
	"github.com/c4s4/neon/neon/build"
	"github.com/c4s4/neon/neon/util"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "yamldecode",
		Func: yamlDecode,
		Help: `Decode given string in YAML format.

Arguments:

- The string in YAML format to decode.

Returns:

- Decoded value.

Examples:

    # decode given list
    yamldecode("['foo', 'bar']")
    # returns list: ["foo", "bar"]`,
	})
}

func yamlDecode(encoded string) interface{} {
	value, err := util.DecodeYAML([]byte(encoded))
	if err != nil {
		panic(err)
	}
	return value
}
