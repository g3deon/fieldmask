package fieldmask

import (
	"strings"
)

const pathSeparator = "."

// removeEmptyPaths filters out empty or whitespace-only strings from the provided slice of paths.
func removeEmptyPaths(paths []string) []string {
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		if strings.TrimSpace(path) != "" {
			result = append(result, path)
		}
	}
	return result
}

// removeDuplicatePaths removes duplicate strings from the input slice of paths and returns a new slice with unique paths.
func removeDuplicatePaths(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		if _, exists := seen[path]; !exists {
			seen[path] = struct{}{}
			result = append(result, path)
		}
	}
	return result
}
