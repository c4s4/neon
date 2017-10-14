package build

import (
	"neon/util"
)

// A task is a function that returns an error
type Task func(context *Context) error

// A task constructor is a function that returns a task and an error
type TaskConstructor func(target *Target, args util.Object) (Task, error)

// A task descriptor is made of a task constructor and an help string
type TaskDescriptor struct {
	Constructor TaskConstructor
	Help        string
}

// Map that gives constructor for given task name
var TaskMap map[string]TaskDescriptor = make(map[string]TaskDescriptor)
