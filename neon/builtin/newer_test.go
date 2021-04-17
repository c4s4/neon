package builtin

import (
	"os"
	"testing"
	"time"
)

func TestNewer(t *testing.T) {
	testDir := BuildDir + "/builtins/newer"
	os.MkdirAll(testDir, 0755)
	old := testDir + "/old"
	new := testDir + "/new"
	Touch(old)
	time.Sleep(1 * time.Millisecond)
	Touch(new)
	if !newer(new, old) {
		t.Errorf("Newer error")
	}
	if newer(old, new) {
		t.Errorf("Newer error")
	}
}

func TestNewerMulti(t *testing.T) {
	testDir := BuildDir + "/builtins/newer"
	os.MkdirAll(testDir, 0755)
	old1 := testDir + "/old1"
	old2 := testDir + "/old2"
	new1 := testDir + "/new1"
	new2 := testDir + "/new2"
	Touch(old1)
	Touch(old2)
	time.Sleep(1 * time.Millisecond)
	Touch(new1)
	Touch(new2)
	if !newer([]string{new1, new2}, []string{old1, old2}) {
		t.Errorf("Newer error")
	}
	if newer([]string{old1, old2}, []string{new1, new2}) {
		t.Errorf("Newer error")
	}
}
