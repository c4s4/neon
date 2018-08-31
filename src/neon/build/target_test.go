package build

import (
	"reflect"
	"testing"
)

func TestTarget(t *testing.T) {
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
	// prepare build and context
	build := &Build{}
	build.Properties = build.GetProperties()
	build.Environment = build.GetEnvironment()
	build.SetDir(".")
	context := NewContext(build)
	err := context.Init()
	// parse steps
	script := `test2 = "This is another test"`
	task := map[interface{}]interface{}{
		"test": "This is a test",
	}
	object := map[string]interface{}{
		"doc": "Test target",
		"depends": []string{},
		"steps": []interface{}{script, task},
	}
	target, err := NewTarget(build, "test", object)
	if err != nil {
		t.Errorf("Error parsing target: %v", err)
	}
	// run this target
	err = target.Run(context)
	if err != nil {
		t.Errorf("Error running target: %v", err)
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
