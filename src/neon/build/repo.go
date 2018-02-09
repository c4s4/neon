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
	PluginSite = "github.com"
)

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
	path := filepath.Join(repository, plugin)
	if util.DirExists(path) {
		Message("Plugin '%s' already installed in '%s'", plugin, path)
		return nil
	}
	absolute := util.ExpandUserHome(path)
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
	} else {
		Message("Plugin '%s' installed in '%s'", plugin, path)
	}
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

// FindTemplates finds a template in given repository.
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
	} else {
		if repo == "" {
			repo = DefaultRepo
		}
		match, _ := regexp.MatchString("/[^/]/[^/]/[^/].tpl", name)
		if match {
			return util.ExpandUserHome(filepath.Join(repo, name)), nil
		} else {
			templates, err := FindTemplate(name, repo)
			if err != nil || len(templates) == 0 {
				return "", fmt.Errorf("template '%s' was not found", name)
			}
			if len(templates) > 1 {
				return "", fmt.Errorf("there are %d templates matching name '%s'", len(templates), name)
			}
			return util.ExpandUserHome(filepath.Join(repo, templates[0])), nil
		}
	}
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
