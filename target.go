package main

import (
	"fmt"
)

type Target struct {
	Name    string
	Build   *Build
	Depends []string
	Steps   []string
}

func (target *Target) Init(build *Build, name string) {
	target.Build = build
	target.Name = name
}

func (target *Target) Run() {
	for _, depend := range target.Depends {
		dependency := target.Build.Target(depend)
		dependency.Run()
	}
	fmt.Printf("# Running target %s\n", target.Name)
	for _, step := range target.Steps {
		output := target.Build.Execute(step)
		fmt.Println(output)
	}
}
