package builtin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFollowLink(t *testing.T) {
	testDir := BuildDir + "/builtins/followlink"
	os.MkdirAll(testDir, 0755)
	src, err := filepath.Abs(testDir + "/spam.txt")
	if err != nil {
		t.Errorf(err.Error())
	}
	dst, err := filepath.Abs(testDir + "/eggs.txt")
	if err != nil {
		t.Errorf(err.Error())
	}
	Touch(src)
	os.Symlink(src, dst)
	path := followLink(dst)
	if path != src {
		t.Errorf("Symlink not folowed %s instead of %s", path, src)
	}
}
