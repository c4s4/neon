package build

import (
	"io/ioutil"
	"neon/util"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// WriteFile writes a file in given directory.
func WriteFile(dir, file, content string) string {
	if !util.DirExists(dir) {
		os.MkdirAll(dir, util.DirFileMode)
	}
	path := filepath.Join(dir, file)
	ioutil.WriteFile(path, []byte(content), util.FileMode)
	return path
}

// Assert make an assertion for testing purpose, failing test if different:
// - actual: actual value
// - expected: expected value
// - t: test
func Assert(actual, expected interface{}, t *testing.T) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual (\"%s\") != expected (\"%s\")", actual, expected)
	}
}
