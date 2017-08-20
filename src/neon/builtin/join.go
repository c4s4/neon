package builtin

import (
	"neon/build"
	"reflect"
	"strings"
)

func init() {
	build.BuiltinMap["join"] = build.BuiltinDescriptor{
		Function: Join,
		Help: `Join strings.

Arguments:

- The strings to join as a list of strings.
- The separator as a string.

Returns:

- Joined strings as a string.

Examples:

    // join "foo" and "bar" with a space
    join(["foo", "bar"], " ")
    // returns: "foo bar"`,
	}
}

func Join(elements interface{}, separator string) string {
	slice := reflect.ValueOf(elements)
	if slice.Kind() == reflect.Interface {
		slice = slice.Elem()
	}
	if slice.Kind() != reflect.Slice {
		panic("Join first argument must ba a list of strings")
	}
	result := make([]string, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		value := slice.Index(i)
		if value.Kind() == reflect.Interface {
			value = value.Elem()
		}
		if value.Kind() != reflect.String {
			panic("Join first argument must ba a list of strings")
		}
		result[i] = value.String()
	}
	return strings.Join(result, separator)
}
