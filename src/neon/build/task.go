package build

import (
	"fmt"
	zglob "github.com/mattn/go-zglob"
	"io"
	"io/ioutil"
	"neon/util"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"
)

const (
	FILE_MODE     = 0644
	DIR_FILE_MODE = 0755
)

type Task func() error

type Constructor func(target *Target, args util.Object) (Task, error)

var tasksMap map[string]Constructor

func init() {
	tasksMap = map[string]Constructor{
		"script": Script,
		"print":  Print,
		"chdir":  Chdir,
		"mkdir":  MkDir,
		"touch":  Touch,
		"link":   Link,
		"copy":   Copy,
		"remove": Remove,
		"delete": Delete,
		"if":     If,
		"for":    For,
		"while":  While,
		"try":    Try,
		"pass":   Pass,
	}
}

// TASKS DEFINITIONS

func Script(target *Target, args util.Object) (Task, error) {
	fields := []string{"script"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	source, err := args.GetString("script")
	if err != nil {
		return nil, fmt.Errorf("parsing script task: %v", err)
	}
	return func() error {
		_, err := target.Build.Context.Evaluate(source)
		if err != nil {
			return fmt.Errorf("evaluating script: %v", err)
		}
		return nil
	}, nil
}

func Print(target *Target, args util.Object) (Task, error) {
	fields := []string{"print"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	message, ok := args["print"].(string)
	if !ok {
		return nil, fmt.Errorf("argument of task print must be a string")
	}
	return func() error {
		evaluated, err := target.Build.Context.ReplaceProperties(message)
		if err != nil {
			return fmt.Errorf("processing print argument: %v", err)
		}
		fmt.Println(evaluated)
		return nil
	}, nil
}

func Chdir(target *Target, args util.Object) (Task, error) {
	fields := []string{"chdir"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	dir, ok := args["chdir"].(string)
	if !ok {
		return nil, fmt.Errorf("argument to task chdir must be a string")
	}
	return func() error {
		directory, err := target.Build.Context.ReplaceProperties(dir)
		fmt.Printf("Changing to directory '%s'\n", directory)
		if err != nil {
			return fmt.Errorf("processing chdir argument: %v", err)
		}
		err = os.Chdir(directory)
		if err != nil {
			return fmt.Errorf("changing to directory '%s': %s", directory, err)
		}
		return nil
	}, nil
}

func MkDir(target *Target, args util.Object) (Task, error) {
	fields := []string{"mkdir"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	dir, ok := args["mkdir"].(string)
	if !ok {
		return nil, fmt.Errorf("argument to task mkdir must be a string")
	}
	return func() error {
		directory, err := target.Build.Context.ReplaceProperties(dir)
		if err != nil {
			return fmt.Errorf("processing mkdir argument: %v", err)
		}
		fmt.Printf("Making directory '%s'\n", directory)
		err = os.MkdirAll(directory, DIR_FILE_MODE)
		if err != nil {
			return fmt.Errorf("making directory '%s': %s", directory, err)
		}
		return nil
	}, nil
}

func Touch(target *Target, args util.Object) (Task, error) {
	fields := []string{"touch"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	files, err := args.GetListStringsOrString("touch")
	if err != nil {
		return nil, fmt.Errorf("argument to task touch must be a string or list of strings")
	}
	return func() error {
		fmt.Printf("Touching %d file(s)\n", len(files))
		for _, file := range files {
			path, err := target.Build.Context.ReplaceProperties(file)
			if err != nil {
				return fmt.Errorf("processing touch argument: %v", err)
			}
			if util.FileExists(path) {
				time := time.Now()
				err = os.Chtimes(path, time, time)
				if err != nil {
					return fmt.Errorf("changing times of file '%s': %v", path, err)
				}
			} else {
				err := ioutil.WriteFile(path, []byte{}, FILE_MODE)
				if err != nil {
					return fmt.Errorf("creating file '%s': %v", path, err)
				}
			}
		}
		return nil
	}, nil
}

func Link(target *Target, args util.Object) (Task, error) {
	fields := []string{"link", "to"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	s, err := args.GetString("link")
	if err != nil {
		return nil, fmt.Errorf("argument link must be a string")
	}
	d, err := args.GetString("to")
	if err != nil {
		return nil, fmt.Errorf("argument to of task link must be a string")
	}
	return func() error {
		source, err := target.Build.Context.ReplaceProperties(s)
		if err != nil {
			return fmt.Errorf("processing link argument: %v", err)
		}
		dest, err := target.Build.Context.ReplaceProperties(d)
		if err != nil {
			return fmt.Errorf("processing to argument of link task: %v", err)
		}
		fmt.Printf("Linking file '%s' to '%s'\n", source, dest)
		err = os.Symlink(source, dest)
		if err != nil {
			return fmt.Errorf("linking files: %v", err)
		}
		return nil
	}, nil
}

func Copy(target *Target, args util.Object) (Task, error) {
	fields := []string{"copy", "to"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	s, err := args.GetString("copy")
	if err != nil {
		return nil, fmt.Errorf("argument copy must be a string")
	}
	d, err := args.GetString("to")
	if err != nil {
		return nil, fmt.Errorf("argument to of task copy must be a string")
	}
	return func() error {
		source, err := target.Build.Context.ReplaceProperties(s)
		if err != nil {
			return fmt.Errorf("processing copy argument: %v", err)
		}
		dest, err := target.Build.Context.ReplaceProperties(d)
		if err != nil {
			return fmt.Errorf("processing to argument of copy task: %v", err)
		}
		fmt.Printf("Copying file '%s' to '%s'\n", source, dest)
		err = CopyFile(source, dest)
		if err != nil {
			return fmt.Errorf("copying files: %v", err)
		}
		return nil
	}, nil
}

func Remove(target *Target, args util.Object) (Task, error) {
	fields := []string{"remove"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	patterns, err := args.GetListStringsOrString("remove")
	if err != nil {
		return nil, fmt.Errorf("remove argument must a string or list of strings")
	}
	return func() error {
		var files []string
		for _, patt := range patterns {
			pattern, err := target.Build.Context.ReplaceProperties(patt)
			if err != nil {
				return fmt.Errorf("evaluating pattern in task remove: %v", err)
			}
			list, _ := zglob.Glob(pattern)
			for _, file := range list {
				files = append(files, file)
			}
		}
		sort.Strings(files)
		fmt.Printf("Removing %d file(s)\n", len(files))
		for _, file := range files {
			if err = os.Remove(file); err != nil {
				return fmt.Errorf("removing file '%s': %v", file, err)
			}
		}
		return nil
	}, nil
}

func Delete(target *Target, args util.Object) (Task, error) {
	fields := []string{"delete"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	directories, err := args.GetListStringsOrString("delete")
	if err != nil {
		return nil, fmt.Errorf("delete argument must be string or list of strings")
	}
	return func() error {
		for _, dir := range directories {
			directory, err := target.Build.Context.ReplaceProperties(dir)
			if err != nil {
				return fmt.Errorf("evaluating directory in task delete: %v", err)
			}
			if _, err := os.Stat(directory); err == nil {
				fmt.Printf("Deleting directory '%s'\n", directory)
				err = os.RemoveAll(directory)
				if err != nil {
					return fmt.Errorf("deleting directory '%s': %v", directory, err)
				}
			}
		}
		return nil
	}, nil
}

func If(target *Target, args util.Object) (Task, error) {
	fields := []string{"if", "then", "else"}
	if err := CheckFields(args, fields, fields[:2]); err != nil {
		return nil, err
	}
	condition, err := args.GetString("if")
	if err != nil {
		return nil, fmt.Errorf("evaluating if construct: %v", err)
	}
	thenSteps, err := ParseSteps(target, args, "then")
	if err != nil {
		return nil, err
	}
	var elseSteps []Step
	if FieldPresent(args, "else") {
		elseSteps, err = ParseSteps(target, args, "else")
		if err != nil {
			return nil, err
		}
	}
	return func() error {
		result, err := target.Build.Context.Evaluate(condition)
		if err != nil {
			return fmt.Errorf("evaluating 'if' condition: %v", err)
		}
		boolean, ok := result.(bool)
		if !ok {
			return fmt.Errorf("evaluating if condition: must return a bool")
		}
		if boolean {
			err := RunSteps(thenSteps)
			if err != nil {
				return err
			}
		} else {
			err := RunSteps(elseSteps)
			if err != nil {
				return err
			}
		}
		return nil
	}, nil
}

func For(target *Target, args util.Object) (Task, error) {
	fields := []string{"for", "in", "do"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	variable, err := args.GetString("for")
	if err != nil {
		return nil, fmt.Errorf("'for' field of a 'for' loop must be a string")
	}
	list, err := args.GetList("in")
	expression := ""
	if err != nil {
		expression, err = args.GetString("in")
		if err != nil {
			return nil, fmt.Errorf("'in' field of 'for' loop must be a list or string")
		}
	}
	steps, err := ParseSteps(target, args, "do")
	if err != nil {
		return nil, err
	}
	return func() error {
		if expression != "" {
			result, err := target.Build.Context.Evaluate(expression)
			if err != nil {
				return fmt.Errorf("evaluating in field of for loop: %v", err)
			}
			list, err = ToList(result)
			if err != nil {
				return fmt.Errorf("'in' field of 'for' loop must be an expression that returns a list")
			}
		}
		for _, value := range list {
			target.Build.Context.SetProperty(variable, value)
			if err != nil {
				return err
			}
			err := RunSteps(steps)
			if err != nil {
				return err
			}
		}
		return nil
	}, nil

}

func While(target *Target, args util.Object) (Task, error) {
	fields := []string{"while", "do"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	condition, err := args.GetString("while")
	if err != nil {
		return nil, fmt.Errorf("'while' field of a 'while' loop must be a string")
	}
	steps, err := ParseSteps(target, args, "do")
	if err != nil {
		return nil, err
	}
	return func() error {
		for {
			result, err := target.Build.Context.Evaluate(condition)
			if err != nil {
				return fmt.Errorf("evaluating 'while' field of 'while' loop: %v", err)
			}
			loop, ok := result.(bool)
			if !ok {
				return fmt.Errorf("evaluating 'while' condition: must return a bool")
			}
			if !loop {
				break
			}
			err = RunSteps(steps)
			if err != nil {
				return err
			}
		}
		return nil
	}, nil
}

func Try(target *Target, args util.Object) (Task, error) {
	fields := []string{"try", "catch"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	trySteps, err := ParseSteps(target, args, "try")
	if err != nil {
		return nil, err
	}
	catchSteps, err := ParseSteps(target, args, "catch")
	if err != nil {
		return nil, err
	}
	return func() error {
		err := RunSteps(trySteps)
		if err != nil {
			target.Build.Context.SetProperty("error", err.Error())
			err = RunSteps(catchSteps)
			if err != nil {
				return err
			}
		}
		return nil
	}, nil
}

func Pass(target *Target, args util.Object) (Task, error) {
	fields := []string{"pass"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	return func() error {
		return nil
	}, nil
}

// UTILITY FUNCTIONS

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

func ParseSteps(target *Target, object util.Object, field string) ([]Step, error) {
	list, err := object.GetList(field)
	if err != nil {
		return nil, err
	}
	var steps []Step
	for index, element := range list {
		target.Build.Log("Parsing step %v in %s field", index, field)
		step, err := NewStep(target, element)
		if err != nil {
			return nil, fmt.Errorf("parsing target '%s': %v", target.Name, err)
		}
		steps = append(steps, step)
	}
	return steps, nil
}

func RunSteps(steps []Step) error {
	for _, step := range steps {
		err := step.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func ToList(object interface{}) ([]interface{}, error) {
	slice := reflect.ValueOf(object)
	if slice.Kind() == reflect.Slice {
		result := make([]interface{}, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			result[i] = slice.Index(i).Interface()
		}
		return result, nil
	} else {
		return nil, fmt.Errorf("must be a list")
	}
}

func CopyFile(source, dest string) error {
	from, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("opening source file '%s': %v", source, err)
	}
	defer from.Close()
	to, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("creating desctination file '%s': %v", dest, err)
	}
	defer to.Close()
	_, err = io.Copy(from, to)
	if err != nil {
		return fmt.Errorf("copying file: %v", err)
	}
	err = to.Sync()
	if err != nil {
		return fmt.Errorf("syncing destination file: %v", err)
	}
	return nil
}
