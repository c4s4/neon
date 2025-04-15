package builtin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFollowLink(t *testing.T) {
	testDir := BuildDir + "/builtins/followlink"
	_ = os.MkdirAll(testDir, 0755)
	src, err := filepath.Abs(testDir + "/spam.txt")
	if err != nil {
		t.Error(err)
	}
	dst, err := filepath.Abs(testDir + "/eggs.txt")
	if err != nil {
		t.Error(err)
	}
	_ = Touch(src)
	_ = os.Symlink(src, dst)
	path := followLink(dst)
	if path != src {
		t.Errorf("Symlink not followed %s instead of %s", path, src)
	}
}
