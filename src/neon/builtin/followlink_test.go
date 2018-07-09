package builtin

import (
	"os"
	"testing"
)

func TestFollowLink(t *testing.T) {
	testDir := BuildDir + "/builtins/followlink"
	os.MkdirAll(testDir, 0755)
	src := testDir + "/spam.txt"
	dst := testDir + "/eggs.txt"
	Touch(src)
	os.Symlink(src, dst)
	path := followLink(dst)
	if path != src {
		t.Errorf("Symlink not folowed %s instead of %s", path, src)
	}
}
