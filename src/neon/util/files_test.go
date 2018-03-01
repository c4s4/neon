package util

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
	"testing"
)

const (
	FileMode = 0644
)

func makeTempFile(dir string, t *testing.T) string {
	tempFile, err := ioutil.TempFile(dir, "files_test.tmp")
	if err != nil {
		t.Fail()
	}
	return tempFile.Name()
}

func writeTempFile(dir string, t *testing.T) string {
	tempFile := makeTempFile(dir, t)
	err := ioutil.WriteFile(tempFile, []byte("test"), FileMode)
	if err != nil {
		t.Fail()
	}
	return tempFile
}

func TestReadFile(t *testing.T) {
	tempFile := writeTempFile("", t)
	defer os.Remove(tempFile)
	text, err := ReadFile(tempFile)
	if err != nil {
		t.Fail()
	}
	if string(text) != "test" {
		t.Fail()
	}
	_, err = ReadFile("file_that_doesnt_exist")
	if err == nil {
		t.Fail()
	}
}

func TestFileExists(t *testing.T) {
	tempFile := writeTempFile("", t)
	defer os.Remove(tempFile)
	if !FileExists(tempFile) {
		t.Fail()
	}
	if FileExists("file_that_doesnt_exist") {
		t.Fail()
	}
}

func TestDirExists(t *testing.T) {
	if !DirExists("../util") {
		t.Fail()
	}
	if DirExists("dir_that_doesnt_exist") {
		t.Fail()
	}
}

func TestCopyFile(t *testing.T) {
	srcFile := writeTempFile("", t)
	defer os.Remove(srcFile)
	dstFile := path.Join(os.TempDir(), "test.tmp")
	defer os.Remove(dstFile)
	err := CopyFile(srcFile, dstFile)
	if err != nil {
		t.Fail()
	}
	text, err := ReadFile(dstFile)
	if err != nil {
		t.Fail()
	}
	if string(text) != "test" {
		t.Fail()
	}
}

func TestCopyFilesToDir(t *testing.T) {
	srcDir := path.Join(os.TempDir(), "src")
	dstDir := path.Join(os.TempDir(), "dst")
	if err := os.MkdirAll(srcDir, DirFileMode); err != nil {
		t.Fail()
	}
	if err := os.MkdirAll(dstDir, DirFileMode); err != nil {
		t.Fail()
	}
	tmpFile1 := makeTempFile(srcDir, t)
	println(tmpFile1)
	defer os.Remove(tmpFile1)
	tmpFile2 := makeTempFile(srcDir, t)
	println(tmpFile2)
	defer os.Remove(tmpFile2)
	err := CopyFilesToDir(srcDir, []string{path.Base(tmpFile1), path.Base(tmpFile2)}, dstDir, false)
	if err != nil {
		t.Fail()
	}
	newFile1 := path.Join(dstDir, path.Base(tmpFile1))
	newFile2 := path.Join(dstDir, path.Base(tmpFile2))
	if !FileExists(newFile1) {
		t.Fail()
	}
	if !FileExists(newFile2) {
		t.Fail()
	}
	os.RemoveAll(srcDir)
	os.RemoveAll(dstDir)
}

func TestMoveFilesToDir(t *testing.T) {
	srcDir := path.Join(os.TempDir(), "src")
	dstDir := path.Join(os.TempDir(), "dst")
	if err := os.MkdirAll(srcDir, DirFileMode); err != nil {
		t.Fail()
	}
	if err := os.MkdirAll(dstDir, DirFileMode); err != nil {
		t.Fail()
	}
	tmpFile1 := makeTempFile(srcDir, t)
	tmpFile2 := makeTempFile(srcDir, t)
	err := MoveFilesToDir(srcDir, []string{path.Base(tmpFile1), path.Base(tmpFile2)}, dstDir, false)
	if err != nil {
		t.Fail()
	}
	newFile1 := path.Join(dstDir, path.Base(tmpFile1))
	newFile2 := path.Join(dstDir, path.Base(tmpFile2))
	if !FileExists(newFile1) {
		t.Fail()
	}
	if !FileExists(newFile2) {
		t.Fail()
	}
	if FileExists(tmpFile1) {
		t.Fail()
	}
	if FileExists(tmpFile2) {
		t.Fail()
	}
	os.RemoveAll(srcDir)
	os.RemoveAll(dstDir)
}

func TestExpandUserHome(t *testing.T) {
	user, _ := user.Current()
	home := user.HomeDir
	expanded := ExpandUserHome("~/foo")
	if !strings.HasPrefix(expanded, home) {
		t.Fail()
	}
	if !strings.HasSuffix(expanded, "foo") {
		t.Fail()
	}
	expanded = ExpandUserHome("foo")
	if expanded != "foo" {
		t.Fail()
	}
}

func TestWindowsToUnix(t *testing.T) {
	Assert(PathToUnix("foo"), "foo", t)
	Assert(PathToUnix("foo\\bar"), "foo/bar", t)
	Assert(PathToUnix("\\foo\\bar"), "/foo/bar", t)
	Assert(PathToUnix("C:\\foo\\bar"), "/C/foo/bar", t)
	Assert(PathToUnix("c:\\foo\\bar"), "/c/foo/bar", t)
}

func TestUnixToWindows(t *testing.T) {
	Assert(PathToWindows("foo"), "foo", t)
	Assert(PathToWindows("foo/bar"), "foo\\bar", t)
	Assert(PathToWindows("/foo/bar"), "\\foo\\bar", t)
	Assert(PathToWindows("/C/foo/bar"), "C:\\foo\\bar", t)
	Assert(PathToWindows("/c/foo/bar"), "c:\\foo\\bar", t)
}
