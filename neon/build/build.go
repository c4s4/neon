package build

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/c4s4/neon/neon/util"
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
	FileDir     string
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

// NewBuild creates a Build from a build file.
// This is the public entry point. It initializes an empty inheritance chain
// and delegates the actual work to newBuildInternal.
func NewBuild(file, base, repo string, template bool) (*Build, error) {
	return newBuildInternal(file, base, repo, template, nil)
}

// newBuildInternal performs the real build creation while tracking the chain
// of build files being loaded to detect cyclic inheritance. The chain is a
// stack: a file may appear several times in the inheritance tree (diamond
// inheritance), but not twice on the same branch.
func newBuildInternal(file, base, repo string, template bool, chain []string) (*Build, error) {
	// Resolve absolute path for cycle detection
	absPath, err := filepath.Abs(util.ExpandUserHome(file))
	if err != nil {
		return nil, fmt.Errorf("resolving absolute path for '%s': %v", file, err)
	}
	// Detect cycle
	for _, path := range chain {
		if path == absPath {
			cycle := append(append([]string{}, chain...), absPath)
			return nil, fmt.Errorf("dependency cycle detected: %s", strings.Join(cycle, " → "))
		}
	}
	chain = append(append([]string{}, chain...), absPath)
	object, build, err := parseBuildFile(file)
	if err != nil {
		return nil, err
	}
	if err := SetDirectories(build, base); err != nil {
		return nil, err
	}
	build.FileDir = filepath.Dir(absPath)
	if err := object.CheckFields(Fields); err != nil {
		return nil, fmt.Errorf("parsing build file: %v", err)
	}
	if err := ParseFields(object, build, repo); err != nil {
		return nil, err
	}
	// Relative dotenv files of parents are relative to the parent build file,
	// those of the root build are relative to the build directory
	dotEnvDir := build.FileDir
	if len(chain) == 1 {
		dotEnvDir = build.Dir
	}
	for i, file := range build.DotEnv {
		file = util.ExpandUserHome(file)
		if !filepath.IsAbs(file) {
			file = filepath.Join(dotEnvDir, file)
		}
		build.DotEnv[i] = file
	}
	// Resolve parents with the current inheritance chain
	build.Parents, err = build.getParentsInternal(chain)
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
	object, err := util.DecodeYAMLObject(source)
	if err != nil {
		return nil, nil, fmt.Errorf("build must be a map with string keys: %v", err)
	}
	return object, build, nil
}

// SetDirectories sets build directory and current directory (Here):
// - build: the build to set directories for.
// - base: the build directory.
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
// context. These properties overwrite those defined in the build file.
// - props: properties as a YAML map
// Return: error if something went wrong
func (build *Build) SetCommandLineProperties(props string) error {
	object, err := util.DecodeYAMLObject([]byte(props))
	if err != nil {
		return fmt.Errorf("parsing command line properties: properties must be a map with string keys: %v", err)
	}
	for name, value := range object {
		build.Properties[name] = value
	}
	return nil
}

// GetParents returns parent build objects.
// Return: list of build objects and an error if any.
func (build *Build) GetParents() ([]*Build, error) {
	var chain []string
	// Put the current build file in the chain to prevent self-reference.
	if build.FileDir != "" {
		chain = append(chain, filepath.Join(build.FileDir, build.File))
	}
	return build.getParentsInternal(chain)
}

// getParentsInternal resolves parent builds with the given inheritance chain.
func (build *Build) getParentsInternal(chain []string) ([]*Build, error) {
	var parents []*Build
	for _, extend := range build.Extends {
		file, err := build.ParentPath(extend)
		if err != nil {
			return nil, fmt.Errorf("searching parent build file '%s': %v", extend, err)
		}
		parent, err := newBuildInternal(file, filepath.Dir(file), build.Repository, build.Template, chain)
		if err != nil {
			return nil, fmt.Errorf("loading parent build file '%s': %v", extend, err)
		}
		parents = append(parents, parent)
	}
	return parents, nil
}

