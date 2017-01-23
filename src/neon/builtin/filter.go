package builtin

import (
	zglob "github.com/mattn/go-zglob"
	"neon/build"
)

func init() {
	build.BuiltinMap["filter"] = build.BuiltinDescriptor{
		Function: Filter,
		Help:     "Filter a list of files with excludes",
	}
}

func Filter(includes []string, excludes ...string) []string {
	var filtered []string
	for _, include := range includes {
		excluded := false
		for _, pattern := range excludes {
			// DEBUG
			println("include:", include, "pattern:", pattern)
			exclude, err := zglob.Match(pattern, include)
			// DEBUG
			println("exclude:", exclude, "err:", err)
			if exclude && err == nil {
				// DEBUG
				println("excluded!")
				excluded = true
				break
			}
		}
		if !excluded {
			filtered = append(filtered, include)
		}
	}
	return filtered
}
