package deglob

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// Target is the main data structure for identifying targets with globs within files
type Target struct {
	start            int
	end              int
	name             string
	content          []string
	globSearchResult GlobSearchResult
	globbedFiles     []string
}

var (
	targetGlobPattern  = regexp.MustCompile(`\s+(?P<attr>.*)\s=\sglob\(\[(?P<files>.*)\]\)`)
	targetNamePattern  = regexp.MustCompile(`name\s=\s\"(?P<name>.*)\"`)
	targetStartPattern = regexp.MustCompile(`^cc_library\(.*`) // TODO: Support cc_binary
	targetEndPattern   = regexp.MustCompile(`^\)\n$`)
)

// ExtractTargetsFromFileContents finds all targets with globs in a specific file and returns
// them as a slice
func ExtractTargetsFromFileContents(fileContents []string, filteredFile string) []Target {
	var targetsWithGlob []Target

	var currentTargetContent []string
	var currentTargetGlobResult GlobSearchResult
	currentTargetName := ""
	currentlyInTarget := false
	currentTargetStartLineNumber := -1

	currentLineNumber := 0
	for _, line := range fileContents {
		if targetStartPattern.MatchString(line) && !currentlyInTarget {
			currentlyInTarget = true
			currentTargetStartLineNumber = currentLineNumber
		}

		if currentlyInTarget {
			currentTargetContent = append(currentTargetContent, line)
		}

		if currentlyInTarget && targetNamePattern.MatchString(line) {
			matches := targetNamePattern.FindStringSubmatch(line)
			nameIndex := targetNamePattern.SubexpIndex("name")
			currentTargetName = matches[nameIndex]
		}

		if currentlyInTarget && targetEndPattern.MatchString(line) {
			// Only track target if the target had a glob
			if currentTargetGlobResult.globFound {
				target := Target{start: currentTargetStartLineNumber, end: currentLineNumber, name: currentTargetName, globSearchResult: currentTargetGlobResult, content: append([]string(nil), currentTargetContent...)} // deep copy the slices
				target.globbedFiles = findFilesFromGlobInTargets(target, filteredFile)

				targetsWithGlob = append(targetsWithGlob, target)

				fmt.Println("Found glob in target: " + currentTargetName + " with attr " + currentTargetGlobResult.globAttr + " on line " + fmt.Sprint(currentTargetGlobResult.lineNumber))
			}

			currentlyInTarget = false
			clear(currentTargetContent)
			currentTargetContent = nil
		}

		// 4) Find any globs in the files and the patterns they capture
		if currentlyInTarget && targetGlobPattern.MatchString(line) {
			currentTargetGlobResult = extractAllGlobPatternsFromLine(line, currentLineNumber)
		}

		currentLineNumber++
	}

	return targetsWithGlob
}

// CreateNewFileContentsIncludingNewTargets takes a slice of targets and the existing file contents and then
// inserts the new de-globbed targets into the existing file and returns the entire file as a slice of strings
func CreateNewFileContentsIncludingNewTargets(existingFileContents []string, targetsWithGlob []Target) []string {
	var newFileContents []string

	currentlyInTarget := false
	var newTargetContent []string
	var currentTarget Target

	for index, line := range existingFileContents {
		lineToAdd := line

		switch {
		// If we're currently in a target, our task is just to update several existing target attributes
		case currentlyInTarget:
			if index == currentTarget.globSearchResult.lineNumber {
				lineToAdd = strings.ReplaceAll(line, currentTarget.globSearchResult.globAttr, "deps")
				// TODO: Work with more glob patterns
				lineToAdd = strings.ReplaceAll(lineToAdd, "glob([\""+currentTarget.globSearchResult.globPatterns[0]+"\"])", "["+createListOfNewTargetNamesFromTarget(currentTarget)+"]")
			}

			// The current target might require the generation of entirely new targets. If so, we can't add it
			// whilst we're still in the target as that will mess up the layout. Therefore we add it to the
			// `newTargetContent` slice so it can be written to file once we're outside of the current target.
			if index == currentTarget.end {
				currentlyInTarget = false
				newTargetContent = createNewTargetsFromGlobbedFiles(currentTarget)
			}
		case len(newTargetContent) > 0:
			newFileContents = append(newFileContents, newTargetContent...)
			clear(newTargetContent)
			newTargetContent = nil
		case !currentlyInTarget && isCurrentLineNumberWithinAnyTarget(index, targetsWithGlob):
			currentlyInTarget = true
			currentTarget = returnTargetInCurrentLineNumber(index, targetsWithGlob)
		}

		newFileContents = append(newFileContents, lineToAdd)
	}

	// If we've reached the end of the previous file, but we still have new targets to add, we add can it now
	// to the very bottom of the file
	if len(newTargetContent) > 0 {
		newFileContents = append(newFileContents, newTargetContent...)
		clear(newTargetContent)
	}

	return newFileContents
}

