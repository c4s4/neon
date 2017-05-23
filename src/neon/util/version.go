package util

import (
	"fmt"
	"strconv"
	"strings"
)

// A version contains its string value and its numbers
type Version struct {
	Name    string
	Numbers []int
}

// Versions is a list of versions
type Versions []Version

// Make a version from its string representation
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

// Tells if version is less that other version
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

// Return numbers of numbers of the version
func (versions Versions) Len() int {
	return len(versions)
}

// Swap two version numbers
func (versions Versions) Swap(i, j int) {
	versions[i], versions[j] = versions[j], versions[i]
}

// Tells if a number is less that the other
func (versions Versions) Less(i, j int) bool {
	return versions[i].Less(versions[j])
}
