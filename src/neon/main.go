package main

import (
	"flag"
	"fmt"
	"neon/build"
	"neon/util"
	"os"
	"path/filepath"
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

func FindBuildFile(name string) (string, error) {
	absolute, err := filepath.Abs(name)
	if err != nil {
		return "", fmt.Errorf("getting build file path: %v", err)
	}
	file := filepath.Base(absolute)
	dir := filepath.Dir(absolute)
	for {
		path := filepath.Join(dir, file)
		if util.FileExists(path) {
			return path, nil
		} else {
			parent := filepath.Dir(dir)
			if parent == dir {
				return "", fmt.Errorf("build file not found")
			}
			dir = parent
		}
	}
}

func main() {
	file, help, debug, targets := ParseCommandLine()
	path, err := FindBuildFile(*file)
	if err != nil {
		util.PrintError(err.Error())
		os.Exit(1)
	}
	build, err := build.NewBuild(path, *debug)
	if err != nil {
		util.PrintError(err.Error())
		os.Exit(2)
	}
	if *help {
		err = build.Help()
		if err != nil {
			util.PrintError(err.Error())
			os.Exit(3)
		}
	} else {
		err = build.Run(targets)
		if err == nil {
			util.PrintOK()
		} else {
			util.PrintError(err.Error())
			os.Exit(4)
		}
	}
}
