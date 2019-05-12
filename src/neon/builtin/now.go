package builtin

import (
	"neon/build"
	"time"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "now",
		Func: now,
		Help: `Return current date and time in ISO format.

Arguments:

- none

Returns:

- ISO date and time as a string.

Examples:

    # put current date and time in dt variable
    now()
    # returns: "2006-01-02 15:04:05"
    # to get date in ISO format
    now()[0:10]
    # returns: "2006-01-02"`,
	})
}

func now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
