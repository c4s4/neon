package builtin

import (
	"github.com/c4s4/neon/neon/build"
	colorize "github.com/fatih/color"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "color",
		Func: color,
		Help: `Colorize string.

Arguments:

- The name of the color (black, red, green, yellow, blue, magenta, cyan or white).
- The string to colorize.

Color capitalization determines the intensity of the color:

- Use standard color names (such as black) for normal colors.
- Use capitalized color names (such as Black) for bold colors.
- Use uppercase color names (such as BLACK) for bold background colors.
- Use uncapitalized uppercase color names (such as bLACK) for background colors.

Returns:

- The colorized string.

Examples:

    # green message
    color("green", "OK")
    # returns: string "OK" colorized in green`,
	})
}

func color(color, text string) string {
	switch color {
	case "black":
		return colorize.New(colorize.FgBlack).Sprint(text)
	case "red":
		return colorize.New(colorize.FgRed).Sprint(text)
	case "green":
		return colorize.New(colorize.FgGreen).Sprint(text)
	case "yellow":
		return colorize.New(colorize.FgYellow).Sprint(text)
	case "blue":
		return colorize.New(colorize.FgBlue).Sprint(text)
	case "magenta":
		return colorize.New(colorize.FgMagenta).Sprint(text)
	case "cyan":
		return colorize.New(colorize.FgCyan).Sprint(text)
	case "white":
		return colorize.New(colorize.FgWhite).Sprint(text)
	case "Black":
		return colorize.New(colorize.FgBlack, colorize.Bold).Sprint(text)
	case "Red":
		return colorize.New(colorize.FgRed, colorize.Bold).Sprint(text)
	case "Green":
		return colorize.New(colorize.FgGreen, colorize.Bold).Sprint(text)
	case "Yellow":
		return colorize.New(colorize.FgYellow, colorize.Bold).Sprint(text)
	case "Blue":
		return colorize.New(colorize.FgBlue, colorize.Bold).Sprint(text)
	case "Magenta":
		return colorize.New(colorize.FgMagenta, colorize.Bold).Sprint(text)
	case "Cyan":
		return colorize.New(colorize.FgCyan, colorize.Bold).Sprint(text)
	case "White":
		return colorize.New(colorize.FgWhite, colorize.Bold).Sprint(text)
	case "bLACK":
		return colorize.New(colorize.BgBlack).Sprint(text)
	case "rED":
		return colorize.New(colorize.BgRed).Sprint(text)
	case "gREEN":
		return colorize.New(colorize.BgGreen).Sprint(text)
	case "yELLOW":
		return colorize.New(colorize.BgYellow).Sprint(text)
	case "bLUE":
		return colorize.New(colorize.BgBlue).Sprint(text)
	case "mAGENTA":
		return colorize.New(colorize.BgMagenta).Sprint(text)
	case "cYAN":
		return colorize.New(colorize.BgCyan).Sprint(text)
	case "wHITE":
		return colorize.New(colorize.BgWhite).Sprint(text)
	case "BLACK":
		return colorize.New(colorize.BgBlack, colorize.Bold).Sprint(text)
	case "RED":
		return colorize.New(colorize.BgRed, colorize.Bold).Sprint(text)
	case "GREEN":
		return colorize.New(colorize.BgGreen, colorize.Bold).Sprint(text)
	case "YELLOW":
		return colorize.New(colorize.BgYellow, colorize.Bold).Sprint(text)
	case "BLUE":
		return colorize.New(colorize.BgBlue, colorize.Bold).Sprint(text)
	case "MAGENTA":
		return colorize.New(colorize.BgMagenta, colorize.Bold).Sprint(text)
	case "CYAN":
		return colorize.New(colorize.BgCyan, colorize.Bold).Sprint(text)
	case "WHITE":
		return colorize.New(colorize.BgWhite, colorize.Bold).Sprint(text)
	default:
		return text
	}
}
