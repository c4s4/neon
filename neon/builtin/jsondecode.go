package builtin

import (
	"encoding/json"

	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "jsondecode",
		Func: jsonDecode,
		Help: `Decode given string in JSON format.

Arguments:

- The string in JSON format to decode.

Returns:

- Decoded value.

Examples:

    # decode given list
    jsondecode('["foo", "bar"]')
    # returns list: ["foo", "bar"]`,
	})
}

func jsonDecode(encoded string) interface{} {
	var value interface{}
	err := json.Unmarshal([]byte(encoded), &value)
	if err != nil {
		panic(err)
	}
	return value
}
