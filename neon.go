package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

const (
	DEFAULT_BUILD_FILE = "build.yml"
)

func ParseCommandLine() (*string, *bool) {
	buildFile := flag.String("file", DEFAULT_BUILD_FILE, "Build file to run")
	buildHelp := flag.Bool("build", false, "Print build help")
	flag.Parse()
	return buildFile, buildHelp
}

func LoadBuildFile(buildFile *string) *Build {
	var build *Build
	source, err := ioutil.ReadFile(*buildFile)
	StopOnError(err, "Error loading build file '"+*buildFile+"'", 1)
	err = yaml.Unmarshal(source, &build)
	StopOnError(err, "Error parsing build file '"+*buildFile+"'", 2)
	absBuildFile, err := filepath.Abs(*buildFile)
	StopOnError(err, "Error getting build file path", 4)
	build.Init(absBuildFile)
	return build
}

func main() {
	buildFile, buildHelp := ParseCommandLine()
	build := LoadBuildFile(buildFile)
	if *buildHelp {
		build.Help()
	} else {
		build.Run()
	}
}
