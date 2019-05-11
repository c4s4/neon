package build

import (
	"reflect"
	"testing"
)

func TestScriptStep(t *testing.T) {
	// parse script task
	script := `test = "This is a test"`
	step, err := NewStep(script)
	if err != nil {
		t.Errorf("Error parsing step: %v", err)
	}
	switch step.(type) {
	default:
		t.Errorf("Bad step type")
	case ScriptStep:
		println("Success")
	}
	// prepare build and context
	build := &Build{}
	build.Properties = build.GetProperties()
	build.Environment = build.GetEnvironment()
	build.SetDir(build.Dir)
	context := NewContext(build)
	err = context.Init()
	if err != nil {
		t.Errorf("Error during context init: %v", err)
	}
	// run this step
	err = step.Run(context)
	if err != nil {
		t.Errorf("Error running step: %v", err)
	}
	value, err := context.GetProperty("test")
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}
	if value != "This is a test" {
		t.Errorf("Bad value: %v", value)
	}
}

func TestTaskStep(t *testing.T) {
	// add task to the task map
	TaskMap = make(map[string]TaskDesc)
	type testArgs struct {
		Test string
	}
	AddTask(TaskDesc{
		Name: "test",
		Func: testFunc,
		Args: reflect.TypeOf(testArgs{}),
		Help: `Task documentation.`,
	})
	// parse task step
	task := map[interface{}]interface{}{
		"test": "This is a test",
	}
	step, err := NewStep(task)
	if err != nil {
		t.Errorf("Error parsing step: %v", err)
	}
	switch step.(type) {
	default:
		t.Errorf("Bad step type")
	case TaskStep:
		println("Success")
	}
	// prepare build and context
	build := &Build{}
	build.Properties = build.GetProperties()
	build.Environment = build.GetEnvironment()
	build.SetDir(build.Dir)
	context := NewContext(build)
	err = context.Init()
	if err != nil {
		t.Errorf("Error during context init: %v", err)
	}
	// run this step
	err = step.Run(context)
	if err != nil {
		t.Errorf("Error running step: %v", err)
	}
	value, err := context.GetProperty("test")
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}
	if value != "This is a test" {
		t.Errorf("Bad value: %v", value)
	}
}

func TestSteps(t *testing.T) {
	// add task to the task map
	TaskMap = make(map[string]TaskDesc)
	type testArgs struct {
		Test string
	}
	AddTask(TaskDesc{
		Name: "test",
		Func: testFunc,
		Args: reflect.TypeOf(testArgs{}),
		Help: `Task documentation.`,
	})
	// parse steps
	script := `test2 = "This is another test"`
	task := map[interface{}]interface{}{
		"test": "This is a test",
	}
	steps, err := NewSteps([]interface{}{script, task})
	if err != nil {
		t.Errorf("Error parsing steps: %v", err)
	}
	// prepare build and context
	build := &Build{}
	build.Properties = build.GetProperties()
	build.Environment = build.GetEnvironment()
	build.SetDir(build.Dir)
	context := NewContext(build)
	err = context.Init()
	if err != nil {
		t.Errorf("Error during context init: %v", err)
	}
	// run this step
	err = steps.Run(context)
	if err != nil {
		t.Errorf("Error running steps: %v", err)
	}
	value, err := context.GetProperty("test")
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}
	if value != "This is a test" {
		t.Errorf("Bad value: %v", value)
	}
	value, err = context.GetProperty("test2")
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}
	if value != "This is another test" {
		t.Errorf("Bad value: %v", value)
	}
}
