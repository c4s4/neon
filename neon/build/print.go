package build

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/c4s4/neon/neon/util"
	"github.com/fatih/color"
)

type colorizer func(a ...interface{}) string

// Gray is a flag that tells if we print on console without color
var Gray = false

// Color definitions
var colorTitle colorizer
var colorOk colorizer
var colorError colorizer

// Message prints a message on console:
// - text: text to print (that might embed fields to print, such as "%s")
func Message(text string) {
	printGray(text)
}

// MessageArgs prints a message on console:
// - text: text to print (that might embed fields to print, such as "%s")
// - args: arguments for the text to print
func MessageArgs(text string, args ...interface{}) {
	printGrayArgs(text, args...)
}

// Info prints an information message on console:
// - text: text to print (that might embed fields to print, such as "%s")
func Info(text string) {
	if Gray {
		printGray(text)
	} else {
		printColor(colorTitle(text))
	}
}

// InfoArgs prints an information message on console:
// - text: text to print (that might embed fields to print, such as "%s")
// - args: arguments for the text to print
func InfoArgs(text string, args ...interface{}) {
	if Gray {
		printGrayArgs(text, args...)
	} else {
		printColorArgs(colorTitle(text), args...)
	}
}

// Title prints a title on the console
// - text: text of the title to print
func Title(text string) {
	length := util.TerminalWidth() - (4 + utf8.RuneCountInString(text))
	if length < 2 {
		length = 2
	}
	message := fmt.Sprintf("%s %s --", strings.Repeat("-", length), text)
	if Gray {
		printGray(message)
	} else {
		printColor(colorTitle(message))
	}
}

// PrintOk prints a green OK on the console
func PrintOk() {
	if Gray {
		printGray("OK")
	} else {
		printColor(colorOk("OK"))
	}
}

// PrintError prints a red ERROR on the console followed with an explanatory
// text
// - text: the explanatory text to print
func PrintError(text string) {
	if Gray {
		printGrayArgs("ERROR %s", text)
	} else {
		printColorArgs("%s %s", colorError("ERROR"), text)
	}
}

// PrintColor prints a string in given color
// - text: the text to print
func printColor(text string) {
	fmt.Println(text)
}

// PrintColor prints a string with arguments in given color
// - text: the text to print
// - args: the arguments for the text to print
func printColorArgs(text string, args ...interface{}) {
	fmt.Fprintf(color.Output, text, args...)
	fmt.Println()
}

// PrintGrey prints a string in gray
// - text: the text to print
func printGray(text string) {
	fmt.Println(text)
}

// PrintGreyArgs prints a string with arguments in gray
// - text: the text to print
// - args: the arguments for the text to print
func printGrayArgs(text string, fields ...interface{}) {
	fmt.Printf(text, fields...)
	fmt.Println()
}
