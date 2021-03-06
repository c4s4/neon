package builtin

import (
	"github.com/c4s4/neon/neon/build"
	"reflect"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "list",
		Func: list,
		Help: `Return a list:
- If the object is already a list, return the object.
- If the object is not a list, wrap it into a list.

Arguments:

- The object to turn into a list.

Returns:

- The list.

Examples:

    # get a list of foo
    list(foo)
	# return foo if already a list or [foo] otherwise`,
	})
}

func list(object interface{}) []interface{} {
	value := reflect.ValueOf(object)
	if value.Kind() == reflect.Slice {
		return value.Interface().([]interface{})
	}
	slice := make([]interface{}, 1)
	slice[0] = object
	return slice
}
