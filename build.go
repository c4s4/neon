package main

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

func NewBuild(file string) (*Build, error) {
	source, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("loading build file '%s': %v", file, err)
	}
	var structure map[string]interface{}
	err = yaml.Unmarshal(source, &structure)
	if err != nil {
		return nil, fmt.Errorf("parsing build file '%s': must be a YAML map", file)
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return nil, fmt.Errorf("getting build file path: %v", err)
	}
	build := &Build{}
	err = build.Init(path, &structure)
	return build, err
}

func (build *Build) Init(file string, structure *map[string]interface{}) error {
	str, err := getString(structure, "name")
	if err != nil {
		return err
	}
	build.Name = str
	str, err = getString(structure, "default")
	if err != nil {
		return err
	}
	build.Default = str
	str, err = getString(structure, "doc")
	if err != nil {
		return err
	}
	build.Doc = str
	build.File = file
	build.Dir = filepath.Dir(build.File)
	//build.Context = NewContext(build)
	//for name, target := range build.Targets {
	//	target.Init(build, name)
	//}
	return nil
}

func getString(dict *map[string]interface{}, name string) (string, error) {
	value, ok := (*dict)[name]
	if !ok {
		return "", nil
	}
	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("field '%s' must be a string", name)
	}
	return str, nil
}

func (build *Build) ParseTargets() ([]string, error) {
	targets := flag.Args()
	if len(targets) == 0 {
		if build.Default != "" {
			targets = []string{build.Default}
		} else {
			return nil, errors.New("no default target")
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
		return nil, fmt.Errorf("target '%s' was not found", name)
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
				return fmt.Errorf("getting property '%s': %v", name, err)
			}
			valueStr, err := Serialize(value)
			if err != nil {
				return fmt.Errorf("formatting property '%s': %v", name, err)
			}
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
