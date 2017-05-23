package util

import (
	"fmt"
	"github.com/fatih/color"
)

// Color definitions
var Red = color.New(color.FgRed, color.Bold).SprintFunc()
var Yellow = color.New(color.FgYellow).SprintFunc()
var Green = color.New(color.FgGreen, color.Bold).SprintFunc()

// Print string with arguments in given color
func PrintColor(format string, fields ...interface{}) {
	fmt.Fprintf(color.Output, format, fields...)
	fmt.Println()
}
