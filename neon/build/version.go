package build

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// NeonVersion is passed while compiling
var NeonVersion string = "UNKNOWN"

// Version encapsulates a software version such as x.y.z
type Version struct {
	String string
	Fields []int
}

// RegexSuffixes is a regexp for version suffixes
var RegexSuffixes = `SNAPSHOT|ALPHA|BETA|RC|snapshot|alpha|beta|rc`

// RegexpVersion is a regexp for version
var RegexpVersion = regexp.MustCompile(`^(\d+(\.\d+)*)(-(` + RegexSuffixes + `)(-\d+)?)?$`)

// NewVersion builds a Version from its string representation
func NewVersion(version string) (*Version, error) {
	if !RegexpVersion.MatchString(version) {
		return nil, fmt.Errorf("%s is not a valid version number", version)
	}
	match := RegexpVersion.FindStringSubmatch(version)
	version = match[1]
	parts := strings.Split(version, ".")
	fields := make([]int, len(parts))
	for i := 0; i < len(parts); i++ {
		field, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid number", parts[i])
		}
		fields[i] = field
	}
	v := Version{
		String: version,
		Fields: fields,
	}
	return &v, nil
}

// Len returns the length of the versions, that is the number of parts
func (v *Version) Len() int {
	return len(v.Fields)
}

// Compare compares two versions.
// Returns:
// - <0 if version is lower than other
// - >0 if version is greater than other
// - =0 if versions are equal
func (v *Version) Compare(o *Version) int {
	min := v.Len()
	if o.Len() < min {
		min = o.Len()
	}
	for i := 0; i < min; i++ {
		c := v.Fields[i] - o.Fields[i]
		if c != 0 {
			return c
		}
	}
	return v.Len() - o.Len()
}
