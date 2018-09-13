package builtin

import (
	"encoding/json"
	"github.com/c4s4/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "jsondecode",
		Func: jsonDecode,
		Help: `Decode given string in Json format.

Arguments:

- The string in Json format to decode.

Returns:

- Decoded string.

Examples:

    # decode given list
    jsondecode("['foo', 'bar']")
    # returns string slice: ["foo", "bar"]`,
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
