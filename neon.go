package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	FILENAME = "build.yml"
)

var directory = "."

func printError(err error, message string) {
	if err != nil {
		println(message)
		os.Exit(2)
	}
}

func execute(cmd string) string {
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

type Target []Step

func (t Target) Run() {
	for _, s := range t {
		s.Run()
	}
}

type Build struct {
	Name    string
	Default string
	Targets map[string]Target
}

func (b Build) Run(t string) {
	fmt.Printf("# Running target %s\n", t)
	target := b.Targets[t]
	target.Run()
}

func main() {
	var build Build
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
