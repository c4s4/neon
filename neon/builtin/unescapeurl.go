package builtin

import (
	"net/url"

	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "unescapeurl",
		Func: unescapeURL,
		Help: `Unescape given URL.

Arguments:

- The URL to unescape.

Returns:

- The unescaped URL.

Examples:

    # unescape given URL
    unescapeurl("foo%20bar")
    # returns: "foo bar"`,
	})
}

func unescapeURL(path string) string {
	unescaped, err := url.PathUnescape(path)
	if err != nil {
		panic(err)
	}
	return unescaped
}
