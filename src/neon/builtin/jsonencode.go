package builtin

import (
	"encoding/json"
	"neon/build"
)

func init() {
	build.BuiltinMap["jsonencode"] = build.BuiltinDescriptor{
		Function: JsonEncode,
		Help: `Encode given variable in Json format.

Arguments:

- The variable to encode in Json format.

Returns:

- Json encoded string.

Examples:

    // encode given list
    jsonencode(["foo", "bar"])
    // returns: "['foo', 'bar']"`,
	}
}

func JsonEncode(object interface{}) string {
	bytes, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
