package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func replaceProperty(expression string) string {
	property := expression[2 : len(expression)-1]
	if value, ok := build.Properties[property]; ok {
		return value
	} else {
		println("Property %s was not found", property)
		os.Exit(3)
		return ""
	}
}

func replaceProperties(cmd string) string {
	r := regexp.MustCompile("#{.*?}")
	replaced := r.ReplaceAllStringFunc(cmd, replaceProperty)
	return replaced
}

func execute(cmd string) string {
	cmd = replaceProperties(cmd)
	var command *exec.Cmd
	if runtime.GOOS == "windows" {
		command = exec.Command("cmd.exe", "/C", cmd)
	} else {
		command = exec.Command("sh", "-c", cmd)
	}
	command.Dir = directory
	output, err := command.CombinedOutput()
	result := strings.TrimSpace(string(output))
	printError(err, "Error running command '"+cmd+"': "+result)
	return result
}

type Step string

func (s Step) Run() {
	output := execute(string(s))
	fmt.Println(output)
}
