package build

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"neon/util"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	// location of the repository root
	REPO_ROOT = "~/.neon"
)

// Build structure
type Build struct {
	File        string
	Dir         string
	Here        string
	Name        string
	Default     []string
	Doc         string
	Repository  string
	Singleton   int
	Scripts     []string
	Extends     []string
	Config      []string
	Properties  util.Object
	Environment map[string]string
	Targets     map[string]*Target
	Context     *Context
	Parents     []*Build
	Index       *Index
	Stack       *Stack
}

// Possible fields for a build file
var FIELDS = []string{"name", "doc", "default", "context", "extends",
	"singleton", "repository", "properties", "configuration",
	"environment", "targets"}

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
	if err := ParseName(object, build); err != nil {
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
		// we return build because it can be used to install plugin
		return build, err
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
	return build, nil
}

// Parse singleton field of the build
func ParseSingleton(object util.Object, build *Build) error {
	if object.HasField("singleton") {
		port, err := object.GetInteger("singleton")
		if err != nil {
			return fmt.Errorf("getting singleton port: %v", err)
		}
		if err := util.Singleton(port); err != nil {
			return fmt.Errorf("another instance of the build is already running")
		}
		build.Singleton = port
	}
	return nil
}

// Parse name field of the build
func ParseName(object util.Object, build *Build) error {
	if object.HasField("name") {
		name, err := object.GetString("name")
		if err != nil {
			return fmt.Errorf("getting build name: %v", err)
		}
		build.Name = name
	}
	return nil
}

// Parse default field of the build
func ParseDefault(object util.Object, build *Build) error {
	if object.HasField("default") {
		list, err := object.GetListStringsOrString("default")
		if err != nil {
			return fmt.Errorf("getting default targets: %v", err)
		}
		build.Default = list
	}
	return nil
}

// Parse doc field of the build
func ParseDoc(object util.Object, build *Build) error {
	if object.HasField("doc") {
		doc, err := object.GetString("doc")
		if err != nil {
			return fmt.Errorf("getting build doc: %v", err)
		}
		build.Doc = doc
	}
	return nil
}

// Parse repository field of the build
func ParseRepository(object util.Object, build *Build) error {
	build.Repository = REPO_ROOT
	if object.HasField("repository") {
		repository, err := object.GetString("repository")
		if err != nil {
			return fmt.Errorf("getting build repository: %v", err)
		}
		build.Repository = repository
	}
	return nil
}

// Parse context field of the build
func ParseContext(object util.Object, build *Build) error {
	if object.HasField("context") {
		scripts, err := object.GetListStringsOrString("context")
		if err != nil {
			return fmt.Errorf("getting context: %v", err)
		}
		build.Scripts = scripts
	}
	return nil
}

// Parse extends field of the build
func ParseExtends(object util.Object, build *Build) error {
	if object.HasField("extends") {
		extends, err := object.GetListStringsOrString("extends")
		if err != nil {
			return fmt.Errorf("parsing parents: %v", err)
		}
		build.Extends = extends
		var parents []*Build
		for _, extend := range build.Extends {
			file := build.PluginPath(extend)
			parent, err := NewBuild(file)
			if err != nil {
				plugin := build.PluginName(extend)
				if plugin != "" {
					return fmt.Errorf("loading parent build file '%s', try installing plugin with 'neon -install=%s'",
						extend, plugin)
				} else {
					return fmt.Errorf("loading parent '%s': %v", extend, err)
				}
			}
			parents = append(parents, parent)
		}
		build.Parents = parents
	}
	return nil
}

// Parse build properties
func ParseProperties(object util.Object, build *Build) error {
	properties := make(map[string]interface{})
	var err error
	if object.HasField("properties") {
		properties, err = object.GetObject("properties")
		if err != nil {
			return fmt.Errorf("parsing properties: %v", err)
		}
	}
	build.Properties = properties
	return nil
}

// Parse build configuration
func ParseConfiguration(object util.Object, build *Build) error {
	if object.HasField("configuration") {
		var config util.Object
		files, err := object.GetListStringsOrString("configuration")
		if err != nil {
			return fmt.Errorf("getting configuration file: %v", err)
		}
		for _, file := range files {
			file = util.ExpandAndJoinToRoot(build.Dir, file)
			source, err := util.ReadFile(file)
			if err != nil {
				return fmt.Errorf("reading configuration file: %v", err)
			}
			err = yaml.Unmarshal(source, &config)
			if err != nil {
				return fmt.Errorf("configuration must be a map with string keys")
			}
			for name, value := range config {
				build.Properties[name] = value
			}
		}
		build.Config = files
	}
	return nil
}

