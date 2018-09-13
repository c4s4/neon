package build

import (
	"fmt"
	"github.com/fatih/color"
)

const (
	defaultTheme = "bee"
)

// Colors define a theme
type Colors struct {
	Title []string
	Ok    []string
	Error []string
}

// Theme defined by colors
type Theme struct {
	Title []color.Attribute
	Ok    []color.Attribute
	Error []color.Attribute
}

// Themes is a map of themes by name
var Themes = map[string]*Theme{
	"bee": {
		Title: []color.Attribute{color.FgYellow},
		Ok:    []color.Attribute{color.FgGreen, color.Bold},
		Error: []color.Attribute{color.FgRed, color.Bold},
	},
	"rgb": {
		Title: []color.Attribute{color.FgBlue},
		Ok:    []color.Attribute{color.FgGreen, color.Bold},
		Error: []color.Attribute{color.FgRed, color.Bold},
	},
	"red": {
		Title: []color.Attribute{color.FgRed},
		Ok:    []color.Attribute{color.FgRed, color.Bold},
		Error: []color.Attribute{color.FgRed, color.Bold, color.ReverseVideo},
	},
	"green": {
		Title: []color.Attribute{color.FgGreen},
		Ok:    []color.Attribute{color.FgGreen, color.Bold},
		Error: []color.Attribute{color.FgGreen, color.Bold, color.ReverseVideo},
	},
	"blue": {
		Title: []color.Attribute{color.FgBlue},
		Ok:    []color.Attribute{color.FgBlue, color.Bold},
		Error: []color.Attribute{color.FgBlue, color.Bold, color.ReverseVideo},
	},
	"yellow": {
		Title: []color.Attribute{color.FgYellow},
		Ok:    []color.Attribute{color.FgYellow, color.Bold},
		Error: []color.Attribute{color.FgYellow, color.Bold, color.ReverseVideo},
	},
	"magenta": {
		Title: []color.Attribute{color.FgMagenta},
		Ok:    []color.Attribute{color.FgMagenta, color.Bold},
		Error: []color.Attribute{color.FgMagenta, color.Bold, color.ReverseVideo},
	},
	"cyan": {
		Title: []color.Attribute{color.FgCyan},
		Ok:    []color.Attribute{color.FgCyan, color.Bold},
		Error: []color.Attribute{color.FgCyan, color.Bold, color.ReverseVideo},
	},
	"fire": {
		Title: []color.Attribute{color.FgRed},
		Ok:    []color.Attribute{color.FgGreen, color.Bold, color.Underline},
		Error: []color.Attribute{color.FgRed, color.Bold, color.Underline},
	},
	"marine": {
		Title: []color.Attribute{color.FgBlue},
		Ok:    []color.Attribute{color.FgGreen, color.Bold, color.Underline},
		Error: []color.Attribute{color.FgRed, color.Bold, color.Underline},
	},
	"nature": {
		Title: []color.Attribute{color.FgGreen},
		Ok:    []color.Attribute{color.FgGreen, color.Bold, color.Underline},
		Error: []color.Attribute{color.FgRed, color.Bold, color.Underline},
	},
	"bold": {
		Title: []color.Attribute{color.FgYellow, color.Bold},
		Ok:    []color.Attribute{color.FgGreen, color.Underline, color.Bold},
		Error: []color.Attribute{color.FgRed, color.Underline, color.Bold},
	},
	"reverse": {
		Title: []color.Attribute{color.ReverseVideo},
		Ok:    []color.Attribute{color.ReverseVideo, color.Bold},
		Error: []color.Attribute{color.ReverseVideo, color.Bold},
	},
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
