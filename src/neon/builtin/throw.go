package builtin

import (
	"fmt"
	"neon/build"
)

func init() {
	build.BuiltinMap["throw"] = build.BuiltinDescriptor{
		Function: Throw,
		Help:     "Throw an error and sets 'error' variable to its value",
	}
}

func Throw(message string) error {
	return fmt.Errorf(message)
}
