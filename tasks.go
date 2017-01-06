package main

import (
	"fmt"
	"os"
)

const (
	DEFAULT_FILE_MODE = 0777
)

var tasksMap = map[string]func(target *Target, args interface{}) (Task, error){
	"mkdir": MkDir,
}

type Task func(target *Target, args interface{}) error

func MkDir(target *Target, args interface{}) (Task, error) {
	dir, ok := args.(string)
	if !ok {
		return nil, fmt.Errorf("Argument to task mkdir must be s string")
	}
	return func(target *Target, args interface{}) error {
		err := os.MkdirAll(dir, DEFAULT_FILE_MODE)
		if err != nil {
			return fmt.Errorf("making directory '%s': %s", dir, err)
		}
		return nil
	}, nil
}
