package build

import (
	"os"
	"reflect"
	"runtime"
	"testing"
)

func TestNewBuild(t *testing.T) {
	// write the build file
	buildFile := `doc: Test build file

properties:
    FOO: 'foo'

targets:

    test:
        doc: Test target
        steps:
        - test: 'This is a test'`
	if _, err := WriteFile("/tmp", "build.yml", buildFile); err != nil {
		t.Fatalf("write file: %v", err)
	}
	defer func() {
		_ = os.Remove("/tmp/build.yml")
	}()
	// define test task
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
	// load the build file
	build, err := NewBuild("/tmp/build.yml", "/tmp", "", false)
	if err != nil {
		t.Errorf("Error parsing build file: %v", err)
	}
	if build.Dir != "/tmp" {
		t.Errorf("Bad build dir: %s", build.Dir)
	}
}

func TestGetShell(t *testing.T) {
	build := &Build{
		Shell: map[string][]string{
			runtime.GOOS: {"foo"},
			"other":      {"bar"},
		},
	}
	shell, err := build.GetShell()
	if err != nil {
		t.Fail()
	}
	Assert(shell, []string{"foo"}, t)
}
