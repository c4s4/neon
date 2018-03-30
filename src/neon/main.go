package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	_build "neon/build"
	_ "neon/builtin"
	_ "neon/task"
	"neon/util"
	"os"
	"path/filepath"
	"time"
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

// Version is passed while compiling
var Version string

// Configuration is loaded configuration
var configuration = &Configuration{}

// LoadConfiguration loads configuration file
func LoadConfiguration() (*Configuration, error) {
	configuration := Configuration{}
	file := util.ExpandUserHome(DefaultConfiguration)
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
	// apply grey
	_build.Grey = configuration.Grey
	// apply theme
	if configuration.Theme != "" {
		err := _build.ApplyThemeByName(configuration.Theme)
		if err != nil {
			return nil, err
		}
	}
	// apply custome theme
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
	return &configuration, nil
}

// ParseCommandLine parses command line and returns parsed options
func ParseCommandLine() (string, bool, bool, string, bool, bool, string, bool, bool, string, bool, string, string, bool,
	string, bool, bool, string, bool, []string) {
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
	refs := flag.Bool("refs", false, "Print tasks and builtins reference")
	install := flag.String("install", "", "Install given plugin")
	repo := flag.String("repo", "", "Neon plugin repository for installation")
	grey := flag.Bool("grey", false, "Print on terminal without colors")
	template := flag.String("template", "", "Run given template")
	templates := flag.Bool("templates", false, "List available templates in repository")
	parents := flag.Bool("parents", false, "List available parent build files in repository")
	theme := flag.String("theme", "", "Apply given color theme")
	themes := flag.Bool("themes", false, "Print all available color themes")
	flag.Parse()
	targets := flag.Args()
	return *file, *info, *version, *props, *timeit, *tasks, *task, *targs, *builtins,
		*builtin, *refs, *install, *repo, *grey, *template, *templates, *parents, *theme, *themes, targets
}

// FindBuildFile finds build file and returns its path
// - name: the name of the build file
// - repo: the repository path
// Return:
// - path of found build file
// - an error if something went wrong
func FindBuildFile(name, repo string) (string, string, error) {
	absolute, err := filepath.Abs(name)
	if err != nil {
		return "", "", fmt.Errorf("getting build file path: %v", err)
	}
	file := filepath.Base(absolute)
	dir := filepath.Dir(absolute)
	for {
		path := filepath.Join(dir, file)
		if util.FileExists(path) {
			return path, dir, nil
		}
		if path, ok := configuration.Links[dir]; ok {
			path = _build.LinkPath(path, repo)
			if util.FileExists(path) {
				return path, dir, nil
			}
		}
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
	configuration, err = LoadConfiguration()
	if err != nil {
		PrintError(fmt.Errorf("loading configuration file '%s': %v", DefaultConfiguration, err), 6)
	}
	// parse command line
	file, info, version, props, timeit, tasks, task, targs, builtins, builtin, refs, install, repo, grey, template,
		templates, parents, theme, themes, targets := ParseCommandLine()
	// options that do not require we load build file
	if repo == "" {
		if configuration.Repo != "" {
			repo = configuration.Repo
		} else {
			repo = _build.DefaultRepo
		}
	}
	if grey {
		_build.Grey = true
	}
	if timeit {
		configuration.Time = true
	}
	if tasks {
		_build.PrintTasks()
		return
	} else if task != "" {
		_build.PrintHelpTask(task)
		return
	} else if builtins {
		_build.PrintBuiltins()
		return
	} else if builtin != "" {
		_build.PrintHelpBuiltin(builtin)
		return
	} else if refs {
		_build.PrintReference()
		return
	} else if version {
		_build.Message(Version)
		return
	} else if install != "" {
		err := _build.InstallPlugin(install, repo)
		PrintError(err, 6)
		return
	} else if templates {
		_build.PrintTemplates(repo)
		return
	} else if parents {
		_build.PrintParents(repo)
		return
	} else if theme != "" {
		err := _build.ApplyThemeByName(theme)
		PrintError(err, 8)
	} else if themes {
		_build.PrintThemes()
		return
	}
	// options that do require we load build file
	if template != "" {
		file, err = _build.TemplatePath(template, repo)
		PrintError(err, 1)
	}
	path, base, err := FindBuildFile(file, repo)
	PrintError(err, 1)
	_build.Message("Build: %s", path)
	build, err := _build.NewBuild(path, base, repo)
	PrintError(err, 2)
	err = build.SetCommandLineProperties(props)
	PrintError(err, 3)
	if targs {
		build.PrintTargets()
		return
	} else if info {
		context := _build.NewContext(build)
		err = context.Init()
		PrintError(err, 4)
		err = build.Info(context)
		PrintError(err, 4)
		return
	} else {
		os.Chdir(build.Dir)
		context := _build.NewContext(build)
		err = context.Init()
		PrintError(err, 5)
		err = build.Run(context, targets)
		duration := time.Now().Sub(start)
		if configuration.Time || duration.Seconds() > 10 {
			_build.Message("Build duration: %s", duration.String())
		}
		PrintError(err, 5)
		_build.PrintOk()
		return
	}
}

// PrintError prints an error and exits if any
// - error: the error to check
// - code: the exit code if error is not nil
func PrintError(err error, code int) {
	if err != nil {
		_build.PrintError(err.Error())
		os.Exit(code)
	}
}
