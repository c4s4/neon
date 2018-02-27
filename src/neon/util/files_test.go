package util

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

const (
	FileMode = 0644
)

func tempFile() string {
	tempFile, _ := ioutil.TempFile("", "files_test.tmp")
	return tempFile.Name()
}

func TestReadFile(t *testing.T) {
	tempFile := tempFile()
	defer os.Remove(tempFile)
	ioutil.WriteFile(tempFile, []byte("test"), FileMode)
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
	tempFile := tempFile()
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
	srcFile := tempFile()
	defer os.Remove(srcFile)
	dstFile := path.Join(os.TempDir(), "test.tmp")
	defer os.Remove(dstFile)
	err := ioutil.WriteFile(srcFile, []byte("test"), FileMode)
	if err != nil {
		t.Fail()
	}
	err = CopyFile(srcFile, dstFile)
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
