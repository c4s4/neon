package builtin

import (
	"fmt"
	"io/ioutil"
	_build "neon/build"
	_ "neon/task"
	"neon/util"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

const (
	TestDir  = "../../../test/builtin"
	BuildDir = "../../../build"
)

// Assert make an assertion for testing purpose, failing test if different:
// - actual: actual value
// - expected: expected value
// - t: test
func Assert(actual, expected interface{}, t *testing.T) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual (\"%s\") != expected (\"%s\")", actual, expected)
	}
}

// TestIntegration runs all test build files (in test directory).
func TestIntegration(t *testing.T) {
	builds, err := util.FindFiles(TestDir, []string{"**/*.yml"}, []string{}, false)
	if err != nil {
		t.Errorf("Error getting test build files: %v", err)
	}
	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting current directory: %v", err)
	}
	defer os.Chdir(dir)
	for _, build := range builds {
		file, _ := filepath.Abs(path.Join(TestDir, build))
		err := RunBuildFile(file, dir)
		if err != nil {
			t.Error(err)
		}
	}
}

// RunBuildFile runs given build file
func RunBuildFile(file, dir string) error {
	defer os.Chdir(dir)
	build, err := _build.NewBuild(file, filepath.Dir(file), "", false)
	if err != nil {
		return err
	}
	os.Chdir(build.Dir)
	context := _build.NewContext(build)
	err = context.Init()
	if err != nil {
		return err
	}
	err = build.Run(context, []string{})
	if err != nil {
		return err
	}
	return nil
}

// Touch touches given file:
// - file: the name of the file
// Return error if something went wrong
func Touch(file string) error {
	if util.FileExists(file) {
		time := time.Now()
		err := os.Chtimes(file, time, time)
		if err != nil {
			return fmt.Errorf("changing times of file '%s': %v", file, err)
		}
	} else {
		err := ioutil.WriteFile(file, []byte{}, 0755)
		if err != nil {
			return fmt.Errorf("creating file '%s': %v", file, err)
		}
	}
	return nil
}
