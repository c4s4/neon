package util

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCopyFileOntoItself(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(file, []byte("content"), FileMode); err != nil {
		t.Fatal(err)
	}
	if err := CopyFile(file, file); err == nil {
		t.Errorf("copying a file onto itself should fail")
	}
	content, err := os.ReadFile(file)
	if err != nil || string(content) != "content" {
		t.Errorf("file should be unchanged: %q", content)
	}
}

func TestExpandUserHomeTilde(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no user home")
	}
	if ExpandUserHome("~") != home {
		t.Errorf("~ should be expanded: %s", ExpandUserHome("~"))
	}
	if ExpandUserHome("~/foo") != filepath.Join(home, "foo") {
		t.Errorf("~/foo should be expanded: %s", ExpandUserHome("~/foo"))
	}
	if ExpandUserHome("~foo") != "~foo" {
		t.Errorf("~foo should not be expanded")
	}
}

func TestFindFilesDuplicatesAndBrokenLinks(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte(""), FileMode); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "missing"), filepath.Join(dir, "broken.go")); err != nil {
		t.Skip("symbolic links not supported")
	}
	files, err := FindFiles(dir, []string{"*.go", "**/*.go"}, nil, false)
	if err != nil {
		t.Fatalf("finding files: %v", err)
	}
	if !reflect.DeepEqual(files, []string{"a.go"}) {
		t.Errorf("bad files: %v", files)
	}
}

func TestMoveFile(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source.txt")
	dest := filepath.Join(dir, "dest.txt")
	if err := os.WriteFile(source, []byte("content"), FileMode); err != nil {
		t.Fatal(err)
	}
	if err := moveFile(source, dest); err != nil {
		t.Fatalf("moving file: %v", err)
	}
	if FileExists(source) || !FileExists(dest) {
		t.Errorf("file should be moved")
	}
}

func TestFindFilesBadPattern(t *testing.T) {
	if _, err := FindFiles(t.TempDir(), []string{"[bad"}, nil, false); err == nil {
		t.Errorf("invalid pattern should be reported")
	}
}
