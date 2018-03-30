package build

import (
	"fmt"
	"neon/util"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
)

const (
	defaultTheme = "bee"
)

type colorizer func(a ...interface{}) string

// Theme defined by colors
type Theme struct {
	Title []color.Attribute
	Ok    []color.Attribute
	Error []color.Attribute
}

// Themes is a map of themes by name
var Themes = map[string]*Theme{
	"bee": &Theme{
		Title: []color.Attribute{color.FgYellow},
		Ok:    []color.Attribute{color.FgGreen, color.Bold},
		Error: []color.Attribute{color.FgRed, color.Bold},
	},
	"red": &Theme{
		Title: []color.Attribute{color.FgRed},
		Ok:    []color.Attribute{color.FgRed, color.Bold},
		Error: []color.Attribute{color.FgRed, color.Bold, color.ReverseVideo},
	},
	"green": &Theme{
		Title: []color.Attribute{color.FgGreen},
		Ok:    []color.Attribute{color.FgGreen, color.Bold},
		Error: []color.Attribute{color.FgGreen, color.Bold, color.ReverseVideo},
	},
	"blue": &Theme{
		Title: []color.Attribute{color.FgBlue},
		Ok:    []color.Attribute{color.FgBlue, color.Bold},
		Error: []color.Attribute{color.FgBlue, color.Bold, color.ReverseVideo},
	},
	"fire": &Theme{
		Title: []color.Attribute{color.FgRed},
		Ok:    []color.Attribute{color.FgGreen, color.Bold, color.Underline},
		Error: []color.Attribute{color.FgRed, color.Bold, color.Underline},
	},
	"marine": &Theme{
		Title: []color.Attribute{color.FgBlue},
		Ok:    []color.Attribute{color.FgGreen, color.Bold, color.Underline},
		Error: []color.Attribute{color.FgRed, color.Bold, color.Underline},
	},
	"nature": &Theme{
		Title: []color.Attribute{color.FgGreen},
		Ok:    []color.Attribute{color.FgGreen, color.Bold, color.Underline},
		Error: []color.Attribute{color.FgRed, color.Bold, color.Underline},
	},
	"bold": &Theme{
		Title: []color.Attribute{color.FgYellow, color.Bold},
		Ok:    []color.Attribute{color.FgGreen, color.Underline, color.Bold},
		Error: []color.Attribute{color.FgRed, color.Underline, color.Bold},
	},
	"reverse": &Theme{
		Title: []color.Attribute{color.ReverseVideo},
		Ok:    []color.Attribute{color.ReverseVideo, color.Bold},
		Error: []color.Attribute{color.ReverseVideo, color.Bold},
	},
}

// Colors define a theme
type Colors struct {
	Title []string
	Ok    []string
	Error []string
}

// Attributes values by name
var Attributes = map[string]color.Attribute{
	"Reset":        color.Reset,
	"Bold":         color.Bold,
	"Faint":        color.Faint,
	"Italic":       color.Italic,
	"Underline":    color.Underline,
	"BlinkSlow":    color.BlinkSlow,
	"BlinkRapid":   color.BlinkRapid,
	"ReverseVideo": color.ReverseVideo,
	"Concealed":    color.Concealed,
	"CrossedOut":   color.CrossedOut,
	"FgBlack":      color.FgBlack,
	"FgRed":        color.FgRed,
	"FgGreen":      color.FgGreen,
	"FgYellow":     color.FgYellow,
	"FgBlue":       color.FgBlue,
	"FgMagenta":    color.FgMagenta,
	"FgCyan":       color.FgCyan,
	"FgWhite":      color.FgWhite,
	"FgHiBlack":    color.FgHiBlack,
	"FgHiRed":      color.FgHiRed,
	"FgHiGreen":    color.FgHiGreen,
	"FgHiYellow":   color.FgHiYellow,
	"FgHiBlue":     color.FgHiBlue,
	"FgHiMagenta":  color.FgHiMagenta,
	"FgHiCyan":     color.FgHiCyan,
	"FgHiWhite":    color.FgHiWhite,
	"BgBlack":      color.BgBlack,
	"BgRed":        color.BgRed,
	"BgGreen":      color.BgGreen,
	"BgYellow":     color.BgYellow,
	"BgBlue":       color.BgBlue,
	"BgMagenta":    color.BgMagenta,
	"BgCyan":       color.BgCyan,
	"BgWhite":      color.BgWhite,
	"BgHiBlack":    color.BgHiBlack,
	"BgHiRed":      color.BgHiRed,
	"BgHiGreen":    color.BgHiGreen,
	"BgHiYellow":   color.BgHiYellow,
	"BgHiBlue":     color.BgHiBlue,
	"BgHiMagenta":  color.BgHiMagenta,
	"BgHiCyan":     color.BgHiCyan,
	"BgHiWhite":    color.BgHiWhite,
}

// ParseAttributes parse attributes
func ParseAttributes(colors []string) ([]color.Attribute, error) {
	var attributes []color.Attribute
	for _, col := range colors {
		attribute, ok := Attributes[col]
		if !ok {
			return nil, fmt.Errorf("unknown attribute '%s'", col)
		}
		attributes = append(attributes, attribute)
	}
	return attributes, nil
}

// ParseTheme parses colors and returns a Theme
func ParseTheme(colors *Colors) (*Theme, error) {
	title, err := ParseAttributes(colors.Title)
	if err != nil {
		return nil, err
	}
	ok, err := ParseAttributes(colors.Ok)
	if err != nil {
		return nil, err
	}
	e, err := ParseAttributes(colors.Error)
	if err != nil {
		return nil, err
	}
	return &Theme{Title: title, Ok: ok, Error: e}, nil
}

// Grey is a flag that tells if we print on console without color
var Grey = false

// Color definitions
var colorTitle colorizer
var colorOk colorizer
var colorError colorizer

// apply default theme
func init() {
	ApplyThemeByName(defaultTheme)
}

// ApplyThemeByName applies named theme
func ApplyThemeByName(name string) error {
	theme, ok := Themes[name]
	if !ok {
		return fmt.Errorf("unknown theme '%s'", name)
	}
	ApplyTheme(theme)
	return nil
}

// ApplyTheme applies given theme
func ApplyTheme(theme *Theme) {
	colorTitle = color.New(theme.Title...).SprintFunc()
	colorOk = color.New(theme.Ok...).SprintFunc()
	colorError = color.New(theme.Error...).SprintFunc()
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
