package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_build "github.com/c4s4/neon/neon/build"
	_ "github.com/c4s4/neon/neon/builtin"
	_ "github.com/c4s4/neon/neon/task"
	"github.com/c4s4/neon/neon/util"

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
		source, err := os.ReadFile(file)
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
	_build.Gray = configuration.Grey
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

// Options holds parsed command line options
type Options struct {
	File         string
	Info         bool
	Version      bool
	Props        string
	Time         bool
	Tasks        bool
	Task         string
	PrintTargets bool
	Builtins     bool
	Builtin      string
	Tree         bool
	TasksRef     bool
	BuiltinsRef  bool
	Install      string
	Repo         string
	Update       bool
	Batch        bool
	Grey         bool
	Template     string
	Templates    bool
	Parents      bool
	Theme        string
	Themes       bool
	Targets      []string
}

// ParseCommandLine parses command line and returns parsed options
func ParseCommandLine() *Options {
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
	tasksRef := flag.Bool("tasks-ref", false, "Print tasks reference")
	builtinsRef := flag.Bool("builtins-ref", false, "Print builtins reference")
	install := flag.String("install", "", "Install given plugin")
	repo := flag.String("repo", "", "Neon plugin repository for installation")
	update := flag.Bool("update", false, "Update neon and repository")
	batch := flag.Bool("batch", false, "Force neon and repository update in batch mode")
	grey := flag.Bool("grey", false, "Print on terminal without colors")
	template := flag.String("template", "", "Run given template")
	templates := flag.Bool("templates", false, "List available templates in repository")
	parents := flag.Bool("parents", false, "List available parent build files in repository")
	theme := flag.String("theme", "", "Apply given color theme")
	themes := flag.Bool("themes", false, "Print all available color themes")
	flag.Parse()
	targets := flag.Args()
	return &Options{
		File:         *file,
		Info:         *info,
		Version:      *version,
		Props:        *props,
		Time:         *timeit,
		Tasks:        *tasks,
		Task:         *task,
		PrintTargets: *targs,
		Builtins:     *builtins,
		Builtin:      *builtin,
		Tree:         *tree,
		TasksRef:     *tasksRef,
		BuiltinsRef:  *builtinsRef,
		Install:      *install,
		Repo:         *repo,
		Update:       *update,
		Batch:        *batch,
		Grey:         *grey,
		Template:     *template,
		Templates:    *templates,
		Parents:      *parents,
		Theme:        *theme,
		Themes:       *themes,
		Targets:      targets,
	}
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
	if err := run(); err != nil {
		_build.PrintError(err.Error())
		os.Exit(1)
	}
}

func run() error {
	var err error
	start := time.Now()
	// load configuration file
	configPath := DefaultConfiguration
	if envPath := os.Getenv("NEON_CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}
	configuration, err := LoadConfiguration(configPath)
	if err != nil {
		return fmt.Errorf("loading configuration file '%s': %v", DefaultConfiguration, err)
	}
	// parse command line
	opts := ParseCommandLine()
	// options that do not require we load build file
	repo := opts.Repo
	if repo == "" {
		if configuration.Repo != "" {
			repo = configuration.Repo
		} else {
			repo = _build.DefaultRepository
		}
		repo = util.ExpandUserHome(repo)
	}
	_build.Gray = opts.Grey
	configuration.Time = opts.Time
	if printInfo(opts, repo) {
		return nil
	}
	if opts.Version {
		_build.Message(_build.NeonVersion)
		return nil
	} else if opts.Install != "" {
		err := _build.InstallPlugin(opts.Install, repo)
		return err
	} else if opts.Theme != "" {
		err := _build.ApplyThemeByName(opts.Theme)
		if err != nil {
			return err
		}
	} else if opts.Update {
		err := _build.Update(repo, opts.Batch)
		return err
	}
	// options that do require we load build file
	file := opts.File
	if opts.Template != "" {
		file, err = _build.TemplatePath(opts.Template, repo)
		if err != nil {
			return err
		}
	}
	path, base, err := FindBuildFile(file, repo, configuration)
	if err != nil {
		return err
	}
	build, err := _build.NewBuild(path, base, repo, opts.Template != "")
	if err != nil {
		return err
	}
	err = build.SetCommandLineProperties(opts.Props)
	if err != nil {
		return err
	}
	if opts.PrintTargets {
		_build.Message(build.FormatTargets())
	} else if opts.Info {
		context := _build.NewContext(build)
		err = context.Init()
		if err != nil {
			return err
		}
		text, err := build.Info(context)
		if err != nil {
			return err
		}
		_build.Message(text)
	} else if opts.Tree {
		build.Tree()
	} else {
		err = os.Chdir(build.Dir)
		if err != nil {
			return err
		}
		context := _build.NewContext(build)
		err = context.Init()
		if err != nil {
			return err
		}
		err = build.Run(context, opts.Targets)
		duration := time.Since(start)
		if configuration.Time || duration.Seconds() > 10 {
			_build.InfoArgs("Build duration: %s", duration.String())
		}
		if err != nil {
			return err
		}
		_build.PrintOk()
	}
	return nil
}

// printInfo prints build information if requested
func printInfo(opts *Options, repo string) bool {
	if opts.Tasks {
		_build.Message(_build.InfoTasks())
		return true
	} else if opts.Task != "" {
		_build.Message(_build.InfoTask(opts.Task))
		return true
	} else if opts.Builtins {
		_build.Message(_build.InfoBuiltins())
		return true
	} else if opts.Builtin != "" {
		_build.Message(_build.InfoBuiltin(opts.Builtin))
		return true
	} else if opts.Templates {
		_build.Message(_build.InfoTemplates(repo))
		return true
	} else if opts.Parents {
		_build.Message(_build.InfoParents(repo))
		return true
	} else if opts.Themes {
		_build.Message(_build.InfoThemes())
		return true
	} else if opts.TasksRef {
		_build.Message(_build.InfoTasksReference())
		return true
	} else if opts.BuiltinsRef {
		_build.Message(_build.InfoBuiltinsReference())
		return true
	}
	return false
}

// (Empty or removed)
