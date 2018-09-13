package build

import (
	"fmt"
	"neon/util"
	"os"
)

// Target is a structure for a target
type Target struct {
	Build   *Build
	Name    string
	Doc     string
	Depends []string
	Steps   Steps
}

// NewTarget makes a new target:
// - build: the build of the target
// - name: the name of the target
// - object: the body of the target as an interface
// Returns:
// - a pointer to the built target
// - an error if something went wrong
func NewTarget(build *Build, name string, object util.Object) (*Target, error) {
	target := &Target{
		Build: build,
		Name:  name,
	}
	err := object.CheckFields([]string{"doc", "depends", "steps"})
	if err != nil {
		return nil, err
	}
	if err := ParseTargetDoc(object, target); err != nil {
		return nil, err
	}
	if err := ParseTargetDepends(object, target); err != nil {
		return nil, err
	}
	if err := ParseTargetSteps(object, target); err != nil {
		return nil, err
	}
	return target, nil
}

// ParseTargetDoc parses documentation of the target:
// - object: body of the target as an interface
// - target: the target to document
// Return: an error if something went wrong
func ParseTargetDoc(object util.Object, target *Target) error {
	if object.HasField("doc") {
		doc, err := object.GetString("doc")
		if err != nil {
			return fmt.Errorf("doc field in target '%s' must be a string", target.Name)
		}
		target.Doc = doc
	}
	return nil
}

// ParseTargetDepends parses target dependencies:
// - object: the target body as an interface
// - target: the target being parsed
// Return: an error if something went wrong
func ParseTargetDepends(object util.Object, target *Target) error {
	if object.HasField("depends") {
		depends, err := object.GetListStringsOrString("depends")
		if err != nil {
			return fmt.Errorf("depends field must be a string or list of strings")
		}
		target.Depends = depends
	}
	return nil
}

// ParseTargetSteps parses steps of a target:
// - object: the target body as an interface
// - target: the target being parsed
// Return: an error if something went wrong
func ParseTargetSteps(object util.Object, target *Target) error {
	if object.HasField("steps") {
		list, err := object.GetList("steps")
		if err != nil {
			return fmt.Errorf("parsig target '%s': steps must be a list", target.Name)
		}
		var steps []Step
		for index, object := range list {
			step, err := NewStep(object)
			if err != nil {
				return fmt.Errorf("in step %d: %v", index+1, err)
			}
			steps = append(steps, step)
		}
		target.Steps = steps
	}
	return nil
}

// Run target in given context:
// - context: the context of the build
// Return: an error if something went wrong
func (target *Target) Run(context *Context) error {
	for _, name := range target.Depends {
		if !context.Stack.Contains(name) {
			err := target.Build.RunTarget(context, name)
			if err != nil {
				return err
			}
		}
	}
	err := context.Stack.Push(target)
	if err != nil {
		return err
	}
	Title(target.Name)
	err = os.Chdir(target.Build.Dir)
	if err != nil {
		return fmt.Errorf("changing to build directory '%s'", target.Build.Dir)
	}
	err = target.Steps.Run(context)
	if err != nil {
		return err
	}
	return nil
}
