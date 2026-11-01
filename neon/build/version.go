package build

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// NeonVersion is set at build time with -ldflags (UNKNOWN otherwise)
var NeonVersion = "UNKNOWN"

// Version encapsulates a software version such as x.y.z
type Version struct {
	String string
	Fields []int
	// Suffix is the rank of the suffix (see suffixRanks), release versions
	// (without suffix) have the highest rank
	Suffix int
	// Number is the number after the suffix (such as 2 in 1.0-RC-2)
	Number int
}

// suffixRanks gives the order of version suffixes: a SNAPSHOT is lower than an
// ALPHA, which is lower than a BETA, then RC and finally release
var suffixRanks = map[string]int{
	"snapshot": 1,
	"alpha":    2,
	"beta":     3,
	"rc":       4,
	"":         5,
}

// RegexSuffixes is a regexp for version suffixes
var RegexSuffixes = `SNAPSHOT|ALPHA|BETA|RC|snapshot|alpha|beta|rc`

// RegexpVersion is a regexp for version, with an optional leading "v" (as in
// git tags such as "v1.2.3")
var RegexpVersion = regexp.MustCompile(`^[vV]?(\d+(\.\d+)*)(-(` + RegexSuffixes + `)([-.]?\d+)?)?$`)

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
	number := 0
	if match[5] != "" {
		digits := strings.TrimLeft(match[5], "-.")
		n, err := strconv.Atoi(digits)
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid number", digits)
		}
		number = n
	}
	v := Version{
		String: version,
		Fields: fields,
		Suffix: suffixRanks[strings.ToLower(match[4])],
		Number: number,
	}
	return &v, nil
}

// IsReleaseVersion tells if given version is a release version number (such
// as "1.2.3"), as opposed to development versions (such as a git hash or
// "UNKNOWN")
func IsReleaseVersion(version string) bool {
	return RegexpVersion.MatchString(version)
}

// Len returns the length of the version, that is the number of fields
func (v *Version) Len() int {
	return len(v.Fields)
}

// Compare compares two versions. Missing fields count as zeros (1.0 is equal
// to 1.0.0), then suffixes are compared (1.0-SNAPSHOT < 1.0-ALPHA < 1.0-BETA <
// 1.0-RC < 1.0) and finally suffix numbers (1.0-RC-1 < 1.0-RC-2).
// Returns:
// - <0 if version is lower than other
// - >0 if version is greater than other
// - =0 if versions are equal
func (v *Version) Compare(o *Version) int {
	length := v.Len()
	if o.Len() > length {
		length = o.Len()
	}
	for i := 0; i < length; i++ {
		c := v.field(i) - o.field(i)
		if c != 0 {
			return c
		}
	}
	if v.Suffix != o.Suffix {
		return v.Suffix - o.Suffix
	}
	return v.Number - o.Number
}

// field returns field at given index, 0 if version has fewer fields
func (v *Version) field(i int) int {
	if i < len(v.Fields) {
		return v.Fields[i]
	}
	return 0
}

// SortVersions sorts a list of versions in place, in ascending order. A
// leading "v" (as in git tags such as "v1.2.3") is ignored. Strings that are
// not valid versions are placed first, in lexical order.
// - versions: the versions to sort
func SortVersions(versions []string) {
	parsed := make(map[string]*Version, len(versions))
	for _, version := range versions {
		if v, err := NewVersion(strings.TrimPrefix(strings.TrimPrefix(version, "v"), "V")); err == nil {
			parsed[version] = v
		}
	}
	sort.SliceStable(versions, func(i, j int) bool {
		vi, vj := parsed[versions[i]], parsed[versions[j]]
		switch {
		case vi == nil && vj == nil:
			return versions[i] < versions[j]
		case vi == nil:
			return true
		case vj == nil:
			return false
		}
		if c := vi.Compare(vj); c != 0 {
			return c < 0
		}
		return versions[i] < versions[j]
	})
}

// CompareNeonVersion compares NeON version with given version. It panics,
// which is the way builtins report errors, if NeON is a development version or
// if given version is not valid.
// - version: the version to compare with
// Return: <0 if NeON version is lower, >0 if greater, 0 if equal
func CompareNeonVersion(version string) int {
	n, err := NewVersion(NeonVersion)
	if err != nil {
		panic(fmt.Sprintf("NeON development version '%s' can't be compared with a version number", NeonVersion))
	}
	v, err := NewVersion(version)
	if err != nil {
		panic(fmt.Sprintf("version '%s' is not a valid version number", version))
	}
	return n.Compare(v)
}
