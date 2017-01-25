package build

import (
	"fmt"
	"neon/util"
	"os"
)

type Target struct {
	Build   *Build
	Name    string
	Doc     string
	Depends []string
	Steps   []Step
}

func NewTarget(build *Build, name string, object util.Object) (*Target, error) {
	build.Log("Parsing target '%s'", name)
	build.Log("Target structure: %#v", object)
	target := &Target{
		Build: build,
		Name:  name,
	}
	build.Log("Reading target '%s' first level fields", name)
	err := object.CheckFields([]string{"doc", "depends", "steps"})
	if err != nil {
		return nil, err
	}
	if object.HasField("doc") {
		doc, err := object.GetString("doc")
		if err != nil {
			return nil, fmt.Errorf("parsing target %s: doc field must be a string", name)
		}
		target.Doc = doc
	}
	if object.HasField("depends") {
		depends, err := object.GetListStringsOrString("depends")
		if err != nil {
			return nil, fmt.Errorf("parsing target %s: depends field must be a string or list of strings", name)
		}
		target.Depends = depends
	}
	if object.HasField("steps") {
		list, err := object.GetList("steps")
		if err != nil {
			return nil, fmt.Errorf("steps must be a list")
		}
		var steps []Step
		for index, object := range list {
			build.Log("Parsing step %v in target '%s'", index, name)
			step, err := NewStep(target, object)
			if err != nil {
				return nil, fmt.Errorf("parsing target '%s': %v", name, err)
			}
			steps = append(steps, step)
		}
		target.Steps = steps
	}
	return target, nil
}

func (target *Target) Run(stack *Stack) error {
	if stack.Contains(target.Name) {
		stack.Push(target.Name)
		return fmt.Errorf("target circular dependency: %v", stack.ToString())
	}
	stack.Push(target.Name)
	if len(target.Depends) > 0 {
		for _, name := range target.Depends {
			err := target.Build.RunTarget(name, stack)
			if err != nil {
				return err
			}
		}
	}
	util.PrintTarget("Running target " + target.Name)
	err := os.Chdir(target.Build.Dir)
	if err != nil {
		return fmt.Errorf("changing to build directory '%s'", target.Build.Dir)
	}
	for index, step := range target.Steps {
		err := step.Run()
		if err != nil {
			return fmt.Errorf("running target '%s' step %v: %v",
				target.Name, index+1, err)
		}
	}
	return nil
}
