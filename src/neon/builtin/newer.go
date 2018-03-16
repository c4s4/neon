package builtin

import (
	"neon/build"
	"os"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "newer",
		Func: newer,
		Help: `Tells if source is newer than result file (if any).

Arguments:

- source: source file that must exist.
- result: result file (may not exist).

Returns:

- A boolean that tells if source is newer than result. If result file doesn't
  exists, this returns true.

Examples:

    # generate PDF if source Markdown changed
    if newer("source.md", "result.pdf") {
    	compile("source.md")
    }`,
	})
}

func newer(source, result string) bool {
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
