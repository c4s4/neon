package builtin

import (
	"os"
	"time"

	"github.com/c4s4/neon/neon/build"
	"github.com/c4s4/neon/neon/util"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "newer",
		Func: newer,
		Help: `Tells if source files are newer than result ones.

Arguments:

- sources: source file(s) (may not exist).
- results: result file(s) (may not exist).

Returns:

- A boolean that tells if source files are newer than result ones.
  If source files don't exist, returns false.
  If result files don't exist, returns true.

Examples:

    # generate PDF if source Markdown changed
    if newer("source.md", "result.pdf") {
    	compile("source.md")
    }
	# generate binary if source files are newer than generated binary
    if newer(find(".", "**/*.go"), "bin/binary") {
    	generateBinary()
    }`,
	})
}

func newer(sources, results interface{}) bool {
	if sources == nil {
		return false
	}
	if results == nil {
		return true
	}
	sourceFiles, err := util.ToSliceString(sources)
	if err != nil {
		panic("source must be a string or list of strings")
	}
	resultFiles, err := util.ToSliceString(results)
	if err != nil {
		panic("result must be a string or list of strings")
	}
	var sourceTime time.Time
	for _, source := range sourceFiles {
		info, err := os.Stat(source)
		if err != nil {
			continue
		}
		t := info.ModTime()
		if sourceTime.IsZero() || t.After(sourceTime) {
			sourceTime = t
		}
	}
	if sourceTime.IsZero() {
		return false
	}
	var resultTime time.Time
	for _, result := range resultFiles {
		info, err := os.Stat(result)
		if err != nil {
			continue
		}
		t := info.ModTime()
		if resultTime.IsZero() || t.Before(resultTime) {
			resultTime = t
		}
	}
	if resultTime.IsZero() {
		return true
	}
	return sourceTime.After(resultTime)
}
