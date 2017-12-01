package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"strings"
	"os"
)

func init() {
	build.TaskMap["classpath"] = build.TaskDescriptor{
		Constructor: Classpath,
		Help: `Build a Java classpath.

Arguments:

- classpath: the name of the property to set with classpath.
- classes: a list of class directories to add in classpath (optional).
- jars: a glob or list of globs of jar files to add to classpath (optional).
# TODO
- dependencies: a list of dependency files to add to classpath (optional).
- repositories: a list of repository URLs to get dependencies from (optional,
  defaults to 'http://repo1.maven.org/maven2').
- scope: the classpath scope (optional, defaults to 'runtime').

Examples:

	# build classpath with classes in build/classes directory
	- classpath: 'classpath'
	  classes:   'build/classes'
    # build classpath with jar files in lib directory
    - classpath: 'classpath'
      jars:      'lib/*.jar'
	# build classpath with a dependencies file
	- classpath:    'classpath'
	  dependencies: 'dependencies.yml'

Notes:

Dependency files should list dependencies as follows:

	- group:    junit
      artifact: junit
      version:  4.12
	  scope:    test

Scopes may be runtime (default), compile, test or provided.`,
	}
}

func Classpath(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"classpath", "classes", "jars"}
	if err := CheckFields(args, fields, fields[:1]); err != nil {
		return nil, err
	}
	classpath, err := args.GetString("classpath")
	if err != nil {
		return nil, fmt.Errorf("argument classpath must be a string")
	}
	var classes []string
	if args.HasField("classes") {
		classes, err = args.GetListStringsOrString("classes")
		if err != nil {
			return nil, fmt.Errorf("argument classes of task classpath must be a string or list of strings")
		}
	}
	var jars []string
	if args.HasField("jars") {
		jars, err = args.GetListStringsOrString("jars")
		if err != nil {
			return nil, fmt.Errorf("argument jars of task classpath must be a string or list of strings")
		}
	}
	return func(context *build.Context) error {
		// evaluate arguments
		_classpath, _err := context.EvaluateString(classpath)
		if _err != nil {
			return fmt.Errorf("evaluating classpath argument: %v", _err)
		}
		var _classes []string
		for _, _class := range classes {
			_c, _err := context.EvaluateString(_class)
			if _err != nil {
				return fmt.Errorf("evaluating classes argument: %v", _err)
			}
			_classes = append(_classes, _c)
		}
		_jars, _err := context.FindFiles(".", jars, nil, false)
		if _err != nil {
			return fmt.Errorf("finding jar files: %v", _err)
		}
		// evaluate classpath
		_elements := append(_classes, _jars...)
		_path := strings.Join(_elements, string(os.PathListSeparator))
		context.SetProperty(_classpath, _path)
		return nil
	}, nil
}
