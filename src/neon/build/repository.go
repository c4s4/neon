package build

import (
	"fmt"
	"io/ioutil"
	"neon/util"
	"path/filepath"
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
	file := filepath.Join(repo.Root, group, version, artifact)
	resource, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("loading resource '%s': %v", file, err)
	}
	return resource, nil
}

func SplitRepositoryPath(path string) (string, string, string, error) {
	if IsRepositoryPath(path) {
		parts := strings.Split(path[1:], "/")
		if len(parts) < 2 || len(parts) > 3 {
			return "", "", "", fmt.Errorf("Bad Neon path '%s'", path)
		}
		if len(parts) == 2 {
			parts = []string{parts[0], "latest", parts[1]}
		}
		return parts[0], parts[1], parts[2], nil
	} else {
		return "", "", "", fmt.Errorf("'%s' is not a repository path", path)
	}
}

func IsRepositoryPath(path string) bool {
	return strings.HasPrefix(path, ":")
}
