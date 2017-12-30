package build

import (
	"strings"
	"testing"
	"reflect"
)

func TestEvaluateString(t *testing.T) {
	context := NewContext()
	if actual, err := context.EvaluateString(`foo`); err != nil || actual != `foo` {
		t.Errorf("TestEvaluateString failure")
	}
	if actual, err := context.EvaluateString(`foo #{"bar"}`); err != nil || actual != `foo bar` {
		t.Errorf("TestEvaluateString failure")
	}
	if actual, err := context.EvaluateString(`foo ={"bar"}`); err != nil || actual != `foo bar` {
		t.Errorf("TestEvaluateString failure")
	}
	if actual, err := context.EvaluateString(`foo ={1+1}`); err != nil || actual != `foo 2` {
		t.Errorf("TestEvaluateString failure")
	}
}

func TestEvaluateStringWithProperties(t *testing.T) {
	context := NewContext()
	properties := map[string]interface{}{
		"FOO": "foo",
		"BAR": "bar",
	}
	build := &Build {
		Dir:        "dir",
		Here:       "here",
		Properties: properties,
	}
	context.InitProperties(build)
	if actual, err := context.EvaluateString(`foo`); err != nil || actual != `foo` {
		t.Errorf("TestEvaluateStringWithProperties failure")
	}
	if actual, err := context.EvaluateString(`={FOO} bar`); err != nil || actual != `foo bar` {
		t.Errorf("TestEvaluateStringWithProperties failure")
	}
	if actual, err := context.EvaluateString(`={FOO} ={BAR}`); err != nil || actual != `foo bar` {
		t.Errorf("TestEvaluateStringWithProperties failure")
	}
	if _, err := context.EvaluateString(`={XXX}`); err == nil || err.Error() != `Undefined symbol 'XXX' (at line 1, column 1)` {
		t.Errorf("TestEvaluateStringWithProperties failure")
	}
}

func TestEvaluateSliceWithProperties(t *testing.T) {
	context := NewContext()
	properties := map[string]interface{}{
		"FOO": "foo",
		"BAR": "bar",
	}
	build := &Build{
		Dir:        "dir",
		Here:       "here",
		Properties: properties,
	}
	context.InitProperties(build)
	actual, err := context.EvaluateObject([]string{`={FOO} BAR`, `FOO ={BAR}`})
	if err != nil {
		t.Fail()
	}
	if reflect.TypeOf(actual) != reflect.SliceOf(reflect.TypeOf("")) {
		t.Fail()
	}
	value := reflect.ValueOf(actual)
	if value.Len() != 2 {
		t.Fail()
	}
	if value.Index(0) == reflect.ValueOf(`foo BAR`) {
		t.Fail()
	}
	if value.Index(1) == reflect.ValueOf(`FOO bar`) {
		t.Fail()
	}
}

func TestEvaluateMapWithProperties(t *testing.T) {
	context := NewContext()
	properties := map[string]interface{}{
		"FOO": "foo",
		"BAR": "bar",
	}
	build := &Build{
		Dir:        "dir",
		Here:       "here",
		Properties: properties,
	}
	context.InitProperties(build)
	actual, err := context.EvaluateObject(map[string]string{"={FOO}": "BAR", "FOO": "={BAR}"})
	if err != nil {
		t.Fail()
	}
	if reflect.TypeOf(actual) != reflect.TypeOf(make(map[string]string)) {
		t.Fail()
	}
	value := reflect.ValueOf(actual)
	if value.Len() != 2 {
		t.Fail()
	}
	if value.MapIndex(reflect.ValueOf("foo")) == reflect.ValueOf("BAR") {
		t.Fail()
	}
	if value.MapIndex(reflect.ValueOf("FOO")) == reflect.ValueOf("bar") {
		t.Fail()
	}
}

func TestGetSetProperty(t *testing.T) {
	context := NewContext()
	context.SetProperty("foo", "bar")
	if p, err := context.GetProperty("foo"); p != "bar" || err != nil {
		t.Fail()
	}
}

func TestEvaluateExpression(t *testing.T) {
	context := NewContext()
	_, err := context.EvaluateExpression(`foo = "BAR"`)
	if err != nil {
		t.Fail()
	}
	r, err := context.GetProperty("foo")
	if err != nil || r != "BAR" {
		t.Fail()
	}
	if r, err = context.EvaluateExpression(`1+2`); err != nil || r != int64(3) {
		t.Fail()
	}
}

func TestEvaluateEnvironmentSimple(t *testing.T) {
	context := NewContext()
	build := &Build{
		Dir: "dir",
		Environment: map[string]string{
			"FOO": "BAR",
		},
	}
	env, err := context.EvaluateEnvironment(build)
	if err != nil {
		t.Errorf("Error getting environment: %v", err)
	}
	for _, line := range env {
		if line == "FOO=BAR" {
			return
		}
	}
	t.Error("Env FOO=BAR not found")
}

func TestEvaluateEnvironmentComplex(t *testing.T) {
	context := NewContext()
	build := &Build{
		Dir: "dir",
		Environment: map[string]string{
			"FOO": "BAR:${HOME}",
		},
	}
	env, err := context.EvaluateEnvironment(build)
	if err != nil {
		t.Errorf("Error getting environment: %v", err)
	}
	var foo string
	for _, line := range env {
		if strings.HasPrefix(line, "HOME=") {
			foo = "FOO=BAR:" + line[5:]
		}
	}
	if foo == "" {
		return
	}
	for _, line := range env {
		if line == foo {
			return
		}
	}
	t.Error("Environment variable FOO not set correctly")
}
