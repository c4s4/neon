package builtin

import (
	"neon/build"
	"time"
)

func init() {
	build.BuiltinMap["now"] = build.BuiltinDescriptor{
		Function: Now,
		Help: `Return current date and time in ISO format.

Arguments:
- none
Returns:
- ISO date ans time as a string.

Examples:
// put current date and time in dt variable
dt = now()`,
	}
}

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
