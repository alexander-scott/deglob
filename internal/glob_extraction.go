package deglob

import (
	"regexp"
	"strings"
)

var (
	basicGlobCheckRegex        = regexp.MustCompile(`glob\(\[.*]\)`)
	globRegex                  = regexp.MustCompile(`glob\(\[[a-zA-Z0-9_\] + \["*\.\/]*\)|\]\)`)
	bracketInGlobRegex         = regexp.MustCompile(`\[[a-zA-Z0-9_ ,"*\.\/]*\]`)
	stringInBracketInGlobRegex = regexp.MustCompile(`\"[a-zA-Z0-9_ ,*\.\/]*\"`)
)

// GlobSearchResult is the main data structure for identifying individual glob patterns on a single line
type GlobSearchResult struct {
	globFound    bool
	globAttr     string
	fullLine     string
	globPatterns []string
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

	return GlobSearchResult{globFound: true, globAttr: attr, globPatterns: foundGlobPatterns, fullLine: line}
}
