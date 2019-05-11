package builtin

import (
	"github.com/c4s4/neon/build"
	"github.com/c4s4/neon/util"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "findinpath",
		Func: findInPath,
		Help: `Find executables in PATH.

Arguments:

- The executable to find.

Returns:

- A list of absolute paths to the executable, in the order of the PATH.

Examples:

    # find python in path
    findinpath("python")
    # returns: ["/opt/python/current/bin/python", /usr/bin/python"]`,
	})
}

func findInPath(executable string) []string {
	return util.FindInPath(executable)
}
