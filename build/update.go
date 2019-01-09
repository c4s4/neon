package build

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"

	"github.com/c4s4/neon/util"
)

const (
	VersionURL         = "http://sweetohm.net/neon/version"
	InstallCommandCurl = "sh -c \"$(curl http://sweetohm.net/neon/install)\""
	InstallCommandWget = "sh -c \"$(wget -O - http://sweetohm.net/neon/install)\""
)

// Update updates Neon and repository:
// - repo: the repository path.
// Return: error if something went wrong
func Update(repo string) error {
	printNewRelease()
	return nil
}

func printNewRelease() {
	response, err := http.Get(VersionURL)
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
				fmt.Println(InstallCommandCurl)
			} else if len(util.FindInPath("wget")) > 0 {
				fmt.Println("You can install it with:")
				fmt.Println(InstallCommandWget)
			}
		}
	} else {
		fmt.Println("Your version of neon is the latest one")
	}
}
