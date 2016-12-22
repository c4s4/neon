package main

import (
	"fmt"
	"os"
)

type Build struct {
	Name       string
	Default    string
	Properties map[string]string
	Targets    map[string]*Target
}

func (b Build) Run(t string) {
	fmt.Printf("# Running target %s\n", t)
	target := b.Targets[t]
	target.Run()
}

func (b Build) Target(t string) *Target {
	if target, ok := b.Targets[t]; ok {
		return target
	} else {
		fmt.Printf("Target %s was not found", t)
		os.Exit(3)
		return nil
	}
}
