package main

import (
	"errors"
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

func (build *Build) ParseTargets() ([]string, error) {
	targets := flag.Args()
	if len(targets) == 0 {
		if build.Default != "" {
			targets = []string{build.Default}
		} else {
			return nil, errors.New("No default target")
		}
	}
	return targets, nil
}

func (build *Build) Run() error {
	targets, err := build.ParseTargets()
	if err != nil {
		return err
	}
	for _, target := range targets {
		target, err := build.Target(target)
		if err != nil {
			return err
		}
		target.Run()
	}
	PrintOK()
	return nil
}

func (build *Build) Target(name string) (*Target, error) {
	if target, ok := build.Targets[name]; ok {
		return target, nil
	} else {
		return nil, fmt.Errorf("Target '%s' was not found", name)
	}
}

func (build *Build) Help() error {
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
			value, err := build.Context.GetProperty(name)
			if err != nil {
				return err
			}
			valueStr, err := Serialize(value)
			return fmt.Errorf("Error serializing property '" + name + "' value")
			PrintTargetHelp(name, valueStr, []string{}, length)
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
			target, err := build.Target(name)
			if err != nil {
				return err
			}
			PrintTargetHelp(name, target.Doc, target.Depends, length)
		}
	}
	return nil
}
