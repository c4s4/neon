package builtin

import (
	"encoding/json"
	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "jsonindent",
		Func: jsonIndent,
		Help: `Indent given string in Json format.

Arguments:

- The string in Json format to indent.

Returns:

- Indented JSON string.

Examples:

    # indent given list
    jsonindent("['foo', 'bar']")
    # returns string:
    # [
    #   "foo",
    #   "bar"
    # ]`,
	})
}

func jsonIndent(encoded string) interface{} {
	var value interface{}
	err := json.Unmarshal([]byte(encoded), &value)
	if err != nil {
		panic(err)
	}
	bytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
