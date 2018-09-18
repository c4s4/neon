package build

import (
	"regexp"
	"testing"
	"os"
	"io/ioutil"
)

func TestMessage(t *testing.T) {
	stdout := os.Stdout
	read, write, _ := os.Pipe()
	os.Stdout = write
	Message("This is a test!")
	os.Stdout = stdout
	write.Close()
	out, _ := ioutil.ReadAll(read)
	if string(out) != "This is a test!\n" {
		t.Errorf("Message failure")
	}
}

func TestInfoNotGrey(t *testing.T) {
	stdout := os.Stdout
	read, write, _ := os.Pipe()
	os.Stdout = write
	Info("This is a test!")
	os.Stdout = stdout
	write.Close()
	out, _ := ioutil.ReadAll(read)
	if string(out) != "This is a test!\n" {
		t.Errorf("Message failure: '%s'", string(out))
	}
}

func TestInfoGrey(t *testing.T) {
	Grey = true
	stdout := os.Stdout
	read, write, _ := os.Pipe()
	os.Stdout = write
	Info("This is a test!")
	os.Stdout = stdout
	Grey = false
	write.Close()
	out, _ := ioutil.ReadAll(read)
	if string(out) != "This is a test!\n" {
		t.Errorf("Message failure: '%s'", string(out))
	}
}

func TestTitle(t *testing.T) {
	stdout := os.Stdout
	read, write, _ := os.Pipe()
	os.Stdout = write
	Title("Test")
	os.Stdout = stdout
	write.Close()
	out, _ := ioutil.ReadAll(read)
	if matched, _ := regexp.Match(`-+ Test -+`, out); !matched {
		t.Errorf("Title failure: '%s'", string(out))
	}
}

func TestPrintOk(t *testing.T) {
	stdout := os.Stdout
	read, write, _ := os.Pipe()
	os.Stdout = write
	PrintOk()
	os.Stdout = stdout
	write.Close()
	out, _ := ioutil.ReadAll(read)
	if string(out) != "OK\n" {
		t.Errorf("PrintOk failure")
	}
}

// FIXME
// func TestPrintError(t *testing.T) {
// 	stdout := os.Stdout
// 	read, write, _ := os.Pipe()
// 	os.Stdout = write
// 	PrintError("Test")
// 	os.Stdout = stdout
// 	write.Close()
// 	out, _ := ioutil.ReadAll(read)
// 	if string(out) != "ERROR Test\n" {
// 		t.Errorf("PrintError failure: '%s'", string(out))
// 	}
// }
