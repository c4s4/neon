package task

import (
	"fmt"
	"github.com/c4s4/changelog/lib"
	"github.com/c4s4/neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "changelog",
		Func: changelog,
		Args: reflect.TypeOf(changelogArgs{}),
		Help: `Load information from semantic changelog file.

Arguments:

- changelog: the name of the changelog file (look for changelog in current
  directory if omitted).

Note:

- The release version is stored in property _changelog_version.
- The release date is stored in property _changelog_date.
- The release summary is stored in property _changelog_summary.

Examples:

    # get changelog information in file 'test.yml':
    - changelog: "test.yml"`,
	})
}

type changelogArgs struct {
	Changelog string `neon:"file,optional"`
}

func changelog(context *build.Context, args interface{}) error {
	params := args.(changelogArgs)
	var file string
	if params.Changelog == "" {
		var err error
		file, err = lib.FindChangelog()
		if err != nil {
			return err
		}
	} else {
		file = params.Changelog
	}
	source, err := lib.ReadChangelog(file)
	if err != nil {
		return err
	}
	changelog, err := lib.ParseChangelog(source)
	if err != nil {
		return err
	}
	if len(changelog) < 1 {
		return fmt.Errorf("the changelog is empty")
	}
	release := changelog[0]
	context.SetProperty("_changelog_version", release.Version)
	context.SetProperty("_changelog_date", release.Date)
	context.SetProperty("_changelog_summary", release.Summary)
	return nil
}
