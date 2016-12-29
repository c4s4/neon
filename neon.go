package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
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

func LoadBuildFile(file string) (*Build, error) {
	var build *Build
	source, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("loading build file '%s': %v", file, err)
	}
	err = yaml.Unmarshal(source, &build)
	if err != nil {
		return nil, fmt.Errorf("parsing build file '%s': %v", file, err)
	}
	absBuildFile, err := filepath.Abs(file)
	if err != nil {
		return nil, fmt.Errorf("getting build file path: %v", err)
	}
	build.Init(absBuildFile)
	return build, nil
}

func main() {
	file, help := ParseCommandLine()
	build, err := LoadBuildFile(*file)
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
		err = build.Run()
		if err != nil {
			PrintError(err.Error())
			os.Exit(2)
		}
	}
}
