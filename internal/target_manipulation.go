package deglob

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// Target is the main data structure for identifying targets with globs within files
type Target struct {
	start              int
	end                int
	name               string
	content            []string
	globAttr           string
	globAttrLineNumber int
	globPatterns       []string
	globbedFiles       []string
}

var (
	targetGlobPattern  = regexp.MustCompile(`\s(?P<attr>.*)\s=\sglob\(\[(?P<files>.*)\]\)`)
	targetNamePattern  = regexp.MustCompile(`name\s=\s\"(?P<name>.*)\"`)
	targetStartPattern = regexp.MustCompile(`^cc_library\(.*`) // TODO: Support cc_binary
	targetEndPattern   = regexp.MustCompile(`^\)\n$`)
)

// ExtractTargetsFromFileContents finds all targets with globs in a specific file and returns
// them as a slice
func ExtractTargetsFromFileContents(fileContents []string, filteredFile string) []Target {
	var targetsWithGlob []Target

	var currentTargetGlobPatterns []string
	var currentTargetContent []string
	currentTargetName := ""
	currentTargetGlobAttr := ""
	currentTargetGlobAttrLineNumber := -1
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
			if len(currentTargetGlobPatterns) > 0 {
				target := Target{start: currentTargetStartLineNumber, end: currentLineNumber, name: currentTargetName, globPatterns: append([]string(nil), currentTargetGlobPatterns...), content: append([]string(nil), currentTargetContent...), globAttr: currentTargetGlobAttr, globAttrLineNumber: currentTargetGlobAttrLineNumber} // deep copy the slices
				target.globbedFiles = findFilesFromGlobInTargets(target, filteredFile)

				targetsWithGlob = append(targetsWithGlob, target)

				fmt.Println("Found glob in target: " + currentTargetName + " with attr " + currentTargetGlobAttr + " on line " + fmt.Sprint(currentTargetGlobAttrLineNumber))
			}

			currentlyInTarget = false
			clear(currentTargetGlobPatterns)
			clear(currentTargetContent)
		}

		// 4) Find any globs in the files and the patterns they capture
		if currentlyInTarget && targetGlobPattern.MatchString(line) {
			matches := targetGlobPattern.FindStringSubmatch(line)
			filesIndex := targetGlobPattern.SubexpIndex("files")
			attrIndex := targetGlobPattern.SubexpIndex("attr")

			// TODO: Tackle multiple strings inside a single glob
			globPattern := strings.Trim(matches[filesIndex], "\"")
			currentTargetGlobAttr = strings.TrimSpace(matches[attrIndex])
			currentTargetGlobPatterns = append(currentTargetGlobPatterns, globPattern)
			currentTargetGlobAttrLineNumber = currentLineNumber
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
		case currentlyInTarget:
			if index == currentTarget.globAttrLineNumber {
				lineToAdd = strings.ReplaceAll(line, currentTarget.globAttr, "deps")
				lineToAdd = strings.ReplaceAll(lineToAdd, "glob([\""+currentTarget.globPatterns[0]+"\"])", "["+createListOfNewTargetNamesFromTarget(currentTarget)+"]")
			}
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

	return newFileContents
}

func findFilesFromGlobInTargets(target Target, filteredFile string) []string {
	packagePath := strings.ReplaceAll(filteredFile, "BUILD", "")

	fmt.Println("Package_path: ", packagePath)
	fmt.Println("Target content: ", target.content)
	fmt.Println("Glob attr: ", target.globAttr)

	var globbedFiles []string

	for _, globPattern := range target.globPatterns {
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
				newSrcLine := strings.ReplaceAll(targetContentLine, "glob([\""+target.globPatterns[0]+"\"])", "[\""+targetGlobbedFile+"\"]")
				newTargetContent = append(newTargetContent, newSrcLine)
			case targetNamePattern.MatchString(targetContentLine):
				// If this is the name line of the target, replace it with a new name
				newNameSuffix := strings.Split(targetGlobbedFile, ".")[0]
				newNameLine := strings.ReplaceAll(targetContentLine, target.name, target.name+"_"+newNameSuffix)
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
				// If this is the name line of the target, replace it with a new name
				newNameSuffix := strings.Split(targetGlobbedFile, ".")[0]
				newTargetNames = append(newTargetNames, "\":"+target.name+"_"+newNameSuffix+"\"")
			}
		}
	}
	return strings.Join(newTargetNames, ", ")
}
