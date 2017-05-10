package build

import (
	"fmt"
	"strings"
)

func ExpandNeonPath(path string) (string, error) {
	if strings.HasPrefix(path, ":") {
		parts := strings.Split(path[1:], "/")
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("Bad Neon path '%s'", path)
		}
		if len(parts) == 2 {
			parts = []string{parts[0], "latest", parts[1]}
		}
		return fmt.Sprintf("~/.neon/%s/%s/%s", parts[0], parts[1], parts[2]), nil
	} else {
		return path, nil
	}
}
