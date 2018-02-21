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

// FindParents finds parent build files in given repository.
// - repo: the NeON repository (defaults to '~/.neon')
// Return:
// - list of parent build files relative to repo.
// - error if something went wrong.
func FindParents(repo string) ([]string, error) {
	if repo == "" {
		repo = DefaultRepo
	}
	repo = util.ExpandUserHome(repo)
	files, err := util.FindFiles(repo, []string{"*/*/*.yml"}, nil, false)
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

// FindParent finds a parent in given repository.
// - parent: the parent to find (such as "golang").
// - repo: the NeON repository (defaults to '~/.neon')
// Return:
// - parent path relative to repo (such as "c4s4/build/golang.tpl").
// - error if something went wrong.
func FindParent(parent, repo string) ([]string, error) {
	files, err := FindParents(repo)
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
	repopath := filepath.Join(repository, plugin)
	if util.DirExists(repopath) {
		Message("Plugin '%s' already installed in '%s'", plugin, repopath)
		return nil
	}
	absolute := util.ExpandUserHome(repopath)
	repo := "git://" + PluginSite + "/" + plugin + ".git"
	command := exec.Command("git", "clone", repo, absolute)
	Message("Running command '%s'...", strings.Join(command.Args, " "))
	output, err := command.CombinedOutput()
	if err != nil {
		re = regexp.MustCompile("\n\n")
		message := re.ReplaceAllString(string(output), "\n")
		message = strings.TrimSpace(message)
		Message(message)
		return fmt.Errorf("installing plugin '%s'", plugin)
	}
	Message("Plugin '%s' installed in '%s'", plugin, repopath)
	return nil
}

// PrintParents prints parent build files in repository:
// - repo: the NeON repository (defaults to '~/.neon')
func PrintParents(repo string) {
	if repo == "" {
		repo = DefaultRepo
	}
	repo = util.ExpandUserHome(repo)
	files, err := util.FindFiles(repo, []string{"*/*/*.yml"}, nil, false)
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
// - repo: the NeON repository (defaults to '~/.neon')
// Return:
// - list of template files relative to repo.
// - error if something went wrong.
func FindTemplates(repo string) ([]string, error) {
	if repo == "" {
		repo = DefaultRepo
	}
	repo = util.ExpandUserHome(repo)
	files, err := util.FindFiles(repo, []string{"*/*/*.tpl"}, nil, false)
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

// FindTemplate finds a template in given repository.
// - template: the template to find (such as "golang").
// - repo: the NeON repository (defaults to '~/.neon')
// Return:
// - templates path relative to repo (such as "c4s4/build/golang.tpl").
// - error if something went wrong.
func FindTemplate(template, repo string) ([]string, error) {
	files, err := FindTemplates(repo)
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
// - repo: the repository for plugins (defaults to '~/.neon')
// Return: template path (as '~/.neon/c4s4/build/golang.tpl')
func TemplatePath(name, repo string) (string, error) {
	if path.IsAbs(name) || strings.HasPrefix(name, "./") {
		return name, nil
	}
	if repo == "" {
		repo = DefaultRepo
	}
	if RegexpTemplateName.MatchString(name) {
		return util.ExpandUserHome(filepath.Join(repo, name)), nil
	}
	templates, err := FindTemplate(name, repo)
	if err != nil || len(templates) == 0 {
		return "", fmt.Errorf("template '%s' was not found", name)
	}
	if len(templates) > 1 {
		return "", fmt.Errorf("there are %d templates matching name '%s'", len(templates), name)
	}
	return util.ExpandUserHome(filepath.Join(repo, templates[0])), nil
}

// PrintTemplates prints templates in repository:
// - repo: the NeON repository (defaults to '~/.neon')
func PrintTemplates(repo string) {
	files, err := FindTemplates(repo)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
}
