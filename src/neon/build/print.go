package build

import (
	"fmt"
	"neon/util"
)

func Info(message string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(message, args...))
}

func Title(message string) {
	util.PrintColor(util.Yellow(message))
}
