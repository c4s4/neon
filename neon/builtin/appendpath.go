package builtin

import (
	"github.com/c4s4/neon/neon/build"
	"path/filepath"
	"reflect"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "appendpath",
		Func: appendPath,
		Help: `Append root directory to paths.

Arguments:

- The root directory.
- The paths to append.

Returns:

- Appended paths as a list.

Examples:

    # append root "foo" to paths "spam" and "eggs"
    appendpath("foo", "spam", "eggs")
	# returns: ["foo/spam", "foo/eggs"] on Linux and
	# ["foo\spam", "foo\eggs"] on Windows`,
	})
}

func appendPath(root string, paths []string) []string {
	p := make([]string, len(paths))
	for i, e := range paths {
		s := reflect.ValueOf(e).String()
		p[i] = filepath.Join(root, s)
	}
	return p
}
