package builtin

import (
	"github.com/c4s4/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "yamlencode",
		Func: yamlEncode,
		Help: `Encode given variable in YAML format.

Arguments:

- The variable to encode in YAML format.

Returns:

- Json encoded string.

Examples:

    # encode given list
    yamlencode(["foo", "bar"])
    # returns: "['foo', 'bar']"`,
	})
}

func yamlEncode(object interface{}) string {
	bytes, err := build.PropertyToString(object, true)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
