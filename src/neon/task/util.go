package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	FILE_MODE     = 0644
	DIR_FILE_MODE = 0755
)

func CheckFields(args util.Object, list, mandatory []string) error {
	task := strings.Join(list, "/")
	fields := args.Fields()
	if err := fieldsList(fields, list); err != nil {
		return fmt.Errorf("building %s task: %v", task, err)
	}
	if err := fieldsMandatory(fields, mandatory); err != nil {
		return fmt.Errorf("building %s task: %v", task, err)
	}
	return nil
}

func fieldsList(fields, list []string) error {
	for _, field := range fields {
		found := false
		for _, e := range list {
			if e == field {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unknown field '%s'", field)
		}
	}
	return nil
}

func fieldsMandatory(fields, mandatory []string) error {
	for _, manda := range mandatory {
		found := false
		for _, field := range fields {
			if manda == field {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("mandatory field '%s' not found", manda)
		}
	}
	return nil
}

func FieldPresent(args util.Object, field string) bool {
	fields := args.Fields()
	for _, f := range fields {
		if f == field {
			return true
		}
	}
	return false
}

func ParseSteps(target *build.Target, object util.Object, field string) ([]build.Step, error) {
	list, err := object.GetList(field)
	if err != nil {
		return nil, err
	}
	var steps []build.Step
	for index, element := range list {
		target.Build.Debug("Parsing step %v in %s field", index, field)
		step, err := build.NewStep(target, element)
		if err != nil {
			return nil, fmt.Errorf("parsing target '%s': %v", target.Name, err)
		}
		steps = append(steps, step)
	}
	return steps, nil
}

func RunSteps(steps []build.Step) error {
	for _, step := range steps {
		err := step.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func sanitizedName(filename string) string {
	if len(filename) > 1 && filename[1] == ':' &&
		runtime.GOOS == "windows" {
		filename = filename[2:]
	}
	filename = filepath.ToSlash(filename)
	filename = strings.TrimLeft(filename, "/.")
	return strings.Replace(filename, "../", "", -1)
}
