package builtin

import (
	"os"
	"testing"
)

func TestFind(t *testing.T) {
	testDir := BuildDir + "/builtins/find"
	os.MkdirAll(testDir, 0755)
	os.MkdirAll(testDir+"/foo", 0755)
	Touch(testDir + "/spam.txt")
	Touch(testDir + "/foo/eggs.txt")
	files := find(testDir, "*.txt")
	if len(files) != 1 {
		t.Errorf("Got %d files while expecting 1", len(files))
	}
	files = find(testDir, "**/*.txt")
	if len(files) != 2 {
		t.Errorf("Got %d files while expecting 1", len(files))
	}
}
