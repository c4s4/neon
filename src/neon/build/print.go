package build

import (
	"fmt"
	"github.com/fatih/color"
	"neon/util"
	"strings"
	"unicode/utf8"
)

const (
	DEFAULT_WIDTH = 80
)

// Flag that tells if we print on console without color
var Grey = false

// Color definitions
var red = color.New(color.FgRed, color.Bold).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var green = color.New(color.FgGreen, color.Bold).SprintFunc()

// Print a message
func Message(text string, args ...interface{}) {
	printGrey(text, args...)
}

// Print a title
func Title(text string) {
	length := util.TerminalWidth() - (4 + utf8.RuneCountInString(text))
	if length < 2 {
		length = 2
	}
	message := fmt.Sprintf("%s %s --", strings.Repeat("-", length), text)
	if Grey {
		printGrey(message)
	} else {
		printColor(yellow(message))
	}
}

// Print OK
func PrintOk() {
	if Grey {
		printGrey("OK")
	} else {
		printColor(green("OK"))
	}
}

// Print ERROR
func PrintError(text string) {
	if Grey {
		printGrey("ERROR %s", text)
	} else {
		printColor("%s %s", red("ERROR"), text)
	}
}

// Print string with arguments in given color
func printColor(text string, fields ...interface{}) {
	if len(fields) > 0 {
		fmt.Fprintf(color.Output, text, fields...)
		fmt.Println()
	} else {
		fmt.Println(text)
	}
}

// Print string with arguments in grey
func printGrey(text string, fields ...interface{}) {
	if len(fields) > 0 {
		fmt.Printf(text, fields...)
		fmt.Println()
	} else {
		fmt.Println(text)
	}
}
