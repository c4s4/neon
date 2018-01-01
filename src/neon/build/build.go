package build

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"neon/util"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"net"
	"time"
)

const (
	// location of the repository root
	DEFAULT_REPO = "~/.neon"
	// regexp for plugin name
	RE_PLUGIN = `[\w-]+/[\w-]+`
)

// Possible fields for a build file
var FIELDS = []string{"doc", "default", "extends", "repository", "context",
	"singleton", "shell", "properties", "configuration", "environment", "targets"}

// Build structure
type Build struct {
	File        string
	Dir         string
	Here        string
	Default     []string
	Doc         string
	Repository  string
	Singleton   string
	Shell       map[string][]string
	Scripts     []string
	Extends     []string
	Config      []string
	Properties  util.Object
	Environment map[string]string
	Targets     map[string]*Target
	Parents     []*Build
}

// Make a build from a build file
func NewBuild(file string) (*Build, error) {
	build := &Build{}
	path := util.ExpandUserHome(file)
	build.File = filepath.Base(path)
	base, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return nil, fmt.Errorf("getting build file directory: %v", err)
	}
	build.Dir = base
	here, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting current directory: %v", err)
	}
	build.Here = here
	source, err := util.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("loading build file '%s': %v", path, err)
	}
	var object util.Object
	err = yaml.Unmarshal(source, &object)
	if err != nil {
		return nil, fmt.Errorf("build must be a map with string keys")
	}
	if err := object.CheckFields(FIELDS); err != nil {
		return nil, fmt.Errorf("parsing build file: %v", err)
	}
	if err := ParseSingleton(object, build); err != nil {
		return nil, err
	}
	if err := ParseShell(object, build); err != nil {
		return nil, err
	}
	if err := ParseDefault(object, build); err != nil {
		return nil, err
	}
	if err := ParseDoc(object, build); err != nil {
		return nil, err
	}
	if err := ParseRepository(object, build); err != nil {
		return nil, err
	}
	if err := ParseContext(object, build); err != nil {
		return nil, err
	}
	if err := ParseExtends(object, build); err != nil {
		return nil, err
	}
	if err := ParseProperties(object, build); err != nil {
		return nil, err
	}
	if err := ParseConfiguration(object, build); err != nil {
		return nil, err
	}
	if err := ParseEnvironment(object, build); err != nil {
		return nil, err
	}
	if err := ParseTargets(object, build); err != nil {
		return nil, err
	}
	build.SetDir(base)
	return build, nil
}

// Return the build properties, including those inherited from parents
func (build *Build) GetProperties() util.Object {
	var properties = make(map[string]interface{})
	for _, parent := range build.Parents {
		for name, value := range parent.GetProperties() {
			properties[name] = value
		}
	}
	for name, value := range build.Properties {
		properties[name] = value
	}
	return properties
}

// Return the build environment, including those inherited from parents
func (build *Build) GetEnvironment() map[string]string {
	var environment = make(map[string]string)
	for _, parent := range build.Parents {
		for name, value := range parent.GetEnvironment() {
			environment[name] = value
		}
	}
	for name, value := range build.Environment {
		environment[name] = value
	}
	return environment
}

// Return the build targets, including those inherited from parents
func (build *Build) GetTargets() map[string]*Target {
	var targets = make(map[string]*Target)
	for _, parent := range build.Parents {
		for name, target := range parent.GetTargets() {
			targets[name] = target
		}
	}
	for name, target := range build.Targets {
		targets[name] = target
	}
	return targets
}

// Return target with given name. If not defined in build, return target
// inherited from parent
func (build *Build) GetTargetByName(name string) *Target {
	target, found := build.Targets[name]
	if found {
		return target
	} else {
		for _, parent := range build.Parents {
			target = parent.GetTargetByName(name)
			if target != nil {
				return target
			}
		}
	}
	return nil
}

