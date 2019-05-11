package builtin

import (
	"fmt"

	"github.com/c4s4/changelog/lib"
	"github.com/c4s4/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "changelog",
		Func: changelog,
		Help: `Return changelog information from file.

Arguments:

- changelog: the name of the changelog file (look for changelog in current
  directory if empty string).

Note:

- The function returns a Changelog that is a list of Releases struct with
  fields Version, Date and Summary.

Examples:

    # get version of last release:
    - 'VERSION = changelog("")[0].Version'`,
	})
}

func changelog(file string) lib.Changelog {
	if file == "" {
		var err error
		file, err = lib.FindChangelog()
		if err != nil {
			panic(err)
		}
	}
	source, err := lib.ReadChangelog(file)
	if err != nil {
		panic(err)
	}
	changelog, err := lib.ParseChangelog(source)
	if err != nil {
		panic(err)
	}
	if len(changelog) < 1 {
		panic(fmt.Errorf("the changelog is empty"))
	}
	return changelog
}
