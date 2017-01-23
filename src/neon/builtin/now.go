package builtin

import (
	"neon/build"
	"time"
)

func init() {
	build.BuiltinMap["now"] = build.BuiltinDescriptor{
		Function: Now,
		Help:     "Return current date and time in ISO format",
	}
}

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
