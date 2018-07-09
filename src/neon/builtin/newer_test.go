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
	time.Sleep(100 * time.Millisecond)
	Touch(new)
	if !newer(new, old) {
		t.Errorf("Newer error")
	}
}
