package builtin

import (
	"gopkg.in/yaml.v2"
	"neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "yamldecode",
		Func: yamlDecode,
		Help: `Decode given string in YAML format.

Arguments:

- The string in YAML format to decode.

Returns:

- Decoded string.

Examples:

    # decode given list
    yamldecode("['foo', 'bar']")
    # returns string slice: ["foo", "bar"]`,
	})
}

func yamlDecode(encoded string) interface{} {
	var value interface{}
	err := yaml.Unmarshal([]byte(encoded), &value)
	if err != nil {
		panic(err)
	}
	return value
}
