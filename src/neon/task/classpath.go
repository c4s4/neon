package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"strings"
	"os"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"net/http"
	"io"
	"path"
)

var LOCAL_REPOSITORY = util.ExpandUserHome("~/.java/repository")
var DEFAULT_REPOSITORY = "http://central.maven.org/maven2"

func init() {
	build.TaskMap["classpath"] = build.TaskDescriptor{
		Constructor: Classpath,
		Help: `Build a Java classpath.

Arguments:

- classpath: the name of the property to set with classpath.
- classes: a list of class directories to add in classpath (optional).
- jars: a glob or list of globs of jar files to add to classpath (optional).
- dependencies: a list of dependency files to add to classpath (optional).
- scopes: the classpath scope (optional, if set will take dependencies without
  scope and listed scopes, if not set, will only take dependencies without
  scope).
- repositories: a list of repository URLs to get dependencies from (optional,
  defaults to 'http://repo1.maven.org/maven2').
- todir: to copy jar files to given directory (optional).

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
	# copy classpath's jar files to 'build/lib' directory
	- classpath:    _
	  dependencies: 'dependencies.yml'
      todir:        'build/lib'

Notes:

Dependency files should list dependencies as follows:

	- group:    junit
      artifact: junit
      version:  4.12
	  scopes:   [test]

Scopes is optional. If not set, dependency will always be included. If set,
dependency will be included for classpath with these scopes.`,
	}
}

func Classpath(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"classpath", "classes", "jars", "dependencies", "scopes", "repositories", "todir"}
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
	var dependencies []string
	if args.HasField("dependencies") {
		dependencies, err = args.GetListStringsOrString("dependencies")
		if err != nil {
			return nil, fmt.Errorf("argument dependencies of task classpath must be a string or list of strings")
		}
	}
	var scopes []string
	if args.HasField("scopes") {
		scopes, err = args.GetListStringsOrString("scopes")
		if err != nil {
			return nil, fmt.Errorf("argument scopes of task classpath must be a string or list of strings")
		}
	}
	var repositories []string
	if args.HasField("repositories") {
		repositories, err = args.GetListStringsOrString("repositories")
		if err != nil {
			return nil, fmt.Errorf("argument repositories of task classpath must be a string or list of strings")
		}
	}
	var todir string
	if args.HasField("todir") {
		todir, err = args.GetString("todir")
		if err != nil {
			return nil, fmt.Errorf("argument todir of task classpath must be a string")
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
		var _dependencies []string
		for _, _dependency := range dependencies {
			_d, _err := context.EvaluateString(_dependency)
			if _err != nil {
				return fmt.Errorf("evaluating dependencies argument: %v", _err)
			}
			_dependencies = append(_dependencies, _d)
		}
		var _scopes []string
		for _, _scope := range scopes {
			_s, _err := context.EvaluateString(_scope)
			if _err != nil {
				return fmt.Errorf("evaluating scopes argument: %v", _err)
			}
			_scopes = append(_scopes, _s)
		}
		var _repositories []string
		for _, _repository := range repositories {
			_r, _err := context.EvaluateString(_repository)
			if _err != nil {
				return fmt.Errorf("evaluating repositories argument: %v", _err)
			}
			_repositories = append(_repositories, _r)
		}
		_todir, _err := context.EvaluateString(todir)
		if _err != nil {
			return fmt.Errorf("evaluating destination directory: %v", _err)
		}
		// get dependencies
		_deps, _err := getDependencies(_dependencies, _scopes, _repositories, context)
		if _err != nil {
			return fmt.Errorf("getting dependencies: %v", _err)
		}
		// evaluate classpath
		var _elements []string
		_elements = append(_elements, _classes...)
		_elements = append(_elements, _jars...)
		_elements = append(_elements, _deps...)
		_path := strings.Join(_elements, string(os.PathListSeparator))
		context.SetProperty(_classpath, _path)
		// copy jar files to destination directory
		if _todir != "" {
			var _jars []string
			for _, _element := range _elements {
				if strings.HasSuffix(_element, ".jar") {
					_jars = append(_jars, _element)
				}
			}
			_err = copyJarsToDir(_jars, _todir)
			if err != nil {
				return fmt.Errorf("copying jar files to destination directory: %v", _err)
			}
		}
		return nil
	}, nil
}

func getDependencies(dependencies, scopes, repositories []string, context *build.Context) ([]string, error) {
	if !util.DirExists(LOCAL_REPOSITORY) {
		os.MkdirAll(LOCAL_REPOSITORY, util.DIR_FILE_MODE)
	}
	var deps []string
	for _, dependency := range dependencies {
		dep, err := getDependency(dependency, scopes, repositories, context)
		if err != nil {
			return nil, err
		}
		deps = append(deps, dep...)
	}
	return deps, nil
}

func getDependency(file string, scopes, repositories []string, context *build.Context) ([]string, error) {
	var dependencies Dependencies
	source, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(source, &dependencies)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, dependency := range dependencies {
		if selected(scopes, dependency.Scopes) {
			path := dependency.Path(LOCAL_REPOSITORY)
			paths = append(paths, path)
			if !util.FileExists(path) {
				err = downloadDependency(dependency, repositories, context)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	if err != nil {
		return nil, err
	}
	return paths, nil
}

func downloadDependency(dependency Dependency, repositories []string, context *build.Context) error {
	context.Message("Downloading dependency '%s'", dependency.String())
	path := dependency.Path(LOCAL_REPOSITORY)
	dir := filepath.Dir(path)
	if !util.DirExists(dir) {
		os.MkdirAll(dir, util.DIR_FILE_MODE)
	}
	if repositories == nil {
		repositories = []string{DEFAULT_REPOSITORY}
	}
	var err error
	for _, repository := range repositories {
		url := dependency.Path(repository)
		err = download(path, url)
		if err == nil {
			return nil
		}
	}
	return err
}

func download(path, url string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("getting '%s': %v", url, err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return fmt.Errorf("getting '%s': %s", url, response.Status)
	}
	output, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("saving dependency '%s': %v", path, err)
	}
	defer output.Close()
	_, err = io.Copy(output, response.Body)
	if err != nil {
		return fmt.Errorf("saving dependency '%s': %v", path, err)
	}
	return nil
}

type Dependency struct {
	Group    string
	Artifact string
	Version  string
	Scopes   []string
}

func (d *Dependency) Path(root string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s-%s.jar", root, strings.Replace(d.Group, ".", "/", -1), d.Artifact, d.Version, d.Artifact, d.Version)
}

func (d *Dependency) String() string {
	return fmt.Sprintf("%s/%s/%s", d.Group, d.Artifact, d.Version)
}

type Dependencies []Dependency

func selected(classpath, dependency []string) bool {
	if dependency == nil {
		return true
	}
	for _, scope1 := range classpath {
		for _, scope2 := range dependency {
			if scope1 == scope2 {
				return true
			}
		}
	}
	return false
}

func copyJarsToDir(jars []string, dir string) error {
	if !util.DirExists(dir) {
		os.MkdirAll(dir, util.DIR_FILE_MODE)
	}
	for _, jar := range jars {
		dest := path.Join(dir, path.Base(jar))
		err := util.CopyFile(jar, dest)
		if err != nil {
			return err
		}
	}
	return nil
}
