package build

import (
	"io/ioutil"
	"neon/util"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestInfo(t *testing.T) {
	build := &Build{
		Doc:        "Test documentation",
		Default:    []string{"default"},
		Repository: "repository",
		Singleton:  "12345",
		Extends:    []string{"foo", "bar"},
		Config:     []string{"foo", "bar"},
		//Scripts:    []string{"foo", "bar"},
		Targets: map[string]*Target{
			"test1": {
				Doc:     "Test 1 doc",
				Depends: []string{"foo", "bar"},
			},
			"test2": {
				Doc: "Test 2 doc",
			},
		},
		Environment: map[string]string {
			"FOO": "SPAM",
			"BAR": "EGGS",
		},
		Properties: map[string]interface{}{
			"foo": "spam",
			"bar": "eggs",
		},
	}
	context := NewContext(build)
	err := context.Init()
	if err != nil {
		t.Errorf("Failure: %v", err)
	}
	//build.Properties = build.GetProperties()
	expected := `doc: Test documentation
default: [default]
repository: repository
singleton: 12345
extends:
- foo
- bar
configuration:
- foo
- bar

environment:
  BAR: "EGGS"
  FOO: "SPAM"

properties:
  bar: "eggs"
  foo: "spam"

targets:
  test1: Test 1 doc [foo, bar]
  test2: Test 2 doc`
	info, err := build.Info(context)
	if err != nil {
		t.Errorf("Failure: %v", err)
	}
	if info != expected {
		t.Errorf("Bad build info: %s", info)
	}
}

func TestInfoDoc(t *testing.T) {
	build := Build{
		Doc: "Test documentation",
	}
	if build.infoDoc() != "doc: Test documentation\n" {
		t.Errorf("Bad build doc: %s", build.infoDoc())
	}
}

func TestInfoDefault(t *testing.T) {
	build := Build{
		Default: []string{"default"},
	}
	if build.infoDefault() != "default: [default]\n" {
		t.Errorf("Bad build default: %s", build.infoDefault())
	}
}

func TestInfoRepository(t *testing.T) {
	build := Build{
		Repository: "repository",
	}
	if build.infoRepository() != "repository: repository\n" {
		t.Errorf("Bad build repository: %s", build.infoRepository())
	}
}

func TestInfoSingleton(t *testing.T) {
	build := &Build{
		Singleton: "12345",
	}
	context := NewContext(build)
	if build.infoSingleton(context) != "singleton: 12345\n" {
		t.Errorf("Bad build singleton: %s", build.infoSingleton(context))
	}
}

func TestInfoExtends(t *testing.T) {
	build := &Build{
		Extends: []string{"foo", "bar"},
	}
	if build.infoExtends() != "extends:\n- foo\n- bar\n" {
		t.Errorf("Bad build extends: %s", build.infoExtends())
	}
}

func TestInfoConfiguration(t *testing.T) {
	build := &Build{
		Config: []string{"foo", "bar"},
	}
	if build.infoConfiguration() != "configuration:\n- foo\n- bar\n" {
		t.Errorf("Bad build config: %s", build.infoConfiguration())
	}
}

func TestInfoContext(t *testing.T) {
	build := &Build{
		Scripts: []string{"foo", "bar"},
	}
	if build.infoContext() != "context:\n- foo\n- bar\n" {
		t.Errorf("Bad build context: %s", build.infoContext())
	}
}

func TestInfoTargets(t *testing.T) {
	build := &Build{
		Targets: map[string]*Target{
			"test1": {
				Doc:     "Test 1 doc",
				Depends: []string{"foo", "bar"},
			},
			"test2": {
				Doc: "Test 2 doc",
			},
		},
	}
	expected := `targets:
  test1: Test 1 doc [foo, bar]
  test2: Test 2 doc
`
	if build.infoTargets() != expected {
		t.Errorf("Bad targets info: '%s'", build.infoTargets())
	}
}

func TestInfoTasks(t *testing.T) {
	TaskMap = make(map[string]TaskDesc)
	type testArgs struct {
		Test string
	}
	AddTask(TaskDesc{
		Name: "task",
		Func: testFunc,
		Args: reflect.TypeOf(testArgs{}),
		Help: `Task documentation.`,
	})
	tasks := InfoTasks()
	if tasks != "task" {
		t.Errorf("Bad tasks: %s", tasks)
	}
}

func TestInfoTask(t *testing.T) {
	TaskMap = make(map[string]TaskDesc)
	type testArgs struct {
		Test string
	}
	AddTask(TaskDesc{
		Name: "task",
		Func: testFunc,
		Args: reflect.TypeOf(testArgs{}),
		Help: `Task documentation.`,
	})
	task := InfoTask("task")
	if task != "Task documentation." {
		t.Errorf("Bad task: %s", task)
	}
}

func TestInfoBuiltins(t *testing.T) {
	BuiltinMap = make(map[string]BuiltinDesc)
	AddBuiltin(BuiltinDesc{
		Name: "test",
		Func: TestInfoBuiltins,
		Help: `Test documentation.`,
	})
	builtins := InfoBuiltins()
	if builtins != "test" {
		t.Errorf("Bad builtins: %s", builtins)
	}
}

func TestInfoBuiltin(t *testing.T) {
	BuiltinMap = make(map[string]BuiltinDesc)
	AddBuiltin(BuiltinDesc{
		Name: "test",
		Func: TestInfoBuiltins,
		Help: `Test documentation.`,
	})
	info := InfoBuiltin("test")
	if info != "Test documentation." {
		t.Errorf("Bad builtin info: %s", info)
	}
}

func TestInfoThemes(t *testing.T) {
	themes := InfoThemes()
	if themes != "bee blue bold cyan fire green magenta marine nature red reverse rgb yellow" {
		t.Errorf("Bad themes")
	}
}

func TestInfoTemplates(t *testing.T) {
	repo := "/tmp/neon"
	writeFile(repo+"/foo/bar", "template1.tpl")
	writeFile(repo+"/foo/bar", "template2.tpl")
	defer os.RemoveAll(repo)
	parents := InfoTemplates(repo)
	if parents != "foo/bar/template1.tpl\nfoo/bar/template2.tpl" {
		t.Errorf("Bad templates info: %s", parents)
	}
}

func TestInfoParents(t *testing.T) {
	repo := "/tmp/neon"
	writeFile(repo+"/foo/bar", "parent1.yml")
	writeFile(repo+"/foo/bar", "parent2.yml")
	defer os.RemoveAll(repo)
	parents := InfoParents(repo)
	if parents != "foo/bar/parent1.yml\nfoo/bar/parent2.yml" {
		t.Errorf("Bad parents info: %s", parents)
	}
}

func testFunc(context *Context, args interface{}) error {
	return nil
}

func TestInfoReference(t *testing.T) {
	BuiltinMap = make(map[string]BuiltinDesc)
	AddBuiltin(BuiltinDesc{
		Name: "builtin",
		Func: TestInfoReference,
		Help: `Builtin documentation.`,
	})
	type testArgs struct {
		Test string
	}
	TaskMap = make(map[string]TaskDesc)
	AddTask(TaskDesc{
		Name: "task",
		Func: testFunc,
		Args: reflect.TypeOf(testArgs{}),
		Help: `Task documentation.`,
	})
	actual := InfoReference()
	expected := `Tasks Reference
===============

task
----

Task documentation.

Builtins Reference
==================

builtin
-------

Builtin documentation.`
	if actual != expected {
		t.Errorf("Bad reference: %s", actual)
	}
}

// Utility functions

func writeFile(dir string, file string) string {
	if !util.DirExists(dir) {
		os.MkdirAll(dir, util.DirFileMode)
	}
	path := filepath.Join(dir, file)
	ioutil.WriteFile(path, []byte("test"), util.FileMode)
	return path
}
