package builtin

import (
	"fmt"
)

func init() {
	Builtins["throw"] = BuiltinDescriptor{
		Function: Throw,
		Help:     "Throw an error and sets 'error' variable to its value",
	}
}

func Throw(message string) error {
	return fmt.Errorf(message)
}
