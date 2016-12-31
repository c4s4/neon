package main

import (
	"flag"
	"os"
)

const (
	DEFAULT_BUILD_FILE = "build.yml"
)

func ParseCommandLine() (*string, *bool, []string) {
	file := flag.String("file", DEFAULT_BUILD_FILE, "Build file to run")
	help := flag.Bool("build", false, "Print build help")
	flag.Parse()
	targets := flag.Args()
	return file, help, targets
}

func main() {
	file, help, targets := ParseCommandLine()
	build, err := NewBuild(*file)
	if err != nil {
		PrintError(err.Error())
		os.Exit(1)
	}
	if *help {
		err = build.Help()
		if err != nil {
			PrintError(err.Error())
			os.Exit(2)
		}
	} else {
		err = build.Run(targets)
		PrintOK()
		if err != nil {
			PrintError(err.Error())
			os.Exit(2)
		}
	}
}
