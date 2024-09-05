package builtin

import (
	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "termlink",
		Func: termlink,
		Help: `Make given string a link when printed in a terminal.

Arguments:

- The URL to link to.
- The text to print on terminal.

Note: if text is empty, the URL is used as text.

Examples:

    # make link for string "Example" to "https://example.com"
    termlink("https://example.com", "Example")
    # returns: string "Example" with escape codes for terminal link to "https://example.com"`,
	})
}

func termlink(url, text string) string {
	if text == "" {
		text = url
	}
	return "\033]8;;" + url + "\033\\" + text + "\033]8;;\033\\"
}
