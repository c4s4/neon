package builtin

import (
	"io/ioutil"
	"neon/build"
)

func init() {
	build.BuiltinMap["read"] = build.BuiltinDescriptor{
		Function: Read,
		Help: `Read given file and return its content as a string.

Arguments:

- The file name to read.

Returns:

- The file content as a string.

Examples:

    // read VERSION file and set variable version with ots content
    version = read("VERSION")`,
	}
}

func Read(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err.Error())
	}
	return string(content)
}