func findFilesFromGlobInTargets(target Target, filteredFile string) []string {
	packagePath := strings.ReplaceAll(filteredFile, "BUILD", "")

	fmt.Println("Package_path: ", packagePath)
	fmt.Println("Target content: ", target.content)
	fmt.Println("Glob attr: ", target.globSearchResult.globAttr)

	var globbedFiles []string

	for _, globPattern := range target.globSearchResult.globPatterns {
		globCmd := strings.ReplaceAll(filteredFile, "BUILD", globPattern)
		fmt.Println("Glob_cmd: ", globCmd)
		files, _ := filepath.Glob(globCmd)

		for _, file := range files {
			globbedFullFilePath := packagePath + file
			globbedPackageRelativeFilePath := strings.ReplaceAll(globbedFullFilePath, packagePath, "")
			fmt.Println("Relative_file_path: ", globbedPackageRelativeFilePath)
			globbedFiles = append(globbedFiles, globbedPackageRelativeFilePath)
		}
	}

	return globbedFiles
}

func createNewTargetsFromGlobbedFiles(target Target) []string {
	var newTargetContent []string
	for _, targetGlobbedFile := range target.globbedFiles {
		newTargetContent = append(newTargetContent, "\n")
		for _, targetContentLine := range target.content {
			switch {
			case targetGlobPattern.MatchString(targetContentLine):
				// If this is the glob line of the target, replace it with the explicit source file
				// TODO: Work with multiple glob patterns
				newSrcLine := strings.ReplaceAll(targetContentLine, "glob([\""+target.globSearchResult.globPatterns[0]+"\"])", "[\""+targetGlobbedFile+"\"]")
				newTargetContent = append(newTargetContent, newSrcLine)
			case targetNamePattern.MatchString(targetContentLine):
				// If this is the name line of the target, replace it with a new name
				newTargetName := generateNewTargetNameForGlobbedFile(target.name, targetGlobbedFile, false, false)
				newNameLine := strings.ReplaceAll(targetContentLine, target.name, newTargetName)
				newTargetContent = append(newTargetContent, newNameLine)
			default:
				newTargetContent = append(newTargetContent, targetContentLine)
			}
		}
	}
	return newTargetContent
}

func createListOfNewTargetNamesFromTarget(target Target) string {
	var newTargetNames []string
	for _, targetGlobbedFile := range target.globbedFiles {
		for _, targetContentLine := range target.content {
			if targetNamePattern.MatchString(targetContentLine) {
				newTargetName := generateNewTargetNameForGlobbedFile(target.name, targetGlobbedFile, true, true)
				newTargetNames = append(newTargetNames, newTargetName)
			}
		}
	}
	return strings.Join(newTargetNames, ", ")
}

func generateNewTargetNameForGlobbedFile(targetName string, globbedFileName string, asLabel bool, wrapWithQuotes bool) string {
	newNameSuffix := strings.Split(globbedFileName, ".")[0]
	newNameSuffix = strings.ReplaceAll(newNameSuffix, "/", "_")
	newTargetName := targetName + "_" + newNameSuffix
	if asLabel {
		newTargetName = ":" + newTargetName
	}
	if wrapWithQuotes {
		newTargetName = "\"" + newTargetName + "\""
	}
	return newTargetName
}
