package task

import (
	"fmt"
	"neon/util"
	"path/filepath"
)

func init() {
	TasksMap["copy"] = Descriptor{
		Constructor: Copy,
		Help:        "Copy file(s)",
	}
}

func Copy(target *Target, args util.Object) (Task, error) {
	fields := []string{"copy", "dir", "to", "todir", "flat"}
	if err := CheckFields(args, fields, fields[:1]); err != nil {
		return nil, err
	}
	patterns, err := args.GetListStringsOrString("copy")
	if err != nil {
		return nil, fmt.Errorf("argument copy must be a string or list of strings")
	}
	var dir string
	if args.HasField("dir") {
		dir, err = args.GetString("dir")
		if err != nil {
			return nil, fmt.Errorf("argument dir of task copy must be a string", err)
		}
	}
	var to string
	if args.HasField("to") {
		to, err = args.GetString("to")
		if err != nil {
			return nil, fmt.Errorf("argument to of task copy must be a string")
		}
	}
	var toDir string
	if args.HasField("todir") {
		toDir, err = args.GetString("todir")
		if err != nil {
			return nil, fmt.Errorf("argument todir of task copy must be a string")
		}
	}
	flat := true
	if args.HasField("flat") {
		flat, err = args.GetBoolean("flat")
		if err != nil {
			return nil, fmt.Errorf("argument flat of task copy must be a boolean")
		}
	}
	if (to == "" && toDir == "") || (to != "" && toDir != "") {
		return nil, fmt.Errorf("copy task must have one of 'to' or 'toDir' argument")
	}
	return func() error {
		// evaluate arguments
		for index, pattern := range patterns {
			eval, err := target.Build.Context.ReplaceProperties(pattern)
			if err != nil {
				return fmt.Errorf("evaluating pattern: %v", err)
			}
			patterns[index] = eval
		}
		eval, err := target.Build.Context.ReplaceProperties(dir)
		if err != nil {
			return fmt.Errorf("evaluating source directory: %v", err)
		}
		dir = eval
		eval, err = target.Build.Context.ReplaceProperties(to)
		if err != nil {
			return fmt.Errorf("evaluating destination file: %v", err)
		}
		to = eval
		eval, err = target.Build.Context.ReplaceProperties(toDir)
		if err != nil {
			return fmt.Errorf("evaluating destination directory: %v", err)
		}
		toDir = eval
		// find source files
		sources, err := target.Build.Context.FindFiles(dir, patterns)
		if err != nil {
			return fmt.Errorf("getting source files for copy task: %v", err)
		}
		if to != "" && len(sources) > 1 {
			return fmt.Errorf("can't copy more than one file to a given file, use todir instead")
		}
		if len(sources) < 1 {
			return nil
		}
		fmt.Printf("Copying %d file(s)\n", len(sources))
		if to != "" {
			file := filepath.Join(dir, sources[0])
			err = util.CopyFile(file, to)
			if err != nil {
				return fmt.Errorf("copying file: %v", err)
			}
		}
		if toDir != "" {
			err = util.CopyFilesToDir(dir, sources, toDir, flat)
			if err != nil {
				return fmt.Errorf("copying file: %v", err)
			}
		}
		return nil
	}, nil
}
