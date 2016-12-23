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

func ParseCommandLine() *string {
	buildFile := flag.String("file", DEFAULT_BUILD_FILE, "build file to run")
	flag.Parse()
	return buildFile
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

func ParseTargets(build *Build) []string {
	targets := flag.Args()
	if len(targets) == 0 {
		if build.Default != "" {
			targets = []string{build.Default}
		} else {
			StopWithError("No default target", 3)
		}
	}
	return targets
}

func main() {
	buildFile := ParseCommandLine()
	build := LoadBuildFile(buildFile)
	targets := ParseTargets(build)
	build.Run(targets)
	PrintOK()
}
