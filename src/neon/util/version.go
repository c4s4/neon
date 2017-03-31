package util

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Name    string
	Numbers []int
}

type Versions []Version

func NewVersion(name string) (Version, error) {
	version := Version{Name: name}
	strs := strings.Split(name, ".")
	numbers := make([]int, len(strs))
	for i, str := range strs {
		var err error
		numbers[i], err = strconv.Atoi(str)
		if err != nil {
			return Version{}, fmt.Errorf("version '%s' is invalid: %v", name, err)
		}
	}
	version.Numbers = numbers
	return version, nil
}

func (version Version) Less(other Version) bool {
	maxLen := len(version.Numbers)
	if len(other.Numbers) > len(version.Numbers) {
		maxLen = len(other.Numbers)
	}
	for i := 0; i < maxLen; i++ {
		if version.Numbers[i] < other.Numbers[i] {
			return true
		}
		if version.Numbers[i] > other.Numbers[i] {
			return false
		}
	}
	if len(version.Numbers) < len(other.Numbers) {
		return true
	} else {
		return false
	}
}

func (versions Versions) Len() int {
	return len(versions)
}

func (versions Versions) Swap(i, j int) {
	versions[i], versions[j] = versions[j], versions[i]
}

func (versions Versions) Less(i, j int) bool {
	return versions[i].Less(versions[j])
}
