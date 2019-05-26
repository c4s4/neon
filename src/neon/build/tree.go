package build

import (
	"fmt"
	"strings"
)

// Tree prints the inheritance tree for given build
func (build *Build) Tree() {
	fmt.Println(name(build.File))
	for index, parent := range build.Parents {
		parent.SubTree("", index < len(build.Parents)-1)
	}
}

// SubTree prints the inheritance SubTree for given build
func (build *Build) SubTree(margin string, next bool) {
	name := name(build.File)
	fmt.Print(margin)
	if next {
		margin += "│ "
		fmt.Println("├─" + name)
	} else {
		margin += "  "
		fmt.Println("└─" + name)
	}
	for index, parent := range build.Parents {
		parent.SubTree(margin, index < len(build.Parents)-1)
	}
}

// name return name of given file without extension
func name(file string) string {
	return file[:strings.LastIndex(file, ".")]
}
