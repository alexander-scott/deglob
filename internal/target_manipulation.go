package deglob

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// GlobInfo holds a single glob search result together with the files it matched
type GlobInfo struct {
	searchResult GlobSearchResult
	globbedFiles []string
}

// Target is the main data structure for identifying targets with globs within files
type Target struct {
	start   int
	end     int
	name    string
	content []string
	globs   []GlobInfo
}

var (
	targetNamePattern  = regexp.MustCompile(`name\s=\s\"(?P<name>.*)\"`)
	targetStartPattern = regexp.MustCompile(`^cc_(library|binary)\(.*`)
	targetEndPattern   = regexp.MustCompile(`^\)\n$`)
)

// ExtractTargetsFromFileContents finds all targets with globs in a specific file and returns
// them as a slice
func ExtractTargetsFromFileContents(fileContents []string, filteredFile string) []Target {
	var targetsWithGlob []Target

	var currentTargetContent []string
	var currentTargetGlobResults []GlobSearchResult
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

		// Find any globs in the files and the patterns they capture
		if currentlyInTarget {
			checkLineForGlob := extractAllGlobPatternsFromLine(line)
			if checkLineForGlob.globFound {
				currentTargetGlobResults = append(currentTargetGlobResults, checkLineForGlob)
			}
		}

		if currentlyInTarget && targetEndPattern.MatchString(line) {
			// Only track target if the target had at least one glob
			if len(currentTargetGlobResults) > 0 {
				var globs []GlobInfo
				for _, globResult := range currentTargetGlobResults {
					files := findFilesFromGlob(globResult, filteredFile)
					globs = append(globs, GlobInfo{searchResult: globResult, globbedFiles: files})
				}
				target := Target{
					start:   currentTargetStartLineNumber,
					end:     currentLineNumber,
					name:    currentTargetName,
					content: append([]string(nil), currentTargetContent...),
					globs:   globs,
				}
				targetsWithGlob = append(targetsWithGlob, target)

				for _, glob := range globs {
					fmt.Println("Found glob in target: " + currentTargetName + " with attr " + glob.searchResult.globAttr)
				}
			}

			currentlyInTarget = false
			clear(currentTargetContent)
			currentTargetContent = nil
			currentTargetGlobResults = nil
		}

		currentLineNumber++
	}

	return targetsWithGlob
}

// CreateNewFileContentsIncludingNewTargets takes a slice of targets and the existing file contents
// and inserts the new de-globbed targets into the existing file and returns the entire file as a
// slice of strings
func CreateNewFileContentsIncludingNewTargets(existingFileContents []string, targetsWithGlob []Target) []string {
	var newFileContents []string

	currentlyInTarget := false
	var newTargetContent []string
	var currentTarget Target

	for index, line := range existingFileContents {
		var linesToAdd []string

		switch {
		// If we're currently in a target, our task is just to update several existing target attributes
		case currentlyInTarget:
			foundGlobIdx := -1
			for globIdx, glob := range currentTarget.globs {
				if line == glob.searchResult.fullLine {
					foundGlobIdx = globIdx
					break
				}
			}

			switch {
			case foundGlobIdx == 0:
				// First glob line: replace with deps listing all new sub-targets
				indent := extractIndent(line)
				linesToAdd = append(linesToAdd, indent+"deps = ["+createListOfAllNewTargetNames(currentTarget)+"],\n")
				// Preserve any explicit includes (e.g. glob([...]) + ["a.h"]) as a separate attribute
				for _, g := range currentTarget.globs {
					if len(g.searchResult.explicitIncludes) > 0 {
						quoted := make([]string, len(g.searchResult.explicitIncludes))
						for i, inc := range g.searchResult.explicitIncludes {
							quoted[i] = "\"" + inc + "\""
						}
						linesToAdd = append(linesToAdd, indent+g.searchResult.globAttr+" = ["+strings.Join(quoted, ", ")+"],\n")
					}
				}
			case foundGlobIdx > 0:
				// Subsequent glob lines in the same target: remove (all deps already listed above)
			default:
				linesToAdd = append(linesToAdd, line)
			}

			// The current target might require the generation of entirely new targets. If so, we
			// can't add it whilst we're still in the target as that will mess up the layout.
			// Therefore we add it to the `newTargetContent` slice so it can be written to file
			// once we're outside of the current target.
			if index == currentTarget.end {
				currentlyInTarget = false
				newTargetContent = createNewTargetsFromGlobbedFiles(currentTarget)
			}
		case len(newTargetContent) > 0:
			newFileContents = append(newFileContents, newTargetContent...)
			clear(newTargetContent)
			newTargetContent = nil
			linesToAdd = append(linesToAdd, line)
		case !currentlyInTarget && isCurrentLineNumberWithinAnyTarget(index, targetsWithGlob):
			currentlyInTarget = true
			currentTarget = returnTargetInCurrentLineNumber(index, targetsWithGlob)
			linesToAdd = append(linesToAdd, line)
		default:
			linesToAdd = append(linesToAdd, line)
		}

		newFileContents = append(newFileContents, linesToAdd...)
	}

	// If we've reached the end of the previous file, but we still have new targets to add, we can
	// add them now to the very bottom of the file
	if len(newTargetContent) > 0 {
		newFileContents = append(newFileContents, newTargetContent...)
		clear(newTargetContent)
	}

	return newFileContents
}

