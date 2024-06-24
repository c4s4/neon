package build

import (
	"fmt"
	"os"
	"reflect"

	"github.com/c4s4/neon/neon/util"
)

// Target is a structure for a target
type Target struct {
	Build   *Build
	Name    string
	Doc     string
	Depends []string
	Unless  string
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
	if err := object.CheckFields([]string{"doc", "depends", "unless", "steps"}); err != nil {
		return nil, err
	}
	if err := ParseTargetDoc(object, target); err != nil {
		return nil, err
	}
	if err := ParseTargetDepends(object, target); err != nil {
		return nil, err
	}
	if err := ParseTargetUnless(object, target); err != nil {
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

// ParseTargetUnless parses unless clause of the target:
// - object: body of the target as an interface
// - target: the target to document
// Return: an error if something went wrong
func ParseTargetUnless(object util.Object, target *Target) error {
	if object.HasField("unless") {
		unless, err := object.GetString("unless")
		if err != nil {
			return fmt.Errorf("unless field in target '%s' must be a string", target.Name)
		}
		target.Unless = unless
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
	// if unless expression returns true, we skip this target
	if target.Unless != "" {
		object, err := context.EvaluateExpression(target.Unless)
		if err != nil {
			return fmt.Errorf("evaluating unless clause of target %s: %v", target.Name, err)
		}
		value := reflect.ValueOf(object)
		if value.Kind() != reflect.Bool {
			return fmt.Errorf("unless clause expression must return a boolean")
		}
		unless := object.(bool)
		if unless {
			Title(target.Name)
			Message("Skipping target, unless clause was matched")
			return nil
		}
	}
	if err := context.Stack.Push(target); err != nil {
		return err
	}
	for _, name := range target.Depends {
		if err := target.Build.Root.RunTarget(context, name); err != nil {
			return err
		}
	}
	Title(target.Name)
	if target.Build.Template {
		if err := os.Chdir(target.Build.Here); err != nil {
			return fmt.Errorf("changing to current directory '%s'", target.Build.Dir)
		}
	} else {
		if err := os.Chdir(target.Build.Dir); err != nil {
			return fmt.Errorf("changing to build directory '%s'", target.Build.Dir)
		}
	}
	run_err := target.Steps.Run(context)
	if err := context.Stack.Pop(); err != nil {
		return err
	}
	return run_err
}
