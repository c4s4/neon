package build

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"

	"github.com/c4s4/neon/neon/util"
)

// ErrorOutput is the writer for error messages (standard error by default)
var ErrorOutput io.Writer = os.Stderr

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

// PrintError prints a red ERROR on standard error followed with an
// explanatory text. Color is used only if standard error is a terminal.
// - text: the explanatory text to print
func PrintError(text string) {
	if Gray || !isTerminal(ErrorOutput) {
		_, _ = fmt.Fprintf(ErrorOutput, "ERROR %s\n", text)
	} else {
		writer := ErrorOutput
		if writer == os.Stderr {
			// handles colors on Windows
			writer = color.Error
		}
		_, _ = fmt.Fprintf(writer, "%s %s\n", colorError("ERROR"), text)
	}
}

// isTerminal tells if given writer is a terminal
func isTerminal(writer io.Writer) bool {
	file, ok := writer.(*os.File)
	if !ok {
		return false
	}
	return isatty.IsTerminal(file.Fd()) || isatty.IsCygwinTerminal(file.Fd())
}

// printColor prints a string, already colorized
// - text: the text to print
func printColor(text string) {
	fmt.Println(text)
}

// printColorArgs prints a string with arguments, already colorized
// - text: the text to print
// - args: the arguments for the text to print
func printColorArgs(text string, args ...interface{}) {
	_, _ = fmt.Fprintf(color.Output, text, args...)
	fmt.Println()
}

// printGray prints a string without color
// - text: the text to print
func printGray(text string) {
	fmt.Println(text)
}

// printGrayArgs prints a string with arguments without color
// - text: the text to print
// - args: the arguments for the text to print
func printGrayArgs(text string, fields ...interface{}) {
	fmt.Printf(text, fields...)
	fmt.Println()
}
