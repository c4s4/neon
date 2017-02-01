package builtin

import (
	"neon/build"
	"strings"
)

func init() {
	build.BuiltinMap["lowercase"] = build.BuiltinDescriptor{
		Function: Lowercase,
		Help: `Put a string in lower case.

Arguments:
- The string to put in lower case.
Returns:
- The string in lower case.

Examples:
// set greetings in lower case
greetings = lowercase(greetings)`,
	}
}

func Lowercase(message string) string {
	return strings.ToLower(message)
}
