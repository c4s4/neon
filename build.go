package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"sort"
	"unicode/utf8"
)

type Build struct {
	File       string
	Dir        string
	Name       string
	Default    string
	Doc        string
	Properties map[string]string
	Targets    map[string]*Target
	Context    *Context
}

func (build *Build) Init(file string) {
	build.File = file
	build.Dir = filepath.Dir(build.File)
	build.Context = NewContext(build)
	for name, target := range build.Targets {
		target.Init(build, name)
	}
}

func (build *Build) ParseTargets() []string {
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

func (build *Build) Run() {
	targets := build.ParseTargets()
	for _, target := range targets {
		target := build.Target(target)
		target.Run()
	}
	PrintOK()
}

func (build *Build) Target(name string) *Target {
	if target, ok := build.Targets[name]; ok {
		return target
	} else {
		StopWithError(fmt.Sprintf("Target '%s' was not found", name), 6)
		return nil
	}
}

func (build *Build) Help() {
	// print build documentation
	if build.Doc != "" {
		fmt.Println(build.Doc)
		fmt.Println()
	}
	// print build properties
	length := 0
	properties := []string{}
	for name, _ := range build.Properties {
		if utf8.RuneCountInString(name) > length {
			length = utf8.RuneCountInString(name)
		}
		properties = append(properties, name)
	}
	sort.Strings(properties)
	fmt.Println("Properties:")
	for _, name := range properties {
		value := build.Context.GetProperty(name).String()
		PrintTargetHelp(name, value, []string{}, length)
	}
	fmt.Println()
	// print targets documentation
	length = 0
	targets := []string{}
	for name, _ := range build.Targets {
		if utf8.RuneCountInString(name) > length {
			length = utf8.RuneCountInString(name)
		}
		targets = append(targets, name)
	}
	sort.Strings(targets)
	fmt.Println("Targets:")
	for _, name := range targets {
		target := build.Target(name)
		PrintTargetHelp(name, target.Doc, target.Depends, length)
	}
}
