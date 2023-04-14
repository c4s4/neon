package build

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/c4s4/neon/neon/util"
)

const (
	versionURL         = "http://sweetohm.net/neon/version"
	installCommandCurl = "sh -c \"$(curl http://sweetohm.net/neon/install)\""
	installCommandWget = "sh -c \"$(wget -O - http://sweetohm.net/neon/install)\""
)

// Update updates Neon and repository:
// - repository: the repository path.
// Return: error if something went wrong
func Update(repository string) error {
	printNewRelease()
	return updateRepository(repository)
}

func printNewRelease() {
	response, err := http.Get(versionURL)
	if err != nil {
		return
	}
	if response.StatusCode != 200 {
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	version := strings.TrimSpace(string(body))
	if version != NeonVersion {
		fmt.Printf("Neon %s is available\n", version)
		if runtime.GOOS != "windows" {
			if len(util.FindInPath("curl")) > 0 {
				fmt.Println("You can install it with:")
				fmt.Println(installCommandCurl)
			} else if len(util.FindInPath("wget")) > 0 {
				fmt.Println("You can install it with:")
				fmt.Println(installCommandWget)
			}
		}
	} else {
		fmt.Println("Your version of neon is the latest one")
	}
}

func updateRepository(repository string) error {
	plugins, err := util.FindFiles(repository, []string{"*/*"}, nil, true)
	if err != nil {
		return fmt.Errorf("searching plugins: %v", err)
	}
	if len(plugins) > 0 {
		fmt.Println("Plugins:")
		for _, plugin := range plugins {
			path := path.Join(repository, plugin)
			// fetch repository
			cmd := exec.Command("git", "fetch")
			cmd.Dir = path
			bytes, err := cmd.CombinedOutput()
			if err != nil {
				println(string(bytes))
				return fmt.Errorf("fetching repository: %v", err)
			}
			// get branch name
			cmd = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
			cmd.Dir = path
			bytes, err = cmd.CombinedOutput()
			if err != nil {
				println(string(bytes))
				return fmt.Errorf("getting branch for plugin: %v", err)
			}
			branch := strings.TrimSpace(string(bytes))
			// get hash for local and remote
			cmd = exec.Command("git", "rev-parse", "HEAD")
			cmd.Dir = path
			bytes, err = cmd.CombinedOutput()
			if err != nil {
				println(string(bytes))
				return fmt.Errorf("getting hash for local: %v", err)
			}
			hashLocal := strings.TrimSpace(string(bytes))
			origin := "origin/" + branch
			cmd = exec.Command("git", "rev-parse", origin)
			cmd.Dir = path
			bytes, err = cmd.CombinedOutput()
			if err != nil {
				println(string(bytes))
				return fmt.Errorf("getting hash for remote: %v", err)
			}
			hashRemote := strings.TrimSpace(string(bytes))
			if hashRemote == hashLocal {
				fmt.Printf("- %s [%s]: OK\n", plugin, branch)
			} else {
				reader := bufio.NewReader(os.Stdin)
				fmt.Printf("- %s [%s]: Update [Y/n]? ", plugin, branch)
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("reading user input: %v", err)
				}
				response := strings.ToLower(strings.TrimSpace(string(input)))
				if response == "" || response == "y" {
					cmd = exec.Command("git", "pull")
					cmd.Dir = path
					bytes, err = cmd.CombinedOutput()
					if err != nil {
						println(string(bytes))
						return fmt.Errorf("updating plugin: %v", err)
					}
				}
			}
		}
	} else {
		fmt.Println("No plugin found in repository")
	}
	return nil
}
