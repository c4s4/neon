package task

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"neon/build"
	"neon/util"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"reflect"
)

const (
	DEFAULT_REPOSITORY = "http://central.maven.org/maven2"
)

var LOCAL_REPOSITORY = util.ExpandUserHome("~/.java/repository")

func init() {
	build.AddTask(build.TaskDesc {
		Name: "classpath",
		Func: Classpath,
		Args: reflect.TypeOf(ClasspathArgs{}),
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
	})
}

type ClasspathArgs struct {
	Classpath    string
	Classes      []string `optional file wrap`
	Jars         []string `optional file wrap`
	Dependencies []string `optional file wrap`
	Scopes       []string `optional wrap`
	Repositories []string `optional wrap`
	Todir        string   `optional file`
}

func Classpath(context *build.Context, args interface{}) error {
	params := args.(ClasspathArgs)
	// get dependencies
	deps, err := getDependencies(params.Dependencies, params.Scopes, params.Repositories, context)
	if err != nil {
		return fmt.Errorf("getting dependencies: %v", err)
	}
	// evaluate classpath
	var elements []string
	elements = append(elements, params.Classes...)
	elements = append(elements, params.Jars...)
	elements = append(elements, deps...)
	path := strings.Join(elements, string(os.PathListSeparator))
	context.SetProperty(params.Classpath, path)
	// copy jar files to destination directory
	if params.Todir != "" {
		var jars []string
		for _, element := range elements {
			if strings.HasSuffix(element, ".jar") {
				jars = append(jars, element)
			}
		}
		err = copyJarsToDir(jars, params.Todir)
		if err != nil {
			return fmt.Errorf("copying jar files to destination directory: %v", err)
		}
	}
	return nil
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
