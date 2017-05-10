package build

import (
	"fmt"
	"neon/util"
	"strings"
	"unicode/utf8"
)

func Info(message string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(message, args...))
}

func PrintColorLine(name, doc string, depends []string, length int) {
	deps := ""
	if len(depends) > 0 {
		deps = "[" + strings.Join(depends, ", ") + "]"
	}
	if doc != "" {
		deps = " " + deps
	}
	util.PrintColor("%s%s %s%s", util.Yellow(name),
		strings.Repeat(" ", length-utf8.RuneCountInString(name)), doc, deps)
}
