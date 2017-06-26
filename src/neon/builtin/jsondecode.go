package builtin

import (
	"encoding/json"
	"neon/build"
)

func init() {
	build.BuiltinMap["jsondecode"] = build.BuiltinDescriptor{
		Function: JsonDecode,
		Help: `Decode given string in Json format.

Arguments:

- The string in Json format to decode.

Returns:

- Decoded string.

Examples:

    // decode given list
    jsondecode("['foo', 'bar']")
    // returns string slice: ["foo", "bar"]`,
	}
}

func JsonDecode(encoded string) interface{} {
	var value interface{}
	err := json.Unmarshal([]byte(encoded), &value)
	if err != nil {
		panic(err)
	}
	return value
}
