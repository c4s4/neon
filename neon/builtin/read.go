package builtin

import (
	"os"
	"strings"

	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "read",
		Func: read,
		Help: `Read given file and return its content as a string.

Arguments:

- The file name to read.

Returns:

- The file content as a string.

Examples:

    # read content of VERSION file
    read("VERSION")
    # returns: the contents of VERSION file`,
	})
}

func read(file string) string {
	content, err := os.ReadFile(file)
	if err != nil {
		panic(err.Error())
	}
	return strings.TrimSpace(string(content))
}
