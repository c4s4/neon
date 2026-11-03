package build

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/c4s4/neon/neon/util"
)

const (
	// versionURL is the URL of the file with latest release version, deployed
	// with binaries by the release (see deploy target in build.yml)
	versionURLDefault  = "https://sweetohm.net/dist/neon/version"
	installCommandCurl = "sh -c \"$(curl https://sweetohm.net/dist/neon/install)\""
	installCommandWget = "sh -c \"$(wget -O - https://sweetohm.net/dist/neon/install)\""
	// updateTimeout is the timeout for the request checking latest version
	updateTimeout = 10 * time.Second
)

// versionURL is the URL of the version file (a variable for testing purpose)
var versionURL = versionURLDefault

// Update prints if a new NeON release is available and updates plugins of
// repository:
// - repository: the repository path.
// - batch: update plugins without asking for confirmation.
// Return: error if something went wrong
func Update(repository string, batch bool) error {
	printNewRelease()
	return updateRepository(repository, batch)
}

func printNewRelease() {
	client := &http.Client{Timeout: updateTimeout}
	response, err := client.Get(versionURL)
	if err != nil {
		return
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		return
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	version := strings.TrimSpace(string(body))
	latest, err := NewVersion(version)
	if err != nil {
		return
	}
	current, err := NewVersion(NeonVersion)
	if err != nil {
		fmt.Printf("You are running a development version of neon (%s), latest release is %s\n", NeonVersion, version)
	} else if latest.Compare(current) > 0 {
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

func updateRepository(repository string, batch bool) error {
	plugins, err := util.FindFiles(repository, []string{"*/*"}, nil, true)
	if err != nil {
		return fmt.Errorf("searching plugins: %v", err)
	}
	if len(plugins) > 0 {
		fmt.Println("Plugins:")
		for _, plugin := range plugins {
			path := filepath.Join(repository, plugin)
			// skip directories that are not git repositories
			if !util.DirExists(filepath.Join(path, ".git")) {
				fmt.Printf("- %s: not a git repository, skipped\n", plugin)
				continue
			}
			// fetch repository
			cmd := exec.Command("git", "fetch")
			cmd.Dir = path
			bytes, err := cmd.CombinedOutput()
			if err != nil {
				_, _ = fmt.Fprintln(ErrorOutput, string(bytes))
				return fmt.Errorf("fetching repository: %v", err)
			}
			// get branch name
			cmd = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
			cmd.Dir = path
			bytes, err = cmd.CombinedOutput()
			if err != nil {
				_, _ = fmt.Fprintln(ErrorOutput, string(bytes))
				return fmt.Errorf("getting branch for plugin: %v", err)
			}
			branch := strings.TrimSpace(string(bytes))
			// get hash for local and remote
			cmd = exec.Command("git", "rev-parse", "HEAD")
			cmd.Dir = path
			bytes, err = cmd.CombinedOutput()
			if err != nil {
				_, _ = fmt.Fprintln(ErrorOutput, string(bytes))
				return fmt.Errorf("getting hash for local: %v", err)
			}
			hashLocal := strings.TrimSpace(string(bytes))
			origin := "origin/" + branch
			cmd = exec.Command("git", "rev-parse", origin)
			cmd.Dir = path
			bytes, err = cmd.CombinedOutput()
			if err != nil {
				_, _ = fmt.Fprintln(ErrorOutput, string(bytes))
				return fmt.Errorf("getting hash for remote: %v", err)
			}
			hashRemote := strings.TrimSpace(string(bytes))
			if hashRemote == hashLocal {
				fmt.Printf("- %s [%s]: OK\n", plugin, branch)
			} else {
				perform := batch
				if !batch {
					reader := bufio.NewReader(os.Stdin)
					fmt.Printf("- %s [%s]: Update [Y/n]? ", plugin, branch)
					input, err := reader.ReadString('\n')
					if err != nil {
						return fmt.Errorf("reading user input: %v", err)
					}
					response := strings.ToLower(strings.TrimSpace(string(input)))
					if response == "" || response == "y" {
						perform = true
					}
				} else {
					fmt.Printf("- %s [%s]: NOK\n", plugin, branch)
				}
				if perform {
					fmt.Printf("  Updating %s [%s]... ", plugin, branch)
					cmd = exec.Command("git", "pull")
					cmd.Dir = path
					bytes, err = cmd.CombinedOutput()
					if err != nil {
						message := strings.ReplaceAll(string(bytes), "\n", "")
						fmt.Println()
						_, _ = fmt.Fprintln(ErrorOutput, message)
						return fmt.Errorf("updating plugin: %v", err)
					}
					fmt.Println("OK")
				}
			}
		}
	} else {
		fmt.Println("No plugin found in repository")
	}
	return nil
}
