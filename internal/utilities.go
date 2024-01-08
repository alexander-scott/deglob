package deglob

import (
	"fmt"
	"regexp"
)

// filterPathsBasedOnRegexPattern filters list of paths based on a provided filter regex pattern
// and then returns that new list
func filterPathsBasedOnRegexPattern(inputPaths []string, filter string) []string {
	var filteredPaths []string
	for _, file := range inputPaths {
		matched, _ := regexp.MatchString(filter, file)
		if matched {
			filteredPaths = append(filteredPaths, file)
		}
	}
	return filteredPaths
}

func isCurrentLineNumberWithinAnyTarget(lineNumber int, targets []target) bool {
	for _, target := range targets {
		if lineNumber >= target.start && lineNumber <= target.end {
			return true
		}
	}
	return false
}

func returnTargetInCurrentLineNumber(lineNumber int, targets []target) target {
	for _, target := range targets {
		if lineNumber >= target.start && lineNumber <= target.end {
			return target
		}
	}
	panic("No target found on line " + fmt.Sprint(lineNumber))
}