// GetProperties returns build properties, including those inherited from
// parents. Parents are already loaded, thus their properties already include
// those of their own parents.
// Return: build properties as an Object
func (build *Build) GetProperties() util.Object {
	var properties = make(map[string]interface{})
	for _, parent := range build.Parents {
		for name, value := range parent.Properties {
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
		for name, value := range parent.Environment {
			environment[name] = value
		}
	}
	for name, value := range build.Environment {
		environment[name] = value
	}
	return environment
}

// GetDotEnv returns the list of dotenv files to load in environment, including
// those inherited from parents. A file inherited several times (diamond
// inheritance) appears once, at its first position.
// Return: list of dotenv files as a slice of strings
func (build *Build) GetDotEnv() []string {
	var dotenv []string
	for _, parent := range build.Parents {
		dotenv = append(dotenv, parent.DotEnv...)
	}
	dotenv = append(dotenv, build.DotEnv...)
	return unique(dotenv)
}

// unique removes duplicates from a list, keeping first occurrences
func unique(list []string) []string {
	seen := make(map[string]bool, len(list))
	var result []string
	for _, element := range list {
		if !seen[element] {
			seen[element] = true
			result = append(result, element)
		}
	}
	return result
}

// GetTargets returns build targets, including those inherited from parents
// Return: targets as a map of targets with their name as keys
func (build *Build) GetTargets() map[string]*Target {
	var targets = make(map[string]*Target)
	// last parent has precedence, as in GetTarget
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

// GetDefault returns default targets. If none is defined in build, return
// those from parent build files.
// Return: default targets as a slice of strings
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

// GetScripts returns a list of context scripts to run, with their path
// resolved relative to the build file that declares them. A script inherited
// several times (diamond inheritance) is loaded once.
// Return: the list of context script paths and an error if any
func (build *Build) GetScripts() ([]string, error) {
	var scripts []string
	for _, parent := range build.Parents {
		parentScripts, err := parent.GetScripts()
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, parentScripts...)
	}
	for _, script := range build.Scripts {
		path, err := build.ScriptPath(script)
		if err != nil {
			return nil, fmt.Errorf("getting script path '%s': %v", script, err)
		}
		scripts = append(scripts, path)
	}
	return unique(scripts), nil
}

// GetSingleton returns singleton port expression. If none is defined in
// build, return the one of parent build files.
// Return: singleton port expression, empty if none
func (build *Build) GetSingleton() string {
	if build.Singleton != "" {
		return build.Singleton
	}
	for i := len(build.Parents) - 1; i >= 0; i-- {
		if singleton := build.Parents[i].GetSingleton(); singleton != "" {
			return singleton
		}
	}
	return ""
}

// GetTarget returns target with given name. If not defined in build,
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

// GetParentTarget returns parent target with given name.
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
// Return: an error if parent target was not found or failed
func (build *Build) RunParentTarget(context *Context, name string) error {
	target, err := build.GetParentTarget(name)
	if err != nil {
		return err
	}
	if err := context.Stack.Push(target); err != nil {
		return err
	}
	defer func() {
		_ = context.Stack.Pop()
	}()
	if err := target.Steps.Run(context); err != nil {
		return fmt.Errorf("running target '%s': %v", name, err)
	}
	return nil
}

// GetShell returns shell for current OS.
// Return:
// - shell as a slice of strings (such as ["sh", "-c"])
// - error if something went wrong
func (build *Build) GetShell() ([]string, error) {
	if shell, ok := build.Shell[runtime.GOOS]; ok {
		return shell, nil
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
	expression := build.GetSingleton()
	if expression == "" {
		return nil, nil
	}
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

// ListenPort listens given port on loopback interface:
// - port: port to listen.
// Return: listener and error if any
func ListenPort(port int) (net.Listener, error) {
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("singleton port must be between 1 and 65535")
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, fmt.Errorf("listening singleton port: %v", err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()
	return listener, nil
}

// CheckVersion evaluates version expressions of the build and its parents to
// check that NeON version is OK. Version check is skipped for development
// builds, which version is not a release number (such as a git hash).
func (build *Build) CheckVersion(context *Context) error {
	if !IsReleaseVersion(NeonVersion) {
		return nil
	}
	for _, parent := range build.Parents {
		if err := parent.CheckVersion(context); err != nil {
			return err
		}
	}
	if build.Version == "" {
		return nil
	}
	expression := build.Version
	if IsExpression(expression) {
		expression = expression[1:]
	}
	result, err := context.EvaluateExpression(expression)
	if err != nil {
		return fmt.Errorf("evaluating version expression: %v", err)
	}
	versionOK, ok := result.(bool)
	if !ok {
		return fmt.Errorf("version expression should return a boolean")
	}
	if !versionOK {
		return fmt.Errorf("neon version '%s' doesn't meet requirements in version field of '%s'", NeonVersion, build.File)
	}
	return nil
}
