package build

import (
	"testing"
	"reflect"
	"fmt"
)

type TestArgs struct {
	String string
	Bool   bool              `optional`
	Int    int               `optional`
	Float  float64           `optional`
	Array  []string          `optional`
	Map    map[string]string `optional`
}

func TestValidateTaskArgsNominal(t *testing.T) {
	args := map[string]interface{} {
		"string": "Hello World!",
		"bool": true,
		"int": 3,
		"float": 3.14,
		"array": []string{"foo", "bar"},
		"map": map[string]string{"foo": "bar"},
	}
	err := ValidateTaskArgs(args, &TestArgs{})
	if err != nil {
		t.Errorf("failed args validation: %#v", err)
	}
}

func TestValidateTaskArgsMissingArg(t *testing.T) {
	args := map[string]interface{} {
		"int": 3,
	}
	err := ValidateTaskArgs(args, &TestArgs{})
	if err == nil || err.Error() != "missing mandatory field 'string'" {
		t.Errorf("failed args validation: %v", err)
	}
}

func TestValidateTaskArgsMissingArgOptional(t *testing.T) {
	args := map[string]interface{} {
		"string": "Hello World!",
	}
	err := ValidateTaskArgs(args, &TestArgs{})
	if err != nil {
		t.Errorf("failed args validation: %#v", err)
	}
}

func TestValidateTaskArgsBadType(t *testing.T) {
	args := map[string]interface{} {
		"string": 1,
	}
	err := ValidateTaskArgs(args, &TestArgs{})
	if err == nil || err.Error() != "field 'string' must be of type 'string' ('int' provided)" {
		t.Errorf("failed args validation")
	}
}

func TestEvaluateTaskArgsNominal(t *testing.T) {
	args := map[string]interface{} {
		"string": "Hello World!",
		"bool": true,
		"int": 3,
		"float": 3.14,
		"array": []string{"foo", "bar"},
		"map": map[string]string{"foo": "bar"},
	}
	params := TestArgs{}
	err := EvaluateTaskArgs(args, &params, nil)
	if err != nil {
		t.Errorf("failed args evaluation: %#v", err)
	}
	if params.String != "Hello World!" || params.Int != 3 || !params.Bool {
		t.Errorf("failed args evaluation: %#v", args)
	}
	if params.Array[0] != "foo" || params.Array[1] != "bar" {
		t.Errorf("failed args evaluation: %#v", args)
	}
	if params.Map["foo"] != "bar" {
		t.Errorf("failed args evaluation: %#v", args)
	}
}

func TestFieldIs(t *testing.T) {
	field := reflect.StructField{Tag: "test"}
	if !FieldIs(field, "test") {
		t.Errorf("failed FieldIs test")
	}
	if FieldIs(field, "foo") {
		t.Errorf("failed FieldIs test")
	}
	field = reflect.StructField{Tag: "foo bar"}
	if !FieldIs(field, "foo") || !FieldIs(field, "bar") {
		t.Errorf("failed FieldIs test")
	}
	if FieldIs(field, "test") {
		t.Errorf("failed FieldIs test")
	}
}

// This test demonstrates how to check task parameters, fill them with
// arguments from build file, define a task and call it with parameters
func TestTaskCall(t *testing.T) {
	// task arguments as parsed in build file
	args := map[string]interface{} {
		"print": "Hello World!",
	}
	// task arguments type
	type PrintArgs struct {
		Print string
	}
	// the task function
	print := func(ctx *Context, args interface{}) error {
		params := args.(*PrintArgs)
		fmt.Println(params.Print)
		return nil
	}
	// the task argument type
	//var typ reflect.Type = reflect.TypeOf(PrintArgs{})
	// get an instance of arguments type
	//var params reflect.Value = reflect.New(typ).Elem()
	params := PrintArgs{}
	// validate task arguments
	err := ValidateTaskArgs(args, &params)
	if err != nil {
		t.Errorf("failed args validation: %v", err)
	}
	// evaluate task arguments
	err = EvaluateTaskArgs(args, &params, nil)
	if err != nil {
		t.Errorf("failed args evaluation: %v", err)
	}
	if params.Print != "Hello World!" {
		t.Errorf("bad args values: %v", err)
	}
	err = print(nil, &params)
	if err != nil {
		t.Errorf("error calling task: %v", err)
	}
}
