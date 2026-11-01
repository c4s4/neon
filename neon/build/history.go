package build

import (
	"strings"
)

// History is the list of names of targets that already ran during a build
type History struct {
	Targets []string
}

// NewHistory makes a new history
// Returns: a pointer to the history
func NewHistory() *History {
	history := History{
		Targets: make([]string, 0),
	}
	return &history
}

// Contains tells if the history contains given target
// - name: name of the target to test
// Returns: a boolean telling if target is in the history
func (history *History) Contains(name string) bool {
	for _, target := range history.Targets {
		if name == target {
			return true
		}
	}
	return false
}

// Push a target on the history
// - target: target to push on the history
// Return: always nil (kept for symmetry with Stack.Push)
func (history *History) Push(target *Target) error {
	history.Targets = append(history.Targets, target.Name)
	return nil
}

// String returns string representation of the history, such as:
// "foo, bar, spam"
// Return: the history as a string
func (history *History) String() string {
	names := make([]string, len(history.Targets))
	copy(names, history.Targets)
	return strings.Join(names, ", ")
}

// Copy returns a copy of the history
// Return: pointer to a copy of the history
func (history *History) Copy() *History {
	another := make([]string, len(history.Targets))
	copy(another, history.Targets)
	return &History{another}
}
