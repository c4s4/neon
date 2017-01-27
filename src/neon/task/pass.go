package task

import (
	"neon/build"
	"neon/util"
)

func init() {
	build.TaskMap["pass"] = build.TaskDescriptor{
		Constructor: Pass,
		Help: `Does nothing.

Arguments:
- none

Examples:
# do nothing
- pass:

Notes:
- Useful when an instruction is mandatory but we want to do nothing (in a catch
  instruction for instance). For instance, to run a command and don't stop the
  build even if it fails, we could write:
  - try:
	- "command-that-doesnt-exist"
	catch:
	- pass:
- This implementation is super optimized for speed.`,
	}
}

func Pass(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"pass"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	return func() error {
		return nil
	}, nil
}
