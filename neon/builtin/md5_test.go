package builtin

import (
	"os"
	"testing"
)

const file = "/tmp/test.txt"

func TestMD5(t *testing.T) {
	os.WriteFile(file, []byte("test"), 0644)
	defer os.Remove(file)
	expected := "098f6bcd4621d373cade4e832627b4f6"
	actual := md5Sum(file)
	if actual != expected {
		t.Fatalf("bad MD5 sum: expected %s, got %s", expected, actual)
	}
}
