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

// Message prints a message on console:
// - text: text to print (that might embed fields to print, such as "%s")
// - args: arguments for the text to print
func Message(text string, args ...interface{}) {
	printGrey(text, args...)
}

// Title prints a title on the console
// - text: text of the title to print
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

// PrintOK prints a green OK on the console
func PrintOk() {
	if Grey {
		printGrey("OK")
	} else {
		printColor(green("OK"))
	}
}

// PrintError prints a red ERROR on the console followed with an explanatory
// text
// - text: the explanatory text to print
func PrintError(text string) {
	if Grey {
		printGrey("ERROR %s", text)
	} else {
		printColor("%s %s", red("ERROR"), text)
	}
}

// PrintColor prints a string with arguments in given color
// - text: the text to print
// - args: the arguments for the text to print
func printColor(text string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Fprintf(color.Output, text, args...)
		fmt.Println()
	} else {
		fmt.Println(text)
	}
}

// PrintGrey prints a string with arguments in grey
// - text: the text to print
// - args: the arguments for the text to print
func printGrey(text string, fields ...interface{}) {
	if len(fields) > 0 {
		fmt.Printf(text, fields...)
		fmt.Println()
	} else {
		fmt.Println(text)
	}
}
