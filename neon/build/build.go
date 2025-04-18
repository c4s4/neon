package build

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/c4s4/neon/neon/util"
	"gopkg.in/yaml.v2"
)

const (
	// DefaultRepository is local repository root directory
	DefaultRepository = "~/.neon"
	// RegexpPlugin is regexp for a plugin name
	RegexpPlugin = `[\w-]+/[\w-]+`
)

// Fields is the list of possible root fields for a build file
var Fields = []string{"doc", "default", "extends", "repository", "context", "singleton",
	"shell", "properties", "configuration", "expose", "environment", "dotenv", "targets", "version"}

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
	Expose      []string
	Properties  util.Object
	Environment map[string]string
	DotEnv      []string
	Targets     map[string]*Target
	Parents     []*Build
	Root        *Build
	Version     string
	Template    bool
}

// NewBuild makes a build from a build file
// - file: path of the build file
// - base: base of the build
// - repo: repository location
// Return:
// - Pointer to the build
// - error if something went wrong
func NewBuild(file, base, repo string, template bool) (*Build, error) {
	object, build, err := parseBuildFile(file)
	if err != nil {
		return nil, err
	}
	if err := SetDirectories(build, base); err != nil {
		return nil, err
	}
	if err := object.CheckFields(Fields); err != nil {
		return nil, fmt.Errorf("parsing build file: %v", err)
	}
	if err := ParseFields(object, build, repo); err != nil {
		return nil, err
	}
	build.Parents, err = build.GetParents()
	if err != nil {
		return nil, err
	}
	build.Properties = build.GetProperties()
	build.Environment = build.GetEnvironment()
	build.DotEnv = build.GetDotEnv()
	build.SetDir(build.Dir)
	build.SetRoot(build)
	build.Template = template
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
	if err := ParseExpose(object, build); err != nil {
		return err
	}
	if err := ParseEnvironment(object, build); err != nil {
		return err
	}
	if err := ParseDotEnv(object, build); err != nil {
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

// SetDirectories sets build and base directories:
// - build: the build directory.
// - base: the base directory.
// Return: an error if something went wrong.
func SetDirectories(build *Build, base string) error {
	build.Dir = base
	here, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting build file directory: %v", err)
	}
	build.Here = here
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

// SetRoot sets the root build, propagating to parents
// - build: root build
func (build *Build) SetRoot(root *Build) {
	build.Root = root
	for _, parent := range build.Parents {
		parent.SetRoot(root)
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

// GetParents returns parent build objects.
// Return list of build objects and an error if any.
func (build *Build) GetParents() ([]*Build, error) {
	var parents []*Build
	for _, extend := range build.Extends {
		file, err := build.ParentPath(extend)
		if err != nil {
			return nil, fmt.Errorf("searching parent build file '%s': %v", extend, err)
		}
		parent, err := NewBuild(file, filepath.Dir(file), build.Repository, build.Template)
		if err != nil {
			return nil, fmt.Errorf("loading parent build file '%s': %v", extend, err)
		}
		parents = append(parents, parent)
	}
	return parents, nil
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

// GetDotEnv returns the list of dotenv files to load in environment, including
// those inherited from parents
// Return: list of dotenv files as a slice of strings
func (build *Build) GetDotEnv() []string {
	var dotenv []string
	for _, parent := range build.Parents {
		dotenv = append(dotenv, parent.GetDotEnv()...)
	}
	dotenv = append(dotenv, build.DotEnv...)
	return dotenv
}

// GetTargets returns build targets, including those inherited from parents
// Return: targets as a map of targets with their name as keys
func (build *Build) GetTargets() map[string]*Target {
	var targets = make(map[string]*Target)
	for i := len(build.Parents) - 1; i >= 0; i-- {
		parent := build.Parents[i]
		for name, target := range parent.GetTargets() {
			targets[name] = target
		}
	}
	for name, target := range build.Targets {
		targets[name] = target
	}
	return targets
}

// GetDefault returns default targets. If none is defined in build, return
// those from parent build files.
// Return: default targets a slice of strings
func (build *Build) GetDefault() []string {
	if len(build.Default) > 0 {
		return build.Default
	}
	for i := len(build.Parents) - 1; i >= 0; i-- {
		parent := build.Parents[i]
		parentDefault := parent.GetDefault()
		if len(parentDefault) > 0 {
			return parentDefault
		}
	}
	return nil
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

// GetTarget return target with given name. If not defined in build,
// return target inherited from parent
// - name: the target name as a string
// Return: found target
func (build *Build) GetTarget(name string) *Target {
	target, found := build.Targets[name]
	if found {
		return target
	}
	for i := len(build.Parents) - 1; i >= 0; i-- {
		parent := build.Parents[i]
		target = parent.GetTarget(name)
		if target != nil {
			return target
		}
	}
	return nil
}

// GetParentTarget return parent target with given name.
// - name: the name of the target to run
// Return:
// - target: found parent target, nil if none was found
// - error: if something went wrong
func (build *Build) GetParentTarget(name string) (*Target, error) {
	for i := len(build.Parents) - 1; i >= 0; i-- {
		parent := build.Parents[i]
		target := parent.GetTarget(name)
		if target != nil {
			return target, nil
		}
	}
	return nil, fmt.Errorf("target '%s' not found in parent build files", name)
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
	var listener net.Listener
	var err error
	if listener, err = build.EnsureSingle(context); err != nil {
		return err
	}
	if listener != nil {
		defer func() {
			_ = listener.Close()
		}()
	}
	if len(targets) == 0 {
		targets = build.GetDefault()
		if len(targets) == 0 {
			allTargets := build.GetTargets()
			if len(allTargets) == 1 {
				for key := range allTargets {
					targets = []string{key}
				}
			} else {
				return fmt.Errorf("no default target")
			}
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
	target := build.GetTarget(name)
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
func (build *Build) RunParentTarget(context *Context, name string) error {
	target, err := build.GetParentTarget(name)
	if err != nil {
		return err
	}
	err = context.Stack.Push(target)
	if err != nil {
		return err
	}
	err = target.Steps.Run(context)
	if err != nil {
		return fmt.Errorf("running target '%s': %v", name, err)
	}
	if err := context.Stack.Pop(); err != nil {
		return err
	}
	return nil
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
// Return: a listener and an error if another instance is running on same port
func (build *Build) EnsureSingle(context *Context) (net.Listener, error) {
	if build.Singleton == "" {
		return nil, nil
	}
	expression := build.Singleton
	if IsExpression(expression) {
		expression = expression[1:]
	}
	singleton, err := context.EvaluateExpression(expression)
	if err != nil {
		return nil, fmt.Errorf("evaluating singleton port expression '%s': %v", expression, err)
	}
	port, ok := singleton.(int64)
	if !ok {
		portInt, ok := singleton.(int)
		if !ok {
			return nil, fmt.Errorf("singleton port expression '%s' must return an integer", expression)
		}
		port = int64(portInt)
	}
	return ListenPort(int(port))
}

// ListenPort listens given port:
// - port: port to listen.
// Return: listener and error if any
func ListenPort(port int) (net.Listener, error) {
	if port < 0 || port > 65535 {
		return nil, fmt.Errorf("singleton port port must be between 0 and 65535")
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("listening singleton port: %v", err)
	}
	go func() {
		for {
			_, _ = listener.Accept()
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return listener, nil
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
