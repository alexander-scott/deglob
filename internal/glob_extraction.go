package deglob

import (
	"regexp"
	"strings"
)

var (
	basicGlobCheckRegex        = regexp.MustCompile(`glob\(\[.*?\]\)`)
	globRegex                  = regexp.MustCompile(`glob\(\[.*?\]\)`)
	bracketInGlobRegex         = regexp.MustCompile(`\[[^\]]*\]`)
	stringInBracketInGlobRegex = regexp.MustCompile(`"[^"]*"`)
)

// GlobSearchResult is the main data structure for identifying individual glob patterns on a single line
type GlobSearchResult struct {
	globFound        bool
	globAttr         string
	fullLine         string
	globPatterns     []string
	explicitIncludes []string
}

func extractAllGlobPatternsFromLine(line string) GlobSearchResult {
	if !(basicGlobCheckRegex.MatchString(line)) {
		return GlobSearchResult{globFound: false}
	}

	attr := strings.TrimSpace(strings.Split(line, " = ")[0])
	value := strings.Split(line, " = ")[1]

	var foundGlobPatterns []string

	// Find all glob functions within the line
	globMatches := globRegex.FindAllStringSubmatch(value, -1)
	for _, globMatch := range globMatches {
		// For each glob match, find the arrays within them
		bracketsMatches := bracketInGlobRegex.FindAllStringSubmatch(globMatch[0], -1)
		for _, bracketMatch := range bracketsMatches {
			// For each bracket match, find the strings within them
			stringMatches := stringInBracketInGlobRegex.FindAllStringSubmatch(bracketMatch[0], -1)
			for _, stringMatch := range stringMatches {
				globPattern := strings.Trim(stringMatch[0], "\"")
				foundGlobPatterns = append(foundGlobPatterns, globPattern)
			}
		}
	}

	if len(foundGlobPatterns) == 0 {
		panic("Could not find any globs in the following line, yet our basic check regex detected a glob: " + line)
	}

	// Extract any explicit includes mixed with the glob on the same line (e.g. glob([...]) + ["a.h"])
	valueWithoutGlobs := globRegex.ReplaceAllString(value, "")
	var foundExplicitIncludes []string
	explicitBrackets := bracketInGlobRegex.FindAllStringSubmatch(valueWithoutGlobs, -1)
	for _, bracketMatch := range explicitBrackets {
		stringMatches := stringInBracketInGlobRegex.FindAllStringSubmatch(bracketMatch[0], -1)
		for _, stringMatch := range stringMatches {
			foundExplicitIncludes = append(foundExplicitIncludes, strings.Trim(stringMatch[0], "\""))
		}
	}

	return GlobSearchResult{globFound: true, globAttr: attr, globPatterns: foundGlobPatterns, explicitIncludes: foundExplicitIncludes, fullLine: line}
}

// JoinMultiLineGlobs pre-processes file contents to join multi-line glob statements onto a single
// line. For example:
//
//	hdrs = glob([
//	    "*.h",
//	]),
//
// becomes: hdrs = glob([ "*.h", ]),
func JoinMultiLineGlobs(lines []string) []string {
	var result []string
	accumulating := false
	var accumulated string

	for _, line := range lines {
		if !accumulating {
			// Check if this line starts a multi-line glob (has glob([ but no ]) on same line)
			if strings.Contains(line, "glob([") && !strings.Contains(line, "])") {
				accumulating = true
				accumulated = strings.TrimRight(line, "\n")
				continue
			}
			result = append(result, line)
		} else {
			// Accumulating a multi-line glob - trim whitespace from each inner line and append
			trimmed := strings.TrimSpace(strings.TrimRight(line, "\n"))
			accumulated += " " + trimmed
			if strings.Contains(trimmed, "])") {
				accumulating = false
				result = append(result, accumulated+"\n")
				accumulated = ""
			}
		}
	}

	return result
}
