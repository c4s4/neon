package build

import (
	"fmt"
	"io/ioutil"
	"neon/util"
	"path/filepath"
	"sort"
	"strings"
)

type Repository interface {
	GetResource(path string) ([]byte, error)
}

type LocalRepository struct {
	Root string
}

func NewLocalRepository() Repository {
	root := util.ExpandUserHome("~/.neon")
	repository := LocalRepository{
		Root: root,
	}
	return repository
}

func (repo LocalRepository) GetResource(path string) ([]byte, error) {
	group, version, artifact, err := SplitRepositoryPath(path)
	if err != nil {
		return nil, err
	}
	directory := filepath.Join(repo.Root, group)
	if !util.FileExists(directory) {
		return nil, fmt.Errorf("plugin '%s' not found (download it with 'neon -get %s')", group, group)
	}
	if version == "" {
		dirs, err := ioutil.ReadDir(directory)
		if err != nil {
			return nil, fmt.Errorf("listing plugin directory: %v", err)
		}
		var versions = make([]util.Version, len(dirs))
		for i, dir := range dirs {
			if !dir.IsDir() {
				return nil, fmt.Errorf("bad '%s' plugin structure: '%s' is not a directory", group, dir.Name())
			}
			versions[i], err = util.NewVersion(dir.Name())
			if err != nil {
				return nil, err
			}
		}
		sort.Sort(util.Versions(versions))
		version = versions[len(versions)-1].Name
	}
	file := filepath.Join(repo.Root, group, version, artifact)
	resource, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("loading resource '%s': %v", file, err)
	}
	return resource, nil
}

func SplitRepositoryPath(path string) (string, string, string, error) {
	if IsRepositoryPath(path) {
		parts := strings.Split(path[1:], ":")
		if len(parts) < 2 || len(parts) > 3 {
			return "", "", "", fmt.Errorf("Bad Neon path '%s'", path)
		}
		if len(parts) == 2 {
			parts = []string{parts[0], "", parts[1]}
		}
		return parts[0], parts[1], parts[2], nil
	} else {
		return "", "", "", fmt.Errorf("'%s' is not a repository path", path)
	}
}

func IsRepositoryPath(path string) bool {
	return strings.HasPrefix(path, ":")
}
