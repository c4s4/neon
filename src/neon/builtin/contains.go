package builtin

import (
	"neon/build"
	"reflect"
)

func init() {
	build.BuiltinMap["contains"] = build.BuiltinDescriptor{
		Function: Contains,
		Help: `Contains strings.

Arguments:

- List of strings to search into.
- Seached string.

Returns:

- A boolean telling if the string is contained in the list.

Examples:

    // Tell if the list contains "bar"
    contains(["foo", "bar"], "bar")
    // returns: true`,
	}
}

func Contains(elements interface{}, s string) bool {
	slice := reflect.ValueOf(elements)
	if slice.Kind() == reflect.Interface {
		slice = slice.Elem()
	}
	if slice.Kind() != reflect.Slice {
		panic("Contains first argument must ba a list of strings")
	}
	for i := 0; i < slice.Len(); i++ {
		value := slice.Index(i)
		if value.Kind() == reflect.Interface {
			value = value.Elem()
		}
		if value.Kind() != reflect.String {
			panic("Contains first argument must ba a list of strings")
		}
		if value.String() == s {
			return true
		}
	}
	return false
}
