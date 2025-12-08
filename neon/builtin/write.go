package builtin

import (
	"os"

	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "write",
		Func: write,
		Help: `Write a string in given file.

Arguments:

- The file name to write.
- The string to write.

Examples:

    # write "1.2.3" in VERSION file
    write("VERSION", "1.2.3")`,
	})
}

func write(file string, text string) {
	err := os.WriteFile(file, []byte(text), 0644)
	if err != nil {
		panic(err.Error())
	}
}
