package task

import (
	"fmt"
	"reflect"

	"github.com/c4s4/neon/neon/build"
	"github.com/fatih/color"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "print",
		Func: print,
		Args: reflect.TypeOf(printArgs{}),
		Help: `Print a message on the console.

Arguments:

- print: text to print (string).
- color: text color (string).

Possible colors are black, red, green, yellow, blue, magenta, cyan and white.

Examples:

    # say hello
    - print: 'Hello World!'
    # say hello in blue
    - print: 'Hello World!'
      color: blue`,
	})
}

type printArgs struct {
	Print string
	Color string `neon:"optional"`
}

// Colors is the color mapping
var Colors = map[string]color.Attribute{
	"black":   color.FgBlack,
	"red":     color.FgRed,
	"green":   color.FgGreen,
	"yellow":  color.FgYellow,
	"blue":    color.FgBlue,
	"magenta": color.FgMagenta,
	"cyan":    color.FgCyan,
	"white":   color.FgWhite,
}

func print(context *build.Context, args interface{}) error {
	params := args.(printArgs)
	if params.Color != "" {
		colorPrint, ok := Colors[params.Color]
		if !ok {
			return fmt.Errorf("color %s not found", params.Color)
		}
		fmt.Fprintf(color.Output, color.New(colorPrint).SprintFunc()(params.Print))
		fmt.Println()
	} else {
		context.Message(params.Print)
	}
	return nil
}
