package builtin

import (
	"neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "jsonencode",
		Func: jsonEncode,
		Help: `Encode given variable in Json format.

Arguments:

- The variable to encode in Json format.

Returns:

- Json encoded string.

Examples:

    # encode given list
    jsonencode(["foo", "bar"])
    # returns: "['foo', 'bar']"`,
	})
}

func jsonEncode(object interface{}) string {
	bytes, err := build.PropertyToString(object, true)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
