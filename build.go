package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Build struct {
	File       string
	Name       string
	Default    string
	Properties map[string]string
	Targets    map[string]*Target
}

func (build Build) Init(file string) {
	build.File = file
}

func (build Build) Dir() string {
	return filepath.Dir(build.File)
}

func (build Build) Run(targets []string) {
	for _, target := range targets {
		fmt.Printf("# Running target %s\n", target)
		target := build.Targets[target]
		target.Run()
	}
}

func (build Build) Target(t string) *Target {
	if target, ok := build.Targets[t]; ok {
		return target
	} else {
		fmt.Printf("Target %s was not found", t)
		os.Exit(3)
		return nil
	}
}