// Set the build directory, propagating to parents
func (build *Build) SetDir(dir string) {
	build.Dir = dir
	for _, parent := range build.Parents {
		parent.SetDir(dir)
	}
}

// Set command line properties, that overwrite build ones
func (build *Build) SetCommandLineProperties(props string) error {
	var object util.Object
	err := yaml.Unmarshal([]byte(props), &object)
	if err != nil {
		return fmt.Errorf("parsing command line properties: properties must be a map with string keys")
	}
	for name, value := range object {
		build.Properties[name] = value
	}
	return nil
}

// Return default targets. If none is defined in build, return those from
// parents
func (build *Build) GetDefault() []string {
	if len(build.Default) > 0 {
		return build.Default
	} else {
		for _, parent := range build.Parents {
			if len(parent.Default) > 0 {
				return parent.Default
			}
		}
		for _, parent := range build.Parents {
			parentDefault := parent.GetDefault()
			if len(parentDefault) > 0 {
				return parentDefault
			}
		}
	}
	return build.Default
}

// Run build given targets. If no target is given, run default one.
func (build *Build) Run(context *Context, targets []string) error {
	if err := build.Listen(context); err != nil {
		return err
	}
	if len(targets) == 0 {
		targets = build.GetDefault()
		if len(targets) == 0 {
			return fmt.Errorf("no default target")
		}
	}
	for _, target := range targets {
		err := build.RunTarget(target, context)
		if err != nil {
			return err
		}
	}
	return nil
}

// Run given target
func (build *Build) RunTarget(name string, context *Context) error {
	target := build.GetTargetByName(name)
	if target == nil {
		return fmt.Errorf("target '%s' not found", name)
	}
	err := target.Run(context)
	if err != nil {
		return fmt.Errorf("running target '%s': %v", name, err)
	}
	return nil
}

// Run parent target
func (build *Build) RunParentTarget(name string, context *Context) (bool, error) {
	for _, parent := range build.Parents {
		target := parent.GetTargetByName(name)
		if target != nil {
			err := target.RunSteps(context)
			if err != nil {
				return true, fmt.Errorf("running target '%s': %v", name, err)
			}
			return true, nil
		} else {
			ok, err := parent.RunParentTarget(name, context)
			if ok || err != nil {
				return ok, err
			}
		}
	}
	return false, nil
}

// Get parent build file path
func (build *Build) PluginPath(name string) string {
	if path.IsAbs(name) {
		return name
	} else if strings.HasPrefix(name, "./") {
		return filepath.Join(build.Dir, name)
	} else {
		repo := util.ExpandAndJoinToRoot(build.Dir, build.Repository)
		return filepath.Join(repo, name)
	}
}

// Get plugin name for given resource
func (build *Build) PluginName(name string) string {
	re := regexp.MustCompile(`^(` + RE_PLUGIN + `)/.+$`)
	if re.MatchString(name) {
		return re.FindStringSubmatch(name)[1]
	} else {
		return ""
	}
}

// GetShell return shell for current os.
func (build *Build) GetShell() ([]string, error) {
	for system, shell := range build.Shell {
		if system != "default" && system == runtime.GOOS {
			return shell, nil
		}
	}
	shell, ok := build.Shell["default"]
	if !ok {
		return nil, fmt.Errorf("no shell found for '%s'", runtime.GOOS)
	}
	return shell, nil
}

// Run a TCP server on given port to ensure that a single instance is running
// on a machine. Fails if another instance is already running on same port.
func (build *Build) Listen(context *Context) error {
	if build.Singleton == "" {
		return nil
	}
	singleton, err := context.EvaluateExpression(build.Singleton)
	if err != nil {
		return fmt.Errorf("evaluating singleton port expression '%s': %v", build.Singleton, err)
	}
	port, ok := singleton.(int)
	if !ok {
		return fmt.Errorf("singleton port expression '%s' must return an integer", build.Singleton)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("another instance of the build is already running")
	}
	go func() {
		for {
			listener.Accept()
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return nil
}
