package builtin

import (
	"neon/build"
	"net/url"
)

func init() {
	build.BuiltinMap["escapeurl"] = build.BuiltinDescriptor{
		Function: EscapeUrl,
		Help: `Escape given URL.

Arguments:

- The URL to escape.

Returns:

- The escaped URL.

Examples:

    // escape given URL
    escapeurl("foo bar")
    // returns: "foo%20bar"`,
	}
}

func EscapeUrl(path string) string {
	return url.PathEscape(path)
}
