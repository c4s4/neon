package builtin

import (
	"neon/build"
	"strings"
)

func init() {
	build.BuiltinMap["uppercase"] = build.BuiltinDescriptor{
		Function: Uppercase,
		Help: `Put a string in upper case.

Arguments:
- The string to put in upper case.
Returns:
- The string in uppercase.

Examples:
// set greetings in upper case
greetings = uppercase(greetings)`,
	}
}

func Uppercase(message string) string {
	return strings.ToUpper(message)
}
