package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	DEFAULT_BUILD_FILE = "build.yml"
)

func StopOnError(err error, message string, code int) {
	if err != nil {
		println(message)
		os.Exit(code)
	}
}

func StopWithError(message string, code int) {
	println(message)
	os.Exit(code)
}

func main() {
	var build *Build
	// parse command line
	buildFile := flag.String("file", DEFAULT_BUILD_FILE, "build file to run")
	flag.Parse()
	source, err := ioutil.ReadFile(*buildFile)
	StopOnError(err, "Error loading build file '"+*buildFile+"'", 1)
	err = yaml.Unmarshal(source, &build)
	StopOnError(err, "Error parsing build file '"+*buildFile+"'", 2)
	// initialize build
	absBuildFile, err := filepath.Abs(*buildFile)
	StopOnError(err, "Error getting build file path", 4)
	build.Init(absBuildFile)
	// parse targets on command line
	targets := flag.Args()
	if len(targets) == 0 {
		if build.Default != "" {
			targets = []string{build.Default}
		} else {
			StopWithError("No default target", 3)
		}
	}
	// run build
	build.Run(targets)
}
