package builtin

import (
	"io/ioutil"
	"neon/build"
	"strings"
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

    # read VERSION file and set variable version with ots content
    read("VERSION")
    # returns: the contents of "VERSION" file`,
	})
}

func read(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err.Error())
	}
	return strings.TrimSpace(string(content))
}
