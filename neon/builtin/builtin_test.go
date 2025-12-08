package builtin

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"
	"time"

	_build "github.com/c4s4/neon/neon/build"
	_ "github.com/c4s4/neon/neon/task"
	"github.com/c4s4/neon/neon/util"
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

// Assert make an assertion for testing purpose, failing test if different:
// - actual: actual value
// - expected: regexp for expected value
// - t: test
func AssertRegexp(actual, expected string, t *testing.T) {
	match, err := regexp.MatchString(expected, actual)
	if err != nil {
		t.Errorf("error matching regexp: %v", err)
	}
	if !match {
		t.Errorf("actual (\"%s\") !~ expected (\"%s\")", actual, expected)
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
	defer func() {
		_ = os.Chdir(dir)
	}()
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
	defer func() {
		_ = os.Chdir(dir)
	}()
	build, err := _build.NewBuild(file, filepath.Dir(file), "", false)
	if err != nil {
		return err
	}
	if err := os.Chdir(build.Dir); err != nil {
		return fmt.Errorf("changing directory to '%s': %v", build.Dir, err)
	}
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
		err := os.WriteFile(file, []byte{}, 0755)
		if err != nil {
			return fmt.Errorf("creating file '%s': %v", file, err)
		}
	}
	return nil
}
