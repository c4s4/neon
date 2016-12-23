package main

import (
	"fmt"
	"github.com/fatih/color"
)

var red = color.New(color.FgRed, color.Bold).SprintFunc()

func PrintTarget(message string) {
	color.Yellow(message)
}

func PrintError(message string) {
	fmt.Fprintf(color.Output, "%s: %s\n", red("ERROR"), message)
}

func PrintOK() {
	color.New(color.FgGreen).Add(color.Bold).Println("OK")
}
