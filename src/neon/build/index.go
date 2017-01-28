package build

import (
	"strconv"
	"strings"
)

type Index struct {
	Index []int
}

func NewIndex() *Index {
	index := Index{
		Index: make([]int, 1),
	}
	return &index
}

func (index *Index) Expand() {
	index.Index = append(index.Index, 0)
}

func (index *Index) Shrink() {
	index.Index = index.Index[:len(index.Index)-1]
}

func (index *Index) Set(i int) {
	index.Index[len(index.Index)-1] = i
}

func (index *Index) String() string {
	var str []string
	for _, i := range index.Index {
		str = append(str, strconv.Itoa(i+1))
	}
	return strings.Join(str, ".")
}

func (index *Index) Len() int {
	return len(index.Index)
}
