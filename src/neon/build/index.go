package build

import (
	"strconv"
	"strings"
)

// Index structure. This keeps track of steps in execution stack. For instance,
// Index will be [1, 2, 3] while running 3rd steps of the second step of the
// first step. This is used to print error location (error in step 1.2.3).
type Index struct {
	Index []int
}

// Make an index
func NewIndex() *Index {
	index := Index{
		Index: make([]int, 1),
	}
	return &index
}

// Expands the index, adding an element after the last
func (index *Index) Expand() {
	index.Index = append(index.Index, 0)
}

// Shrinks the index, removing the last element
func (index *Index) Shrink() {
	index.Index = index.Index[:len(index.Index)-1]
}

// Set the value of the last element
func (index *Index) Set(i int) {
	index.Index[len(index.Index)-1] = i
}

// Return string representation of the index ("1.2.3" for instance)
func (index *Index) String() string {
	var str []string
	for _, i := range index.Index {
		str = append(str, strconv.Itoa(i+1))
	}
	return strings.Join(str, ".")
}

// Return length of the index as an integer
func (index *Index) Len() int {
	return len(index.Index)
}

// Copy return a copy of the index
func (index *Index) Copy() *Index {
	copy := make([]int, len(index.Index))
	for i := 0; i < len(index.Index); i++ {
		copy[i] = index.Index[i]
	}
	return &Index{copy}
}
