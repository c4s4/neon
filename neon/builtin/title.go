package builtin

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/c4s4/neon/neon/build"
	"github.com/c4s4/neon/neon/util"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "title",
		Func: title,
		Help: `Return a title string.

Arguments:

- The title text as a string.
- The title character as a string.

Returns:

- Title text surrounded by title separators with terminal width.

Examples:

    # get title "Hello World!"
    title("Hello World!", "#")
    # returns: "## Hello World! #############################################"`,
	})
}

func title(text, separator string) string {
	length := util.TerminalWidth() - (4 + utf8.RuneCountInString(text))
	if length < 2 {
		length = 2
	}
	return fmt.Sprintf("%s %s %s", strings.Repeat(separator, 2), text, strings.Repeat(separator, length))
}
