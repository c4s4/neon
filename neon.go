package main

import (
	"flag"
	"os"
)

const (
	DEFAULT_BUILD_FILE = "build.yml"
)

func ParseCommandLine() (*string, *bool) {
	file := flag.String("file", DEFAULT_BUILD_FILE, "Build file to run")
	help := flag.Bool("build", false, "Print build help")
	flag.Parse()
	return file, help
}

func main() {
	file, help := ParseCommandLine()
	build, err := NewBuild(*file)
	if err != nil {
		PrintError(err.Error())
		os.Exit(1)
	}
	if *help {
		err = build.Help()
	} else {
		err = build.Run()
	}
	if err != nil {
		PrintError(err.Error())
		os.Exit(2)
	}
}
