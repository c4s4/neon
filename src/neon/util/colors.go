package util

import (
	"fmt"
	"github.com/fatih/color"
)

var Red = color.New(color.FgRed, color.Bold).SprintFunc()
var Yellow = color.New(color.FgYellow).SprintFunc()
var Green = color.New(color.FgGreen, color.Bold).SprintFunc()

func PrintColor(format string, fields ...interface{}) {
	fmt.Fprintf(color.Output, format, fields...)
	fmt.Println()
}
