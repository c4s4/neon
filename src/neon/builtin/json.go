package builtin

import (
	"encoding/json"
	"neon/build"
)

func init() {
	build.BuiltinMap["json"] = build.BuiltinDescriptor{
		Function: Json,
		Help: `Json encodin.

Arguments:

- The variable to encode in Json format.

Returns:

- Json encoded string.

Examples:

    // encode given list
    json(["foo", "bar"])
    // returns: "['foo', 'bar']"`,
	}
}

func Json(object []interface{}) string {
	bytes, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
