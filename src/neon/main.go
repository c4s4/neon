package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	_build "neon/build"
	_ "neon/builtin"
	_ "neon/task"
	"neon/util"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	// DefaultBuildFile is the default name for build file
	DefaultBuildFile = "build.yml"
	// DefaultConfiguration is the default location for configuration file
	DefaultConfiguration = "~/.neon/settings.yml"
)

// Configuration holds configuration properties
type Configuration struct {
	// Grey disables color output
	Grey bool
	// Theme applies named theme
	Theme string
	// Colors of custom theme
	Colors *_build.Colors
	// Time will print execution time
	Time bool
	// Repo location
	Repo string
	// Links associates build files to directories
	Links map[string]string
}

// ParseConfiguration parses configuration file:
// - file: the configuration file to parse.
// Return: built Configuration struct and error if any
func ParseConfiguration(file string) (*Configuration, error) {
	var configuration Configuration
	file = util.ExpandUserHome(file)
	if util.FileExists(file) {
		source, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(source, &configuration)
		if err != nil {
			return nil, err
		}
	}
	return &configuration, nil
}

// LoadConfiguration loads configuration file:
// - file: configuration file to load
// Return: configuration and error if any
func LoadConfiguration(file string) (*Configuration, error) {
	configuration, err := ParseConfiguration(file)
	if err != nil {
		return nil, err
	}
	// apply grey
	_build.Grey = configuration.Grey
	// apply theme
	if configuration.Theme != "" {
		err := _build.ApplyThemeByName(configuration.Theme)
		if err != nil {
			return nil, err
		}
	}
	// apply custom theme
	if configuration.Colors != nil {
		theme, err := _build.ParseTheme(configuration.Colors)
		if err != nil {
			return nil, err
		}
		_build.ApplyTheme(theme)
	}
	// expand user homes in files
	abs := make(map[string]string)
	for dir, build := range configuration.Links {
		abs[util.ExpandUserHome(dir)] = util.ExpandUserHome(build)
	}
	configuration.Links = abs
	return configuration, nil
}

// ParseCommandLine parses command line and returns parsed options
func ParseCommandLine() (string, bool, bool, string, bool, bool, string, bool, bool, string,
	bool, bool, string, string, bool, bool, string, bool, bool, string, bool, []string) {
	file := flag.String("file", DefaultBuildFile, "Build file to run")
	info := flag.Bool("info", false, "Print build information")
	version := flag.Bool("version", false, "Print neon version")
	props := flag.String("props", "", "Build properties")
	timeit := flag.Bool("time", false, "Print build duration")
	tasks := flag.Bool("tasks", false, "Print tasks list")
	task := flag.String("task", "", "Print help on given task")
	targs := flag.Bool("targets", false, "Print targets list")
	builtins := flag.Bool("builtins", false, "Print builtins list")
	builtin := flag.String("builtin", "", "Print help on given builtin")
	tree := flag.Bool("tree", false, "Print inheritance tree")
	refs := flag.Bool("refs", false, "Print tasks and builtins reference")
	install := flag.String("install", "", "Install given plugin")
	repo := flag.String("repo", "", "Neon plugin repository for installation")
	update := flag.Bool("update", false, "Update neon and repository")
	grey := flag.Bool("grey", false, "Print on terminal without colors")
	template := flag.String("template", "", "Run given template")
	templates := flag.Bool("templates", false, "List available templates in repository")
	parents := flag.Bool("parents", false, "List available parent build files in repository")
	theme := flag.String("theme", "", "Apply given color theme")
	themes := flag.Bool("themes", false, "Print all available color themes")
	flag.Parse()
	targets := flag.Args()
	return *file, *info, *version, *props, *timeit, *tasks, *task, *targs, *builtins,
		*builtin, *tree, *refs, *install, *repo, *update, *grey, *template, *templates,
		*parents, *theme, *themes, targets
}

