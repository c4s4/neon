package build

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"neon/util"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	// DefaultRepo is local repository root directory
	DefaultRepo = "~/.neon"
	// RegexpPlugin is regexp for a plugin name
	RegexpPlugin = `[\w-]+/[\w-]+`
)

// Fields is the list of possible root fields for a build file
var Fields = []string{"doc", "default", "extends", "repository", "context", "singleton",
	"shell", "properties", "configuration", "environment", "targets", "version"}

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
	Version     string
}

// NewBuild makes a build from a build file
// - file: path of the build file
// - base: base of the build
// - repo: repository location
// Return:
// - Pointer to the build
// - error if something went wrong
func NewBuild(file, base, repo string) (*Build, error) {
	object, build, err := parseBuildFile(file)
	if err != nil {
		return nil, err
	}
	if err := setDirectories(build, base); err != nil {
		return nil, err
	}
	if err := object.CheckFields(Fields); err != nil {
		return nil, fmt.Errorf("parsing build file: %v", err)
	}
	if err := ParseFields(object, build, repo); err != nil {
		return nil, err
	}
	build.Properties = build.GetProperties()
	build.Environment = build.GetEnvironment()
	build.SetDir(build.Dir)
	return build, nil
}

// ParseFields parses build file fields:
// - object: the build as an object.
// - build: the build object.
// - repo: the repository.
// Return: an error if something went wrong.
func ParseFields(object util.Object, build *Build, repo string) error {
	if err := ParseSingleton(object, build); err != nil {
		return err
	}
	if err := ParseShell(object, build); err != nil {
		return err
	}
	if err := ParseDefault(object, build); err != nil {
		return err
	}
	if err := ParseDoc(object, build); err != nil {
		return err
	}
	if err := ParseRepository(object, build, repo); err != nil {
		return err
	}
	if err := ParseContext(object, build); err != nil {
		return err
	}
	if err := ParseExtends(object, build); err != nil {
		return err
	}
	if err := ParseProperties(object, build); err != nil {
		return err
	}
	if err := ParseConfiguration(object, build); err != nil {
		return err
	}
	if err := ParseEnvironment(object, build); err != nil {
		return err
	}
	if err := ParseTargets(object, build); err != nil {
		return err
	}
	return ParseVersion(object, build)
}

func parseBuildFile(file string) (util.Object, *Build, error) {
	build := &Build{}
	file = util.ExpandUserHome(file)
	build.File = filepath.Base(file)
	source, err := util.ReadFile(file)
	if err != nil {
		return nil, nil, fmt.Errorf("loading build file '%s': %v", file, err)
	}
	var object util.Object
	if err = yaml.Unmarshal(source, &object); err != nil {
		return nil, nil, fmt.Errorf("build must be a map with string keys: %v", err)
	}
	return object, build, nil
}

func setDirectories(build *Build, base string) error {
	build.Dir = base
	here, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting build file directory: %v", err)
	}
	build.Here = here
	return nil
}

// GetProperties returns build properties, including those inherited from
// parents
// Return: build properties as an Object
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

// GetEnvironment returns the build environment, including the environment
// inherited from parents
// Return: environment as a map with string keys and values
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

// GetTargets returns build targets, including those inherited from parents
// Return: targets as a map of targets with their name as keys
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

// GetTargetByName return target with given name. If not defined in build,
// return target inherited from parent
// - name: the target name as a string
// Return: found target
func (build *Build) GetTargetByName(name string) *Target {
	target, found := build.Targets[name]
	if found {
		return target
	}
	for _, parent := range build.Parents {
		target = parent.GetTargetByName(name)
		if target != nil {
			return target
		}
	}
	return nil
}

// SetDir sets the build directory, propagating to parents
// - dir: build directory as a string
func (build *Build) SetDir(dir string) {
	build.Dir = dir
	for _, parent := range build.Parents {
		parent.SetDir(dir)
	}
}

