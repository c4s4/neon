package build

import (
	"fmt"
	"neon/util"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	// PluginSite is default site for plugins
	PluginSite = "github.com"
)

// RegexpParentName is regexp for parent name
var RegexpParentName = regexp.MustCompile(`[^/]+/[^/]+/[^/]+.yml`)

// RegexpTemplateName is regexp for template name
var RegexpTemplateName = regexp.MustCompile(`[^/]+/[^/]+/[^/]+.tpl`)

// RegexpScriptName is regexp for script name
var RegexpScriptName = regexp.MustCompile(`[^/]+/[^/]+/[^/]+.ank`)

// RegexpLinkName is regexp for parent name
var RegexpLinkName = regexp.MustCompile(`[^/]+/[^/]+/[^/]+.yml`)

// FindParents finds parent build files in given repository.
// - repository: the NeON repository (defaults to '~/.neon')
// Return:
// - list of parent build files relative to repo.
// - error if something went wrong.
func FindParents(repository string) ([]string, error) {
	files, err := util.FindFiles(repository, []string{"*/*/*.yml"}, nil, false)
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

// FindParent finds a parent in given repository.
// - parent: the parent to find (such as "golang").
// - repository: the NeON repository (defaults to '~/.neon')
// Return:
// - parent path relative to repository (such as "c4s4/build/golang.tpl").
// - error if something went wrong.
func FindParent(parent, repository string) ([]string, error) {
	files, err := FindParents(repository)
	if err != nil {
		return nil, err
	}
	var parents []string
	for _, file := range files {
		start := strings.LastIndex(file, "/") + 1
		end := strings.LastIndex(file, ".")
		name := file[start:end]
		if name == parent {
			parents = append(parents, file)
		}
	}
	return parents, nil
}

// ParentPath returns file path for plugin with given name.
// - name: the name of the plugin (as "c4s4/build/foo.yml" or "foo")
// Return:
// - the plugin path as a string (as /home/casa/.neon/c4s4/build/foo.yml)
// - error if something went wrong
func (build *Build) ParentPath(name string) (string, error) {
	if path.IsAbs(name) {
		return name, nil
	}
	if strings.HasPrefix(name, "./") {
		return filepath.Join(build.Dir, name), nil
	}
	if RegexpParentName.MatchString(name) {
		return util.ExpandUserHome(filepath.Join(build.Repository, name)), nil
	}
	parents, err := FindParent(name, build.Repository)
	if err != nil || len(parents) == 0 {
		return "", fmt.Errorf("parent '%s' was not found", name)
	}
	if len(parents) > 1 {
		return "", fmt.Errorf("there are %d parents matching name '%s'", len(parents), name)
	}
	return util.ExpandUserHome(filepath.Join(build.Repository, parents[0])), nil
}

// InstallPlugin installs given plugin in repository:
// - plugin: the plugin name such as c4s4/build. First part us Github user name
//   and second is repository name for the plugin.
// - repository: plugin repository, defaults to ~/.neon.
// Return: an error if something went wrong downloading plugin.
func InstallPlugin(plugin, repository string) error {
	re := regexp.MustCompile(`^` + RegexpPlugin + `$`)
	if !re.MatchString(plugin) {
		return fmt.Errorf("plugin name '%s' is invalid", plugin)
	}
	pluginPath := filepath.Join(repository, plugin)
	if util.DirExists(pluginPath) {
		Message("Plugin '%s' already installed in '%s'", plugin, pluginPath)
		return nil
	}
	gitRepository := "git://" + PluginSite + "/" + plugin + ".git"
	command := exec.Command("git", "clone", gitRepository, pluginPath)
	Message("Running command '%s'...", strings.Join(command.Args, " "))
	output, err := command.CombinedOutput()
	if err != nil {
		re = regexp.MustCompile("\n\n")
		message := re.ReplaceAllString(string(output), "\n")
		message = strings.TrimSpace(message)
		Message(message)
		return fmt.Errorf("installing plugin '%s'", plugin)
	}
	Message("Plugin '%s' installed in '%s'", plugin, pluginPath)
	return nil
}

// PrintParents prints parent build files in repository:
// - repository: the NeON repository (defaults to '~/.neon')
func PrintParents(repository string) {
	files, err := util.FindFiles(repository, []string{"*/*/*.yml"}, nil, false)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		name := path.Base(file)
		if name != "CHANGELOG.yml" && name != "build.yml" {
			fmt.Println(file)
		}
	}
}

// FindTemplates finds templates in given repository.
// - repository: the NeON repository (defaults to '~/.neon')
// Return:
// - list of template files relative to repo.
// - error if something went wrong.
func FindTemplates(repository string) ([]string, error) {
	files, err := util.FindFiles(repository, []string{"*/*/*.tpl"}, nil, false)
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

// FindTemplate finds a template in given repository.
// - template: the template to find (such as "golang").
// - repository: the NeON repository (defaults to '~/.neon')
// Return:
// - templates path relative to repo (such as "c4s4/build/golang.tpl").
// - error if something went wrong.
func FindTemplate(template, repository string) ([]string, error) {
	files, err := FindTemplates(repository)
	if err != nil {
		return nil, err
	}
	var templates []string
	for _, file := range files {
		start := strings.LastIndex(file, "/") + 1
		end := strings.LastIndex(file, ".")
		name := file[start:end]
		if name == template {
			templates = append(templates, file)
		}
	}
	return templates, nil
}

// TemplatePath return the template path:
// - name: the name of the template (such as 'c4s4/build/golang.tpl')
// - repository: the repository for plugins (defaults to '~/.neon')
// Return: template path (as '~/.neon/c4s4/build/golang.tpl')
func TemplatePath(name, repository string) (string, error) {
	if path.IsAbs(name) || strings.HasPrefix(name, "./") {
		return name, nil
	}
	if RegexpTemplateName.MatchString(name) {
		return filepath.Join(repository, name), nil
	}
	templates, err := FindTemplate(name, repository)
	if err != nil || len(templates) == 0 {
		return "", fmt.Errorf("template '%s' was not found", name)
	}
	if len(templates) > 1 {
		return "", fmt.Errorf("there are %d templates matching name '%s'", len(templates), name)
	}
	return filepath.Join(repository, templates[0]), nil
}

// PrintTemplates prints templates in repository:
// - repository: the NeON repository (defaults to '~/.neon')
func PrintTemplates(repository string) {
	files, err := FindTemplates(repository)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
}

// LinkPath return the link path:
// - name: the name of the build file (such as 'c4s4/build/build.yml')
// - repository: the repository for plugins (defaults to '~/.neon')
// Return: link path (as '~/.neon/c4s4/build/build.yml')
func LinkPath(name, repository string) string {
	if path.IsAbs(name) {
		return name
	}
	if strings.HasPrefix(name, "./") {
		return filepath.Join(repository, name[2:])
	}
	return filepath.Join(repository, name)
}

// ScriptPath returns file path for script with given name.
// - name: the name of the script (as "c4s4/build/foo.ank")
// Return:
// - the script path as a string (as /home/casa/.neon/c4s4/build/foo.ank)
// - error if something went wrong
func (build *Build) ScriptPath(name string) (string, error) {
	if path.IsAbs(name) {
		return name, nil
	}
	if strings.HasPrefix(name, "./") {
		return filepath.Join(build.Dir, name), nil
	}
	if RegexpScriptName.MatchString(name) {
		return util.ExpandUserHome(filepath.Join(build.Repository, name)), nil
	}
	return "", fmt.Errorf("script '%s' was not found", name)
}
