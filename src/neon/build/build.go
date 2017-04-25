package build

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"neon/util"
	"os"
	"path/filepath"
)

type Build struct {
	File         string
	Dir          string
	Here         string
	Name         string
	Default      []string
	Doc          string
	Scripts      []string
	Properties   util.Object
	Environment  map[string]string
	Targets      map[string]*Target
	Context      *Context
	Parents      []*Build
	Index        *Index
	Repositories Repositories
}

func NewBuild(file string) (*Build, error) {
	build := &Build{}
	build.Repositories = NewRepositories()
	source, err := build.Repositories.GetResource(file)
	if err != nil {
		return nil, fmt.Errorf("loading build file '%s': %v", file, err)
	}
	var object util.Object
	err = yaml.Unmarshal(source, &object)
	if err != nil {
		return nil, fmt.Errorf("build must be a map with string keys")
	}
	err = object.CheckFields([]string{"name", "doc", "default", "context",
		"extends", "singleton", "properties", "configuration", "environment",
		"targets"})
	if err != nil {
		return nil, fmt.Errorf("parsing build file: %v", err)
	}
	if object.HasField("singleton") {
		port, err := object.GetInteger("singleton")
		if err != nil {
			return nil, fmt.Errorf("getting singleton port: %v", err)
		}
		err = util.Singleton(port)
		if err != nil {
			return nil, fmt.Errorf("another instance of the build is already running")
		}
	}
	if object.HasField("name") {
		name, err := object.GetString("name")
		if err != nil {
			return nil, fmt.Errorf("getting build name: %v", err)
		}
		build.Name = name
	}
	if object.HasField("default") {
		list, err := object.GetListStringsOrString("default")
		if err != nil {
			return nil, fmt.Errorf("getting default targets: %v", err)
		}
		build.Default = list
	}
	if object.HasField("doc") {
		doc, err := object.GetString("doc")
		if err != nil {
			return nil, fmt.Errorf("getting build doc: %v", err)
		}
		build.Doc = doc
	}
	if object.HasField("context") {
		scripts, err := object.GetListStringsOrString("context")
		if err != nil {
			return nil, fmt.Errorf("getting context: %v", err)
		}
		build.Scripts = scripts
	}
	if object.HasField("extends") {
		parents, err := object.GetListStringsOrString("extends")
		if err != nil {
			return nil, fmt.Errorf("parsing parents: %v", err)
		}
		var extends []*Build
		for _, parent := range parents {
			extend, err := NewBuild(parent)
			if err != nil {
				return nil, fmt.Errorf("parsing parent '%s': %v", parent, err)
			}
			extends = append(extends, extend)
		}
		build.Parents = extends
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return nil, fmt.Errorf("getting build file path: %v", err)
	}
	build.File = filepath.Base(path)
	build.Dir = filepath.Dir(path)
	here, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting current directory: %v", err)
	}
	build.Here = here
	properties := make(map[string]interface{})
	if object.HasField("properties") {
		properties, err = object.GetObject("properties")
		if err != nil {
			return nil, fmt.Errorf("parsing properties: %v", err)
		}
	}
	build.Properties = properties
	if object.HasField("configuration") {
		var config util.Object
		file, err := object.GetString("configuration")
		if err != nil {
			return nil, fmt.Errorf("getting configuration file: %v", err)
		}
		file = util.FilePath(build.Dir, file)
		source, err := util.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("reading configuration file: %v", err)
		}
		err = yaml.Unmarshal(source, &config)
		if err != nil {
			return nil, fmt.Errorf("configuration must be a map with string keys")
		}
		for name, value := range config {
			build.Properties[name] = value
		}
	}
	environment := make(map[string]string)
	if object.HasField("environment") {
		env, err := object.GetObject("environment")
		if err != nil {
			return nil, fmt.Errorf("parsing environmen: %v", err)
		}
		environment, err = env.ToMapStringString()
		if err != nil {
			return nil, fmt.Errorf("getting environment: %v", err)
		}
	}
	build.Environment = environment
	targets := util.Object(make(map[string]interface{}))
	if object.HasField("targets") {
		targets, err = object.GetObject("targets")
		if err != nil {
			return nil, fmt.Errorf("parsing targets: %v", err)
		}
	}
	build.Targets = make(map[string]*Target)
	for name, _ := range targets {
		object, err := targets.GetObject(name)
		if err != nil {
			return nil, fmt.Errorf("parsing target '%s': %v", name, err)
		}
		target, err := NewTarget(build, name, object)
		if err != nil {
			return nil, fmt.Errorf("parsing target '%s': %v", name, err)
		}
		build.Targets[name] = target
	}
	return build, nil
}

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

func (build *Build) SetContext(context *Context) {
	build.Context = context
	for _, parent := range build.Parents {
		parent.SetContext(context)
	}
}

func (build *Build) SetDir(dir string) {
	build.Dir = dir
	for _, parent := range build.Parents {
		parent.SetDir(dir)
	}
}

func (build *Build) SetProperties(props string) error {
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

func (build *Build) Init() error {
	os.Chdir(build.Dir)
	context, err := NewContext(build)
	if err != nil {
		return fmt.Errorf("evaluating context: %v", err)
	}
	build.SetDir(build.Dir)
	build.SetContext(context)
	return nil
}

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

func (build *Build) Run(targets []string) error {
	if len(targets) == 0 {
		targets = build.GetDefault()
		if len(targets) == 0 {
			return fmt.Errorf("no default target")
		}
	}
	for _, target := range targets {
		err := build.RunTarget(target, NewStack())
		if err != nil {
			return err
		}
	}
	return nil
}

func (build *Build) GetTarget(name string) *Target {
	target, found := build.Targets[name]
	if found {
		return target
	} else {
		for _, parent := range build.Parents {
			target = parent.GetTarget(name)
			if target != nil {
				return target
			}
		}
	}
	return nil
}

func (build *Build) RunTarget(name string, stack *Stack) error {
	target := build.GetTarget(name)
	if target == nil {
		return fmt.Errorf("target '%s' not found", name)
	}
	err := target.Run(stack)
	if err != nil {
		return fmt.Errorf("running target '%s': %v", name, err)
	}
	return nil
}
