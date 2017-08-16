package build

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"neon/util"
)

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

// Parse shell field of the build
func ParseShell(object util.Object, build *Build) error {
	if object.HasField("shell") {
		field := object["shell"]
		if util.IsMap(field) {
			shell := make(map[string][]string)
			mapInterface, err := util.ToMapStringInterface(field)
			if err != nil {
				return err
			}
			for os, v := range mapInterface {
				s, err := util.ToSliceString(v)
				if err != nil {
					return err
				}
				shell[os] = s
			}
			build.Shell = shell
		} else if util.IsSlice(field) {
			slice, err := util.ToSliceString(field)
			if err != nil {
				return err
			}
			build.Shell = map[string][]string {
				"default": slice,
			}
		} else {
			return fmt.Errorf("shell must be a list of strings or a map of list of strings")
		}
	} else {
		build.Shell = map[string][]string {
			"default": []string{"sh", "-c"},
			"windows": []string{"cmd", "/c"},
		}
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
				return fmt.Errorf("loading parent build file '%s': %v", extend, err)
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
			return fmt.Errorf("parsing environment: %v", err)
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
