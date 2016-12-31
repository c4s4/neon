package main

import (
	"fmt"
)

type Target struct {
	Build   *Build
	Name    string
	Doc     string
	Depends []string
	Steps   []*Step
}

func NewTarget(build *Build, name string, object Object) (*Target, error) {
	target := &Target{
		Build: build,
		Name:  name,
	}
	err := object.CheckFields([]string{"doc", "depends", "steps"})
	if err != nil {
		return nil, err
	}
	doc, err := object.GetString("doc")
	if err != nil {
		return nil, err
	}
	target.Doc = doc
	depends, err := object.GetListStrings("depends")
	if err != nil {
		return nil, fmt.Errorf("depends must be a list of strings")
	}
	target.Depends = depends
	list, err := object.GetListStrings("steps")
	if err != nil {
		return nil, fmt.Errorf("steps must be a list of strings or maps")
	}
	var steps []*Step
	for _, object := range list {
		step, err := NewStep(target, object)
		if err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}
	target.Steps = steps
	return target, nil
}

func (target *Target) Run() error {
	if len(target.Depends) > 0 {
		target.Build.Run(target.Depends)
	}
	PrintTarget("Running target " + target.Name)
	for _, step := range target.Steps {
		err := step.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
