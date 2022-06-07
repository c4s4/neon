package build

import (
	"github.com/c4s4/neon/neon/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// WriteFile writes a file in given directory.
func WriteFile(dir, file, content string) (string, error) {
	if !util.DirExists(dir) {
		if err := os.MkdirAll(dir, util.DirFileMode); err != nil {
			return "", err
		}
	}
	path := filepath.Join(dir, file)
	if err := ioutil.WriteFile(path, []byte(content), util.FileMode); err != nil {
		return "", err
	}
	return path, nil
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
