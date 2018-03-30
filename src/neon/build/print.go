package build

import (
	"fmt"
	"github.com/fatih/color"
	"neon/util"
	"strings"
	"unicode/utf8"
)

const (
	defaultTheme = "bee"
)

type colorizer func(a ...interface{}) string

// Themes is a map of themes by name
var Themes = map[string]map[string][]color.Attribute{
	"bee": {
		"title": {color.FgYellow},
		"ok":    {color.FgGreen, color.Bold},
		"error": {color.FgRed, color.Bold},
	},
	"marine": {
		"title": {color.FgBlue},
		"ok":    {color.FgGreen, color.BgBlack, color.Bold},
		"error": {color.FgRed, color.BgBlack, color.Bold},
	},
	"bold": {
		"title": {color.FgYellow, color.Bold},
		"ok":    {color.FgGreen, color.Underline, color.Bold},
		"error": {color.FgRed, color.Underline, color.Bold},
	},
}

// Grey is a flag that tells if we print on console without color
var Grey = false

// Color definitions
var colorTitle colorizer
var colorOk colorizer
var colorError colorizer

// apply default theme
func init() {
	ApplyTheme(defaultTheme)
}

// ApplyTheme applies named theme
func ApplyTheme(theme string) error {
	if _, ok := Themes[theme]; !ok {
		return fmt.Errorf("unknown theme '%s'", theme)
	}
	colorTitle = color.New(Themes[theme]["title"]...).SprintFunc()
	colorOk = color.New(Themes[theme]["ok"]...).SprintFunc()
	colorError = color.New(Themes[theme]["error"]...).SprintFunc()
	return nil
}

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
		printColor(colorTitle(message))
	}
}

// PrintOk prints a green OK on the console
func PrintOk() {
	if Grey {
		printGrey("OK")
	} else {
		printColor(colorOk("OK"))
	}
}

// PrintError prints a red ERROR on the console followed with an explanatory
// text
// - text: the explanatory text to print
func PrintError(text string) {
	if Grey {
		printGrey("ERROR %s", text)
	} else {
		printColor("%s %s", colorError("ERROR"), text)
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
