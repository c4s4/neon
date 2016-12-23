package main

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
	"unicode/utf8"
)

var red = color.New(color.FgRed, color.Bold).SprintFunc()
var cyan = color.New(color.FgYellow).SprintFunc()

func PrintTarget(message string) {
	color.Yellow(message)
}

func PrintTargetHelp(name, doc string, length int) {
	if doc != "" {
		fmt.Fprintf(color.Output, "%s%s %s\n", cyan(name),
			strings.Repeat(" ", length-utf8.RuneCountInString(name)), doc)
	} else {
		fmt.Fprintf(color.Output, "%s\n", cyan(name))
	}
}

func PrintError(message string) {
	fmt.Fprintf(color.Output, "%s: %s\n", red("ERROR"), message)
}

func PrintOK() {
	color.New(color.FgGreen).Add(color.Bold).Println("OK")
}
