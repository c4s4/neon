package build

import (
	"fmt"
	"github.com/fatih/color"
	"neon/util"
	"strings"
	"syscall"
	"unicode/utf8"
	"unsafe"
	"github.com/nsf/termbox-go"
)

// Size of a terminal window
type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

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
	width, _ := termWidth()
	length := width - (4 + utf8.RuneCountInString(text))
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
func printColor(format string, fields ...interface{}) {
	fmt.Fprintf(color.Output, format, fields...)
	fmt.Println()
}

// Print string with arguments in grey
func printGrey(format string, fields ...interface{}) {
	fmt.Printf(format, fields...)
	fmt.Println()
}

// Get terminal width
func termWidth() (int, error) {
	if util.Windows() {
		if err := termbox.Init(); err != nil {
			return 80, err
		}
		width, _ := termbox.Size()
		termbox.Close()
		return width, nil
	} else {
		ws := &winsize{}
		retCode, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
			uintptr(syscall.Stdin),
			uintptr(syscall.TIOCGWINSZ),
			uintptr(unsafe.Pointer(ws)))
		if int(retCode) == -1 {
			return 80, fmt.Errorf("getting terminal width")
		}
		return int(ws.Col), nil
	}
}
