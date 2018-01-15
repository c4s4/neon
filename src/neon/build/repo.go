package build

import (
	"fmt"
	"neon/util"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
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

// PrintTemplates prints templates in repository:
// - repo: the NeON repository (defaults to '~/.neon')
func PrintTemplates(repo string) {
	if repo == "" {
		repo = DefaultRepo
	}
	repo = util.ExpandUserHome(repo)
	files, err := util.FindFiles(repo, []string{"*/*/*.tpl"}, nil, false)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
}
