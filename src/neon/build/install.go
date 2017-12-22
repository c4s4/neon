package build

import (
	"neon/util"
	"regexp"
	"fmt"
	"os/exec"
	"strings"
	"path/filepath"
)

// Install given plugin
// - plugin: the plugin name such as c4s4/build. First part us Github user name
//   and second is repository name for the plugin.
// - repository: plugin repository, defaults to ~/.neon.
func InstallPlugin(plugin, repository string) error {
	re := regexp.MustCompile(`^` + RE_PLUGIN + `$`)
	if !re.MatchString(plugin) {
		return fmt.Errorf("plugin name '%s' is invalid", plugin)
	}
	path := filepath.Join(repository, plugin)
	if util.DirExists(path) {
		Message("Plugin '%s' already installed in '%s'", plugin, path)
		return nil
	}
	absolute := util.ExpandUserHome(path)
	repo := "git://github.com/" + plugin + ".git"
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