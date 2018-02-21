package builtin

import (
	"neon/build"
	"net/url"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "escapeurl",
		Func: escapeUrl,
		Help: `Escape given URL.

Arguments:

- The URL to escape.

Returns:

- The escaped URL.

Examples:

    # escape given URL
    escapeurl("/foo bar")
    # returns: "/foo%20bar"`,
	})
}

func escapeUrl(path string) string {
	url, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	return url.EscapedPath()
}