func findFilesFromGlob(searchResult GlobSearchResult, filteredFile string) []string {
	packagePath := strings.ReplaceAll(filteredFile, "BUILD", "")

	fmt.Println("Package_path: ", packagePath)
	fmt.Println("Glob attr: ", searchResult.globAttr)

	var globbedFiles []string

	for _, globPattern := range searchResult.globPatterns {
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
	for _, glob := range target.globs {
		for _, globbedFile := range glob.globbedFiles {
			newTargetContent = append(newTargetContent, "\n")
			for _, contentLine := range target.content {
				switch {
				case isOtherGlobLine(contentLine, glob, target.globs):
					// Skip other glob lines when building this sub-target
					continue
				case contentLine == glob.searchResult.fullLine:
					// Replace this glob line with the explicit source file reference
					indent := extractIndent(contentLine)
					newTargetContent = append(newTargetContent, indent+glob.searchResult.globAttr+" = [\""+globbedFile+"\"],\n")
				case targetNamePattern.MatchString(contentLine):
					// Replace the target name with the generated sub-target name
					newTargetName := generateNewTargetNameForGlobbedFile(target.name, globbedFile)
					newNameLine := strings.ReplaceAll(contentLine, target.name, newTargetName)
					newTargetContent = append(newTargetContent, newNameLine)
				case strings.HasPrefix(contentLine, "cc_binary("):
					// Sub-targets derived from a cc_binary are cc_library targets
					newTargetContent = append(newTargetContent, strings.Replace(contentLine, "cc_binary(", "cc_library(", 1))
				default:
					newTargetContent = append(newTargetContent, contentLine)
				}
			}
		}
	}
	return newTargetContent
}

func isAnyGlobLine(line string, globs []GlobInfo) bool {
	for _, glob := range globs {
		if line == glob.searchResult.fullLine {
			return true
		}
	}
	return false
}

// isOtherGlobLine returns true when line is a glob line belonging to a different glob than currentGlob.
func isOtherGlobLine(line string, currentGlob GlobInfo, allGlobs []GlobInfo) bool {
	return isAnyGlobLine(line, allGlobs) && line != currentGlob.searchResult.fullLine
}

func createListOfAllNewTargetNames(target Target) string {
	var newTargetNames []string
	for _, glob := range target.globs {
		for _, globbedFile := range glob.globbedFiles {
			newTargetNames = append(newTargetNames, "\":"+generateNewTargetNameForGlobbedFile(target.name, globbedFile)+"\"")
		}
	}
	return strings.Join(newTargetNames, ", ")
}

func extractIndent(line string) string {
	return line[:len(line)-len(strings.TrimLeft(line, " \t"))]
}

func generateNewTargetNameForGlobbedFile(targetName string, globbedFileName string) string {
	newNameSuffix := strings.ReplaceAll(globbedFileName, ".", "_")
	newNameSuffix = strings.ReplaceAll(newNameSuffix, "/", "_")
	return targetName + "_" + newNameSuffix
}
