package build

import (
	"bufio"
	"fmt"
	"github.com/c4s4/neon/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"runtime"
)

const (
	BaseURL = "http://sweetohm.net/neon"
	VersionURL = BaseURL + "/version"
	BinaryURL = BaseURL + "/neon-%s-%s"
	BinaryPath = "/tmp/neon"
)

// Update updates Neon and repository:
// - repo: the repository path.
// Return: error if something went wrong
func Update(repo string) error {
	response, err := http.Get(VersionURL)
	if err != nil {
		return fmt.Errorf("getting neon version at '%s': %s", VersionURL, err.Error())
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("getting neon version at '%s': bad status code %d", VersionURL, response.StatusCode)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("reading neon version: '%s': %s", VersionURL, err.Error())
	}
	version := strings.TrimSpace(string(body))
	if version != NeonVersion {
		message := fmt.Sprintf("A new version (%s) of neon is available, install[Y/n]? ", version)
		var value string
		if value, err = prompt(message); err != nil {
			return err
		}
		if value == "" || strings.ToLower(value) == "y" {
			var path string
			if path, err = downloadBinary(); err != nil {
				return err
			}
			var installDir string
			if installDir, err = chooseInstallationDir(); err != nil {
				return err
			}
			if err = installBinary(path, installDir); err != nil {
				return err
			}
		}
	}
	return nil
}

func downloadBinary() (string, error) {
	url := fmt.Sprintf(BinaryURL, runtime.GOOS, runtime.GOARCH)
	path := BinaryPath
	if runtime.GOOS == "windows" {
		path = "./neon"
	}
	out, err := os.Create(path)
	if err != nil  {
	  return "", fmt.Errorf("creating destination binary file: %s", err.Error())
	}
	defer out.Close()
	fmt.Print("Downloading neon binary... ")
	response, err := http.Get(url)
	if err != nil {
	  return "", fmt.Errorf("downloading neon binary: %s", err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
	  return "", fmt.Errorf("downloading neon binary, bad status code: %s", response.Status)
	}
	_, err = io.Copy(out, response.Body)
	if err != nil  {
	  return "", fmt.Errorf("saving neon binary: %s", err.Error())
	}
	fmt.Println("OK")
	return path, nil
}

func chooseInstallationDir() (string, error) {
	installDirs := util.FindInPath("neon")
	installDir := ""
	var err error
	if len(installDirs) == 0 {
		if installDir, err = prompt("Installation directory? "); err != nil {
			return "", err
		}
	} else if len(installDirs) == 1 {
		installDir = installDirs[0]
	} else {
		fmt.Println("Please choose installation directory:")
		for i := 0; i < len(installDirs); i++ {
			fmt.Printf("%d) %s\n", i+1, installDirs[i])
		}
		var sel string
		if sel, err = prompt("Installation directory [1]? "); err != nil {
			return "", err
		}
		if sel == "" {
			sel = "1"
		}
		index, err := strconv.Atoi(sel)
		if err != nil || index > len(installDirs) {
			return "", fmt.Errorf("Bad selection")
		}
		installDir = installDirs[index-1]
	}
	return installDir, nil
}

func installBinary(path, installDir string) error {
	fmt.Printf("Installing neon binary in '%s'... ", installDir)
	fmt.Println("OK")
	return nil
}

func prompt(message string) (string, error) {
	fmt.Print(message)
	value, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("reading user input: %v", err)
	}
	return strings.TrimSpace(value), nil
}
