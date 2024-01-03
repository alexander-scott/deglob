package deglob

import (
	"fmt"
	"regexp"
)

func FilterPathsBasedOnRegexPattern(input_paths []string, filter string) []string {
	// 2) Filter paths based on regex pattern
	var filtered_paths []string
	for _, file := range input_paths {
		matched, _ := regexp.MatchString(filter, file)
		if matched {
			filtered_paths = append(filtered_paths, file)
		}
	}
	return filtered_paths
}

func is_current_line_number_within_any_target(line_number int, targets []Target) bool {
	for _, target := range targets {
		if line_number >= target.start && line_number <= target.end {
			return true
		}
	}
	return false
}

func return_target_in_current_line_number(line_number int, targets []Target) Target {
	for _, target := range targets {
		if line_number >= target.start && line_number <= target.end {
			return target
		}
	}
	panic("No target found on line " + fmt.Sprint(line_number))
}