// Parse build environment
func ParseEnvironment(object util.Object, build *Build) error {
	environment := make(map[string]string)
	if object.HasField("environment") {
		env, err := object.GetObject("environment")
		if err != nil {
			return fmt.Errorf("parsing environmen: %v", err)
		}
		environment, err = env.ToMapStringString()
		if err != nil {
			return fmt.Errorf("getting environment: %v", err)
		}
	}
	build.Environment = environment
	return nil
}

// Parse build targets
func ParseTargets(object util.Object, build *Build) error {
	targets := util.Object(make(map[string]interface{}))
	var err error
	if object.HasField("targets") {
		targets, err = object.GetObject("targets")
		if err != nil {
			return fmt.Errorf("parsing targets: %v", err)
		}
	}
	build.Targets = make(map[string]*Target)
	for name := range targets {
		object, err := targets.GetObject(name)
		if err != nil {
			return fmt.Errorf("parsing target '%s': %v", name, err)
		}
		target, err := NewTarget(build, name, object)
		if err != nil {
			return fmt.Errorf("parsing target '%s': %v", name, err)
		}
		build.Targets[name] = target
	}
	return nil
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

// Initialize build:
// - Set build dir and change to it
// - Create context and set it for build
func (build *Build) Init() error {
	os.Chdir(build.Dir)
	context, err := NewContext(build)
	if err != nil {
		return fmt.Errorf("evaluating context: %v", err)
	}
	build.SetDir(build.Dir)
	build.SetContext(context)
	build.SetStack(NewStack())
	return nil
}

// Parse extends field of the build
func (build *Build) LoadParents() error {
	var parents []*Build
	for _, extend := range build.Extends {
		file := build.PluginPath(extend)
		parent, err := NewBuild(file)
		if err != nil {
			plugin := build.PluginName(extend)
			if plugin != "" {
				return fmt.Errorf("loading parent build file '%s', try installing plugin with 'neon -install=%s'",
					extend, plugin)
			} else {
				return fmt.Errorf("loading parent '%s': %v", extend, err)
			}
		}
		parents = append(parents, parent)
	}
	build.Parents = parents
	return nil
}

// Set the build directory, propagating to parents
func (build *Build) SetDir(dir string) {
	build.Dir = dir
	for _, parent := range build.Parents {
		parent.SetDir(dir)
	}
}

// Set the build context, propagating to parents
func (build *Build) SetContext(context *Context) {
	build.Context = context
	for _, parent := range build.Parents {
		parent.SetContext(context)
	}
}

// Set the build stack, propagating to parents
func (build *Build) SetStack(stack *Stack) {
	build.Stack = stack
	for _, parent := range build.Parents {
		parent.SetStack(stack)
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
func (build *Build) Run(targets []string) error {
	if len(targets) == 0 {
		targets = build.GetDefault()
		if len(targets) == 0 {
			return fmt.Errorf("no default target")
		}
	}
	for _, target := range targets {
		err := build.RunTarget(target)
		if err != nil {
			return err
		}
	}
	return nil
}

// Run given target
func (build *Build) RunTarget(name string) error {
	target := build.GetTargetByName(name)
	if target == nil {
		return fmt.Errorf("target '%s' not found", name)
	}
	err := target.Run()
	if err != nil {
		return fmt.Errorf("running target '%s': %v", name, err)
	}
	return nil
}

// Run parent target
func RunParentTarget(build *Build, name string) (bool, error) {
	for _, parent := range build.Parents {
		target := parent.GetTargetByName(name)
		if target != nil {
			err := target.RunSteps()
			if err != nil {
				return true, fmt.Errorf("running target '%s': %v", name, err)
			}
			return true, nil
		} else {
			ok, err := RunParentTarget(parent, name)
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
	re := regexp.MustCompile(`^(\w+/\w+)/.+$`)
	if re.MatchString(name) {
		return re.FindStringSubmatch(name)[1]
	} else {
		return ""
	}
}

// Install given plugin
func (build *Build) Install(plugin string) error {
	re := regexp.MustCompile(`^\w+/\w+$`)
	if !re.MatchString(plugin) {
		return fmt.Errorf("plugin '%s' is invalid", plugin)
	}
	path := build.PluginPath(plugin)
	if util.DirExists(path) {
		Info("Plugin '%s' already installed in '%s'", plugin, path)
		return nil
	}
	absolute := util.ExpandUserHome(path)
	repo := "git@github.com:" + plugin + ".git"
	command := exec.Command("git", "clone", repo, absolute)
	output, err := command.CombinedOutput()
	if err != nil {
		re = regexp.MustCompile("\n\n")
		message := re.ReplaceAllString(string(output), "\n")
		message = strings.TrimSpace(message)
		Info(message)
		return fmt.Errorf("installing plugin '%s'", plugin)
	} else {
		Info("Plugin '%s' installed in '%s'", plugin, path)
	}
	return nil
}
