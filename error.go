package main

import (
	"fmt"
	"os"
)

func StopOnError(err error, message string, code int) {
	if err != nil {
		PrintError(fmt.Sprintf("%s (%s)", message, err.Error()))
		os.Exit(code)
	}
}

func StopWithError(message string, code int) {
	PrintError(message)
	os.Exit(code)
}
