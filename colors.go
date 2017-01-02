package main

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
	"unicode/utf8"
)

var red = color.New(color.FgRed, color.Bold).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var green = color.New(color.FgGreen, color.Bold)

func PrintTarget(message string) {
	color.Yellow(message)
}

func PrintTargetHelp(name, doc string, depends []string, length int) {
	deps := ""
	if len(depends) > 0 {
		deps = "[" + strings.Join(depends, ", ") + "]"
	}
	if doc != "" {
		deps = " " + deps
	}
	fmt.Fprintf(color.Output, "%s%s %s%s\n", yellow(name),
		strings.Repeat(" ", length-utf8.RuneCountInString(name)), doc, deps)
}

func PrintError(message string) {
	fmt.Fprintf(color.Output, "%s %s\n", red("ERROR"), message)
}

func PrintOK() {
	green.Println("OK")
}