// SetCommandLineProperties defines properties passed on command line in the
// context. These properties overwrite those define in the build file.
// - props: properties as a YAML map
// Return: error if something went wrong
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

// GetDefault returns default targets. If none is defined in build, return
// those from parent build files.
// Return: default targets a slice of strings
func (build *Build) GetDefault() []string {
	if len(build.Default) > 0 {
		return build.Default
	}
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
	return build.Default
}

// GetScripts return a list of context scripts to run.
// Return: the list of context scripts
func (build *Build) GetScripts() []string {
	var scripts []string
	for _, parent := range build.Parents {
		scripts = append(scripts, parent.GetScripts()...)
	}
	scripts = append(scripts, build.Scripts...)
	return scripts
}

// Run runs given targets in a build context. If no target is given, runs
// default one.
// - context: the context to run into
// - targets: targets to run as a slice of strings
// Return: error if something went wrong
func (build *Build) Run(context *Context, targets []string) error {
	if err := build.CheckVersion(context); err != nil {
		return err
	}
	if err := build.EnsureSingle(context); err != nil {
		return err
	}
	if len(targets) == 0 {
		targets = build.GetDefault()
		if len(targets) == 0 {
			return fmt.Errorf("no default target")
		}
	}
	for _, target := range targets {
		context.Stack = NewStack()
		err := build.RunTarget(context, target)
		if err != nil {
			return err
		}
	}
	return nil
}

// RunTarget runs given target in a build context.
// - context: build context
// - name: name of the target to run as a string
// Return: an error if something went wrong
func (build *Build) RunTarget(context *Context, name string) error {
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

// RunParentTarget runs parent target with given name in a build context.
// - context: build context
// - name: the name of the target to run
// Return:
// - boolean: that tells if parent target was found
// - error: if something went wrong
func (build *Build) RunParentTarget(context *Context, name string) (bool, error) {
	for _, parent := range build.Parents {
		target := parent.GetTargetByName(name)
		if target != nil {
			err := context.Stack.Push(target)
			if err != nil {
				return false, err
			}
			err = target.Steps.Run(context)
			if err != nil {
				return true, fmt.Errorf("running target '%s': %v", name, err)
			}
			return true, nil
		}
	}
	return false, nil
}

// GetShell return shell for current os.
// Return:
// - shell as a slice of strings (such as ["sh", "-c"])
// - error if something went wrong
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

// EnsureSingle runs a TCP server on given port to ensure that a single
// instance is running on a machine. Fails if another instance is already
// running on same port.
// - context: build context
// Return: an error if another instance is running on same port
func (build *Build) EnsureSingle(context *Context) error {
	if build.Singleton == "" {
		return nil
	}
	expression := build.Singleton
	if IsExpression(expression) {
		expression = expression[1:]
	}
	singleton, err := context.EvaluateExpression(expression)
	if err != nil {
		return fmt.Errorf("evaluating singleton port expression '%s': %v", expression, err)
	}
	port, ok := singleton.(int64)
	if !ok {
		portInt, ok := singleton.(int)
		if !ok {
			return fmt.Errorf("singleton port expression '%s' must return an integer", expression)
		}
		port = int64(portInt)
	}
	if port < 0 || port > 65535 {
		return fmt.Errorf("singleton port port must be between 0 and 65535")
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("listening singleton port: %v", err)
	}
	go func() {
		for {
			listener.Accept()
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return nil
}

// CheckVersion checks evaluates version expression to check that NeON version is OK
func (build *Build) CheckVersion(context *Context) error {
	if build.Version == "" {
		return nil
	}
	result, err := context.EvaluateExpression(build.Version)
	if err != nil {
		return fmt.Errorf("evaluating version expression: %v", err)
	}
	versionOK, ok := result.(bool)
	if !ok {
		return fmt.Errorf("version expression should return a boolean")
	}
	if !versionOK {
		return fmt.Errorf("neon version '%s' doesn't meet requirements in version field", NeonVersion)
	}
	return nil
}
