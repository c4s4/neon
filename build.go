package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
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
}

func (build *Build) Init(file string) {
	build.File = file
	build.Dir = filepath.Dir(build.File)
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

func (build *Build) ReplaceProperty(expression string) string {
	property := expression[2 : len(expression)-1]
	if value, ok := build.Properties[property]; ok {
		return value
	} else {
		StopWithError(fmt.Sprintf("Property '%s' was not found", property), 7)
		return ""
	}
}

func (build *Build) ReplaceProperties(command string) string {
	r := regexp.MustCompile("#{.*?}")
	replaced := r.ReplaceAllStringFunc(command, build.ReplaceProperty)
	return replaced
}

func (build *Build) Execute(cmd string) {
	cmd = build.ReplaceProperties(cmd)
	var command *exec.Cmd
	if runtime.GOOS == "windows" {
		command = exec.Command("cmd.exe", "/C", cmd)
	} else {
		command = exec.Command("sh", "-c", cmd)
	}
	command.Dir = build.Dir
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	StopOnError(err, "Error running command '"+cmd+"'", 5)
}

func (build *Build) Help() {
	if build.Doc != "" {
		fmt.Println(build.Doc)
		fmt.Println()
	}
	length := 0
	targets := []string{}
	for name, _ := range build.Targets {
		if utf8.RuneCountInString(name) > length {
			length = len(name)
		}
		targets = append(targets, name)
	}
	sort.Strings(targets)
	for _, target := range targets {
		PrintTargetHelp(target, build.Target(target).Doc, length)
	}
}
