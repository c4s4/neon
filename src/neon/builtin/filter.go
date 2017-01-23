package builtin

import (
	"neon/build"
)

func init() {
	build.BuiltinMap["filter"] = build.BuiltinDescriptor{
		Function: Filter,
		Help:     "Filter a list of files with excludes",
	}
}

func Filter(includes, excludes []string) []string {
	var filtered []string
	for _, include := range includes {
		excluded := false
		for _, exclude := range excludes {
			if include == exclude {
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
