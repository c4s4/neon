package build

import (
	"fmt"
	"neon/util"
	"strings"
)

type Repositories []Repository

func NewRepositories() Repositories {
	repositories := []Repository{NewLocalRepository()}
	return repositories
}

func (repos Repositories) GetResource(path string) ([]byte, error) {
	if strings.HasPrefix(path, ":") {
		for _, repo := range repos {
			resource, err := repo.GetResource(path)
			if err != nil {
				return nil, err
			}
			if resource != nil {
				return resource, nil
			}
		}
	} else {
		bytes, err := util.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("loading resource '%s': %v", path, err)
		}
		return bytes, nil
	}
	return nil, nil
}
