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
	args := TaskArgs {
		"string": "Hello World!",
		"bool": true,
		"int": 3,
		"float": 3.14,
		"array": []string{"foo", "bar"},
		"map": map[string]string{"foo": "bar"},
	}
	err := ValidateTaskArgs(args, reflect.TypeOf(TestArgs{}))
	if err != nil {
		t.Errorf("failed args validation: %#v", err)
	}
}

func TestValidateTaskArgsMissingArg(t *testing.T) {
	args := TaskArgs {
		"int": 3,
	}
	err := ValidateTaskArgs(args, reflect.TypeOf(TestArgs{}))
	if err == nil || err.Error() != "missing mandatory field 'string'" {
		t.Errorf("failed args validation: %v", err)
	}
}

func TestValidateTaskArgsMissingArgOptional(t *testing.T) {
	args := TaskArgs {
		"string": "Hello World!",
	}
	err := ValidateTaskArgs(args, reflect.TypeOf(TestArgs{}))
	if err != nil {
		t.Errorf("failed args validation: %#v", err)
	}
}

func TestValidateTaskArgsBadType(t *testing.T) {
	args := TaskArgs {
		"string": 1,
	}
	err := ValidateTaskArgs(args, reflect.TypeOf(TestArgs{}))
	if err == nil || err.Error() != "field 'string' must be of type 'string' ('int' provided)" {
		t.Errorf("failed args validation")
	}
}

func TestEvaluateTaskArgsNominal(t *testing.T) {
	args := TaskArgs {
		"string": "Hello World!",
		"bool": true,
		"int": 3,
		"float": 3.14,
		"array": []string{"foo", "bar"},
		"map": map[string]string{"foo": "bar"},
	}
	res, err := EvaluateTaskArgs(args, reflect.TypeOf(TestArgs{}), nil)
	params := res.(TestArgs)
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

func TestIsExpression(t *testing.T) {
	if !IsExpression("=foo") {
		t.Errorf("failed IsExpression test")
	}
	if IsExpression("foo") {
		t.Errorf("failed IsExpression test")
	}
	if IsExpression("={foo}") {
		t.Errorf("failed IsExpression test")
	}
}

// This test demonstrates how to check task parameters, fill them with
// arguments from build file, define a task and call it with parameters
func TestTaskCall(t *testing.T) {
	// task arguments as parsed in build file
	args := TaskArgs {
		"print": "Hello World!",
	}
	// task arguments type
	type PrintArgs struct {
		Print string
	}
	// the task function
	print := func(ctx *Context, args interface{}) error {
		params := args.(PrintArgs)
		fmt.Println(params.Print)
		return nil
	}
	// task arguments type
	typ := reflect.TypeOf(PrintArgs{})
	// validate task arguments
	err := ValidateTaskArgs(args, typ)
	if err != nil {
		t.Errorf("failed args validation: %v", err)
	}
	// evaluate task arguments
	params, err := EvaluateTaskArgs(args, typ, nil)
	if err != nil {
		t.Errorf("failed args evaluation: %v", err)
	}
	if params.(PrintArgs).Print != "Hello World!" {
		t.Errorf("bad args values: %v", err)
	}
	err = print(nil, params)
	if err != nil {
		t.Errorf("error calling task: %v", err)
	}
}

func TestValueOfType(t *testing.T) {
	if !ValueOfType(1, reflect.TypeOf(1)) {
		t.Fail()
	}
	if !ValueOfType("foo", reflect.TypeOf("foo")) {
		t.Fail()
	}
	if !ValueOfType([]interface{}{"foo", "bar"}, reflect.TypeOf([]string{"foo"})) {
		t.Fail()
	}
	if !ValueOfType(map[string]interface{}{"foo": "bar"},
					reflect.TypeOf(map[string]string{"foo": "bar"})) {
		t.Fail()
	}
}
