package builtin

import (
	"neon/build"
	"neon/util"
	"os"
	"path/filepath"
	"strings"
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
	path := os.Getenv("PATH")
	dirs := strings.Split(path, string(os.PathListSeparator))
	var paths []string
	for _, dir := range dirs {
		file := filepath.Join(dir, executable)
		if util.FileIsExecutable(file) {
			paths = append(paths, file)
		}
	}
	return paths
}
