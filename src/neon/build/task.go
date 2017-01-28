package build

import (
	"neon/util"
)

type Task func() error

type TaskConstructor func(target *Target, args util.Object) (Task, error)

type TaskDescriptor struct {
	Constructor TaskConstructor
	Help        string
}

var TaskMap map[string]TaskDescriptor = make(map[string]TaskDescriptor)
