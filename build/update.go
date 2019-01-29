package build

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/c4s4/neon/util"
)

const (
	versionURL         = "http://sweetohm.net/neon/version"
	installCommandCurl = "sh -c \"$(curl http://sweetohm.net/neon/install)\""
	installCommandWget = "sh -c \"$(wget -O - http://sweetohm.net/neon/install)\""
)

// Update updates Neon and repository:
// - repo: the repository path.
// Return: error if something went wrong
func Update(repo string) error {
	printNewRelease()
	updateRepository(repo)
	return nil
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
	body, err := ioutil.ReadAll(response.Body)
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
			// get branch name
			cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
			cmd.Dir = path
			branch, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("getting branch for plugin: %v", err)
			}
			// get hash for local and remote
			cmd = exec.Command("git", "rev-parse", "HEAD")
			cmd.Dir = path
			hashLocal, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("getting hash for local: %v", err)
			}
			cmd = exec.Command("git", "rev-parse", string(branch))
			cmd.Dir = path
			hashRemote, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("getting hash for remote: %v", err)
			}
			println("local:", string(hashLocal))
			println("remote:", string(hashRemote))
			if output == nil {
				fmt.Printf("- %s: OK", plugin)
			} else {
				reader := bufio.NewReader(os.Stdin)
				fmt.Printf("- %s: Update [Y/n]? ", plugin)
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("reading user input: %v", err)
				}
				response := string(input)
				if response == "" || strings.ToLower(response) == "y" {
					cmd = exec.Command("git", "pull")
					cmd.Dir = path
					err := cmd.Run()
					if err != nil {
						return fmt.Errorf("updating plugin: %v", err)
					}
				}
			}
			// FIXME
			println(path)
		}
	} else {
		fmt.Println("No plugin found in repository")
	}
	return nil
}
