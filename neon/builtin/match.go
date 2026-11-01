package builtin

import (
	"regexp"

	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "match",
		Func: match,
		Help: `Tell if given string matches a regular expression.

Arguments:

- The regular expression.
- The string to test.

Returns:

- A boolean telling if string matches regular expression.

Examples:

    # tell if string "neon" matches "n..n" regular expression:
    match("n..n", "neon")
    # returns: true`,
	})
}

func match(r, s string) bool {
	m, _ := regexp.MatchString(r, s)
	return m
}
