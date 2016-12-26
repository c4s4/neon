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
	Properties map[string]interface{}
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
	newLine := false
	// print build documentation
	if build.Doc != "" {
		fmt.Println(build.Doc)
		newLine = true
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
	if len(properties) > 0 {
		if newLine {
			fmt.Println()
		}
		fmt.Println("Properties:")
		for _, name := range properties {
			value := build.Context.GetProperty(name)
			valueString := fmt.Sprintf("%v", value)
			PrintTargetHelp(name, valueString, []string{}, length)
		}
		newLine = true
	}
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
	if len(targets) > 0 {
		if newLine {
			fmt.Println()
		}
		fmt.Println("Targets:")
		for _, name := range targets {
			target := build.Target(name)
			PrintTargetHelp(name, target.Doc, target.Depends, length)
		}
	}
}
