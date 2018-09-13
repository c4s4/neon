package builtin

import (
	"neon/build"
	"os"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "env",
		Func: env,
		Help: `Get environment variable.

Arguments:

- The name of the environment variable to get value for.

Returns:

- The value of this environment variable.

Examples:

    # get PATH environment variable
    env("PATH")
    # returns: value of the environment variable PATH`,
	})
}

func env(variable string) string {
	return os.Getenv(variable)
}
