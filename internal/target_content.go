package deglob

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// target is the main data structure for identifying targets with globs within files
type target struct {
	start            int
	end              int
	name             string
	content          []string
	globSearchResult globSearchResult
	globbedFiles     []string
}

var (
	targetNamePattern  = regexp.MustCompile(`name\s=\s\"(?P<name>.*)\"`)
	targetStartPattern = regexp.MustCompile(`^cc_library\(.*`) // TODO: Support cc_binary
	targetEndPattern   = regexp.MustCompile(`^\)\n$`)
)

// extractTargetsFromFileContents finds all targets with globs in a specific file and returns
// them as a slice
func extractTargetsFromFileContents(fileContents []string, filteredFile string) []target {
	var targetsWithGlob []target

	var currentTargetContent []string
	var currentTargetGlobResult globSearchResult
	currentTargetName := ""
	currentlyInTarget := false
	currentTargetStartLineNumber := -1

	currentLineNumber := 0
	for _, line := range fileContents {
		if !currentlyInTarget && targetStartPattern.MatchString(line) {
			currentlyInTarget = true
			currentTargetStartLineNumber = currentLineNumber
		}

		if currentlyInTarget {
			currentTargetContent = append(currentTargetContent, line)

			switch {
			// If we're at the name of a target
			case targetNamePattern.MatchString(line):
				currentTargetName = returnTargetNameFromLine(line)

			// If there's a glob on this line
			case basicGlobCheckRegex.MatchString(line):
				// Check the current line for a glob
				// TODO: Allow a target to have globs on multiple lines/attributes
				checkLineForGlob := extractAllGlobPatternsFromLine(line)
				if checkLineForGlob.globFound {
					currentTargetGlobResult = checkLineForGlob
				}

			// If we're at the end of the target
			case targetEndPattern.MatchString(line):
				if currentTargetGlobResult.globFound {
					target := target{start: currentTargetStartLineNumber, end: currentLineNumber, name: currentTargetName, globSearchResult: currentTargetGlobResult, content: append([]string(nil), currentTargetContent...)} // deep copy the slices
					target.globbedFiles = findFilesFromGlobInTargets(target, filteredFile)

					targetsWithGlob = append(targetsWithGlob, target)

					fmt.Println("Found glob in target: " + currentTargetName + " with attr " + currentTargetGlobResult.globAttr)
				}

				currentlyInTarget = false
				clear(currentTargetContent)
				currentTargetContent = nil
			}
		}

		currentLineNumber++
	}

	return targetsWithGlob
}

// createNewFileContentsIncludingNewTargets takes a slice of targets and the existing file contents and then
// inserts the new de-globbed targets into the existing file and returns the entire file as a slice of strings
func createNewFileContentsIncludingNewTargets(existingFileContents []string, targetsWithGlob []target) []string {
	var newFileContents []string

	currentlyInTarget := false
	var newTargetContent []string
	var currentTarget target

	for index, line := range existingFileContents {
		lineToAdd := line

		switch {
		// If we're currently in a target, our task is just to update several existing target attributes
		case currentlyInTarget:
			if line == currentTarget.globSearchResult.fullLine {
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

func findFilesFromGlobInTargets(t target, filteredFile string) []string {
	packagePath := strings.ReplaceAll(filteredFile, "BUILD", "")

	fmt.Println("Package_path: ", packagePath)
	fmt.Println("Target content: ", t.content)
	fmt.Println("Glob attr: ", t.globSearchResult.globAttr)

	var globbedFiles []string

	for _, globPattern := range t.globSearchResult.globPatterns {
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

func createNewTargetsFromGlobbedFiles(t target) []string {
	var newTargetContent []string
	for _, targetGlobbedFile := range t.globbedFiles {
		newTargetContent = append(newTargetContent, "\n")
		for _, targetContentLine := range t.content {
			switch {
			case targetContentLine == t.globSearchResult.fullLine:
				// If this is the glob line of the target, replace it with the explicit source file
				// TODO: Work with multiple glob patterns
				newSrcLine := strings.ReplaceAll(targetContentLine, "glob([\""+t.globSearchResult.globPatterns[0]+"\"])", "[\""+targetGlobbedFile+"\"]")
				newTargetContent = append(newTargetContent, newSrcLine)
			case targetNamePattern.MatchString(targetContentLine):
				// If this is the name line of the target, replace it with a new name
				newTargetName := generateNewTargetNameForGlobbedFile(t.name, targetGlobbedFile, false, false)
				newNameLine := strings.ReplaceAll(targetContentLine, t.name, newTargetName)
				newTargetContent = append(newTargetContent, newNameLine)
			default:
				newTargetContent = append(newTargetContent, targetContentLine)
			}
		}
	}
	return newTargetContent
}
