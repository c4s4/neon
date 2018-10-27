package builtin

import (
	"github.com/c4s4/neon/build"
	"regexp"
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

- A boolean telling string matches regular expression.

Examples:

    # tell if string "neon" marchs "n..n" regular expression:
    match("n..n", "neon")
    # return true`,
	})
}

func match(r, s string) bool {
	m, _ := regexp.MatchString(r, s)
	return m
}
