package deglob

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	basicGlobCheckRegex        = regexp.MustCompile(`glob\(\[.*]\)`)
	globRegex                  = regexp.MustCompile(`glob\(\[[a-zA-Z0-9_\] + \["*\.\/]*\)|\]\)`)
	bracketInGlobRegex         = regexp.MustCompile(`\[[a-zA-Z0-9_ ,"*\.\/]*\]`)
	stringInBracketInGlobRegex = regexp.MustCompile(`\"[a-zA-Z0-9_ ,*\.\/]*\"`)
)

type GlobSearchResult struct {
	globFound    bool
	globAttr     string
	lineNumber   int
	globPatterns []string
}

func extractAllGlobPatternsFromLine(line string, lineNumber int) GlobSearchResult {
	// 1) Confirm that we are on the glob line
	// 2) Get a list of all of the glob() functions
	// 3) Get a list of all of the brackets in each glob function
	// 4) Get a list of everything between "" in each brackets
	// 5) Return

	if !(basicGlobCheckRegex.MatchString(line)) {
		return GlobSearchResult{globFound: false}
	}

	attr := strings.TrimSpace(strings.Split(line, " = ")[0])
	value := strings.Split(line, " = ")[1]

	var foundGlobPatterns []string

	// Find all glob functions within the line
	globMatches := globRegex.FindAllStringSubmatch(value, -1)
	for _, globMatch := range globMatches {
		fmt.Println("FOund glob match: " + globMatch[0])
		// For each glob match, find the arrays within them
		bracketsMatches := bracketInGlobRegex.FindAllStringSubmatch(globMatch[0], -1)
		for _, bracketMatch := range bracketsMatches {
			fmt.Println("FOund bracket in glob match: " + bracketMatch[0])
			// For each bracket match, find the strings within them
			stringMatches := stringInBracketInGlobRegex.FindAllStringSubmatch(bracketMatch[0], -1)
			for _, stringMatch := range stringMatches {
				fmt.Println("FOund string in bracket in glob match: " + stringMatch[0])
				globPattern := strings.Trim(stringMatch[0], "\"")
				foundGlobPatterns = append(foundGlobPatterns, globPattern)
			}
		}
	}

	return GlobSearchResult{globFound: true, globAttr: attr, globPatterns: foundGlobPatterns, lineNumber: lineNumber}
}
