package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

type Build struct {
	File       string
	Name       string
	Default    string
	Properties map[string]string
	Targets    map[string]*Target
}

func (build *Build) Init(file string) {
	build.File = file
	for name, target := range build.Targets {
		target.Init(build, name)
	}
}

func (build *Build) Dir() string {
	return filepath.Dir(build.File)
}

func (build *Build) Run(targets []string) {
	for _, target := range targets {
		target := build.Targets[target]
		target.Run()
	}
}

func (build *Build) Target(t string) *Target {
	if target, ok := build.Targets[t]; ok {
		return target
	} else {
		fmt.Printf("Target %s was not found", t)
		os.Exit(3)
		return nil
	}
}

func (build *Build) ReplaceProperty(expression string) string {
	property := expression[2 : len(expression)-1]
	if value, ok := build.Properties[property]; ok {
		return value
	} else {
		println("Property %s was not found", property)
		os.Exit(3)
		return ""
	}
}

func (build *Build) ReplaceProperties(command string) string {
	r := regexp.MustCompile("#{.*?}")
	replaced := r.ReplaceAllStringFunc(command, build.ReplaceProperty)
	return replaced
}

func (build *Build) Execute(cmd string) string {
	cmd = build.ReplaceProperties(cmd)
	var command *exec.Cmd
	if runtime.GOOS == "windows" {
		command = exec.Command("cmd.exe", "/C", cmd)
	} else {
		command = exec.Command("sh", "-c", cmd)
	}
	command.Dir = filepath.Dir(build.File)
	output, err := command.CombinedOutput()
	result := strings.TrimSpace(string(output))
	StopOnError(err, "Error running command '"+cmd+"': "+result, 5)
	return result
}
