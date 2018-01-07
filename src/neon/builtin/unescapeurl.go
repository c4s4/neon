package builtin

import (
	"neon/build"
	"net/url"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "unescapeurl",
		Func: UnescapeUrl,
		Help: `Unescape given URL.

Arguments:

- The URL to unescape.

Returns:

- The unescaped URL.

Examples:

    // unescape given URL
    escapeurl("foo%20bar")
    // returns: "foo bar"`,
	})
}

func UnescapeUrl(path string) string {
	unescaped, err := url.PathUnescape(path)
	if err != nil {
		panic(err)
	}
	return unescaped
}
