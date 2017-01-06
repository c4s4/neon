package main

import (
	"fmt"
	"os"
)

const (
	DEFAULT_FILE_MODE = 0777
)

var tasksMap = map[string]func(target *Target, args interface{}) (Task, error){
	"print": Print,
	"echo":  Print,
	"mkdir": MkDir,
}

type Task func() error

func Print(target *Target, args interface{}) (Task, error) {
	message, ok := args.(string)
	if !ok {
		return nil, fmt.Errorf("argument of task print must be a string")
	}
	return func() error {
		evaluated, err := target.Build.Context.ReplaceProperties(message)
		if err != nil {
			return fmt.Errorf("processing print argument")
		}
		fmt.Println(evaluated)
		return nil
	}, nil
}

func MkDir(target *Target, args interface{}) (Task, error) {
	dir, ok := args.(string)
	if !ok {
		return nil, fmt.Errorf("argument to task mkdir must be s string")
	}
	return func() error {
		evaluated, err := target.Build.Context.ReplaceProperties(dir)
		fmt.Printf("Making directory '%s'\n", evaluated)
		if err != nil {
			return fmt.Errorf("processing mkdir argument")
		}
		err = os.MkdirAll(evaluated, DEFAULT_FILE_MODE)
		if err != nil {
			return fmt.Errorf("making directory '%s': %s", dir, err)
		}
		return nil
	}, nil
}