// FindBuildFile finds build file and returns its path
// - name: the name of the build file
// - repo: the repository path
// Return:
// - path of found build file
// - an error if something went wrong
func FindBuildFile(name, repo string, configuration *Configuration) (string, string, error) {
	absolute, err := filepath.Abs(name)
	if err != nil {
		return "", "", fmt.Errorf("getting build file path: %v", err)
	}
	file := filepath.Base(absolute)
	dir := filepath.Dir(absolute)
	for {
		// first look in configuration file links
		if path, ok := configuration.Links[dir]; ok {
			path = _build.LinkPath(path, repo)
			if util.FileExists(path) {
				return path, dir, nil
			}
		}
		// if not found, look in current directory
		path := filepath.Join(dir, file)
		if util.FileExists(path) {
			return path, dir, nil
		}
		// if not found, loop in parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", "", fmt.Errorf("build file not found")
		}
		dir = parent
	}
}

// Program entry point
func main() {
	var err error
	start := time.Now()
	// load configuration file
	configuration, err := LoadConfiguration(DefaultConfiguration)
	if err != nil {
		PrintError(fmt.Errorf("loading configuration file '%s': %v", DefaultConfiguration, err))
	}
	// parse command line
	file, info, version, props, timeit, tasks, task, targs, builtins, builtin, tree, refs, install, repo,
		update, grey, template, templates, parents, theme, themes, targets := ParseCommandLine()
	// options that do not require we load build file
	if repo == "" {
		if configuration.Repo != "" {
			repo = configuration.Repo
		} else {
			repo = _build.DefaultRepository
		}
		repo = util.ExpandUserHome(repo)
	}
	_build.Grey = grey
	configuration.Time = timeit
	if printInfo(tasks, builtins, templates, parents, themes, refs, task, builtin, repo) {
		return
	}
	if version {
		_build.Message(_build.NeonVersion)
		return
	} else if install != "" {
		err := _build.InstallPlugin(install, repo)
		PrintError(err)
		return
	} else if theme != "" {
		err := _build.ApplyThemeByName(theme)
		PrintError(err)
	} else if update {
		err := _build.Update(repo)
		PrintError(err)
		return
	}
	// options that do require we load build file
	if template != "" {
		file, err = _build.TemplatePath(template, repo)
		PrintError(err)
	}
	path, base, err := FindBuildFile(file, repo, configuration)
	PrintError(err)
	_build.Info("Build %s", path)
	build, err := _build.NewBuild(path, base, repo, template != "")
	PrintError(err)
	err = build.SetCommandLineProperties(props)
	PrintError(err)
	if targs {
		_build.Message(build.FormatTargets())
	} else if info {
		context := _build.NewContext(build)
		err = context.Init()
		PrintError(err)
		text, err := build.Info(context)
		PrintError(err)
		_build.Message(text)
	} else if tree {
		build.Tree()
	} else {
		os.Chdir(build.Dir)
		context := _build.NewContext(build)
		err = context.Init()
		PrintError(err)
		err = build.Run(context, targets)
		duration := time.Now().Sub(start)
		if configuration.Time || duration.Seconds() > 10 {
			_build.Info("Build duration: %s", duration.String())
		}
		PrintError(err)
		_build.PrintOk()
	}
}

func printInfo(tasks, builtins, templates, parents, themes, refs bool, task, builtin, repo string) bool {
	if tasks {
		_build.Message(_build.InfoTasks())
		return true
	} else if task != "" {
		_build.Message(_build.InfoTask(task))
		return true
	} else if builtins {
		_build.Message(_build.InfoBuiltins())
		return true
	} else if builtin != "" {
		_build.Message(_build.InfoBuiltin(builtin))
		return true
	} else if templates {
		_build.Message(_build.InfoTemplates(repo))
		return true
	} else if parents {
		_build.Message(_build.InfoParents(repo))
		return true
	} else if themes {
		_build.Message(_build.InfoThemes())
		return true
	} else if refs {
		_build.Message(_build.InfoReference())
		return true
	}
	return false
}

// PrintError prints an error and exits if any
// - error: the error to check
func PrintError(err error, args ...string) {
	if err != nil {
		_build.PrintError(err.Error())
		os.Exit(1)
	}
}
