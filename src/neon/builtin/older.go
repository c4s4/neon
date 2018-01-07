package builtin

import (
	"neon/build"
	"os"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "older",
		Func: Older,
		Help: `Tells if source is older than result file (if any).

Arguments:

- source: source file that must exist.
- result: result file (may not exist).

Returns:

- A boolean that tells if source is older tha result. If result file doesn't
  exists, this returns true.

Examples:

    // generate PDF if source Markdown changed
    if older("source.md", "result.pdf") {
    	compile("source.md")
    }`,
	})
}

func Older(source, result string) bool {
	info, err := os.Stat(source)
	if err != nil {
		panic("could no get info about source file")
	}
	sourceTime := info.ModTime()
	info, err = os.Stat(result)
	if err != nil {
		return true
	}
	resultTime := info.ModTime()
	return sourceTime.After(resultTime)
}
