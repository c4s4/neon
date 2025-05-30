package task

import (
	"fmt"
	"io"
	"net/http"
	"os"
	p "path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/c4s4/neon/neon/build"
	"github.com/c4s4/neon/neon/util"

	"gopkg.in/yaml.v2"
)

const (
	// DefaultRepository is default repository location
	DefaultRepository = "http://central.maven.org/maven2"
)

// LocalRepository is default location for local repository
var LocalRepository = util.ExpandUserHome("~/.java/repository")

func init() {
	build.AddTask(build.TaskDesc{
		Name: "classpath",
		Func: classpath,
		Args: reflect.TypeOf(classpathArgs{}),
		Help: `Build a Java classpath.

Arguments:

- classpath: the property to set with classpath (string).
- classes: class directories to add in classpath (strings, optional, file,
  wrap).
- jars: globs of jar files to add to classpath (strings, optional, file, wrap).
- dependencies: dependency files to add to classpath (strings, optional, file,
  wrap).
- scopes: classpath scope (strings, optional, wrap). If set, will take
  dependencies without scope and listed scopes, if not set, will only take
  dependencies without scope).
- repositories: repository URLs to get dependencies from, defaults to
  'http://repo1.maven.org/maven2' (strings, optional, wrap).
- todir: directory to copy jar files into (string, optional, file).

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

- Dependency files should list dependencies with YAML syntax as follows:

    - group:    junit
      artifact: junit
      version:  4.12
      scopes:   [test]

- Scopes are optional. If not set, dependency will always be included. If set,
  dependency will be included for classpath with these scopes.`,
	})
}

type classpathArgs struct {
	Classpath    string
	Classes      []string `neon:"optional,file,wrap"`
	Jars         []string `neon:"optional,file,wrap"`
	Dependencies []string `neon:"optional,file,wrap"`
	Scopes       []string `neon:"optional,wrap"`
	Repositories []string `neon:"optional,wrap"`
	Todir        string   `neon:"optional,file"`
}

// Classpath builds a Java classpath.
func classpath(context *build.Context, args interface{}) error {
	params := args.(classpathArgs)
	// get dependencies
	var err error
	var jars []string
	if len(params.Jars) > 0 {
		jars, err = util.FindFiles("", params.Jars, []string{}, false)
		if err != nil {
			return fmt.Errorf("getting jars files: %v", err)
		}
	}
	deps, err := getDependencies(params.Dependencies, params.Scopes, params.Repositories, context)
	if err != nil {
		return fmt.Errorf("getting dependencies: %v", err)
	}
	// evaluate classpath
	var elements []string
	elements = append(elements, params.Classes...)
	elements = append(elements, jars...)
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
	if !util.DirExists(LocalRepository) {
		if err := os.MkdirAll(LocalRepository, util.DirFileMode); err != nil {
			return nil, err
		}
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
	var dependencies dependencies
	source, err := os.ReadFile(file)
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
			path := dependency.Path(LocalRepository)
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

func downloadDependency(dependency dependency, repositories []string, context *build.Context) error {
	context.MessageArgs("Downloading dependency '%s'", dependency.String())
	path := dependency.Path(LocalRepository)
	dir := filepath.Dir(path)
	if !util.DirExists(dir) {
		if err := os.MkdirAll(dir, util.DirFileMode); err != nil {
			return err
		}
	}
	if repositories == nil {
		repositories = []string{DefaultRepository}
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
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != 200 {
		return fmt.Errorf("getting '%s': %s", url, response.Status)
	}
	output, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("saving dependency '%s': %v", path, err)
	}
	defer func() {
		_ = output.Close()
	}()
	_, err = io.Copy(output, response.Body)
	if err != nil {
		return fmt.Errorf("saving dependency '%s': %v", path, err)
	}
	return nil
}

type dependency struct {
	Group    string
	Artifact string
	Version  string
	Scopes   []string
}

func (d *dependency) Path(root string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s-%s.jar", root, strings.ReplaceAll(d.Group, ".", "/"), d.Artifact, d.Version, d.Artifact, d.Version)
}

func (d *dependency) String() string {
	return fmt.Sprintf("%s/%s/%s", d.Group, d.Artifact, d.Version)
}

type dependencies []dependency

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
		if err := os.MkdirAll(dir, util.DirFileMode); err != nil {
			return err
		}
	}
	for _, jar := range jars {
		dest := p.Join(dir, p.Base(jar))
		err := util.CopyFile(jar, dest)
		if err != nil {
			return err
		}
	}
	return nil
}
