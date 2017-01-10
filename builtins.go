package main

import (
	zglob "github.com/mattn/go-zglob"
	"os"
)

func Find(dir string, patterns ...string) []string {
	oldDir, err := os.Getwd()
	if err != nil {
		return nil
	}
	err = os.Chdir(dir)
	if err != nil {
		return nil
	}
	var files []string
	for _, pattern := range patterns {
		f, _ := zglob.Glob(pattern)
		for _, e := range f {
			files = append(files, e)
		}
	}
	os.Chdir(oldDir)
	return files
}
