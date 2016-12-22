package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const (
	FILENAME = "build.yml"
)

var directory = "."
var build Build

func printError(err error, message string) {
	if err != nil {
		println(message)
		os.Exit(2)
	}
}

func main() {
	source, err := ioutil.ReadFile(FILENAME)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &build)
	if err != nil {
		panic(err)
	}
	var targets []string
	if len(os.Args) > 1 {
		targets = os.Args[1:]
	} else {
		if build.Default != "" {
			targets = []string{build.Default}
		} else {
			println("No default target")
			os.Exit(2)
		}
	}
	for _, t := range targets {
		build.Run(t)
	}
}
