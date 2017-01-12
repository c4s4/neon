package main

import (
	"flag"
	"neon/build"
	"neon/util"
	"os"
)

const (
	DEFAULT_BUILD_FILE = "build.yml"
)

func ParseCommandLine() (*string, *bool, *bool, []string) {
	file := flag.String("file", DEFAULT_BUILD_FILE, "Build file to run")
	help := flag.Bool("build", false, "Print build help")
	debug := flag.Bool("debug", false, "Output debugging information")
	flag.Parse()
	targets := flag.Args()
	return file, help, debug, targets
}

func main() {
	file, help, debug, targets := ParseCommandLine()
	build, err := build.NewBuild(*file, *debug)
	if err != nil {
		util.PrintError(err.Error())
		os.Exit(1)
	}
	if *help {
		err = build.Help()
		if err != nil {
			util.PrintError(err.Error())
			os.Exit(2)
		}
	} else {
		err = build.Run(targets)
		if err == nil {
			util.PrintOK()
		} else {
			util.PrintError(err.Error())
			os.Exit(2)
		}
	}
}
