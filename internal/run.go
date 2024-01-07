// Package deglob allows you to remove all globs from files in a specified workspace path
// 1) Find every file in a specified path
// 2) Filter files by pattern
// 3) For every filtered file, read the file contents
// 4) Find any globs in the files and the patterns they capture
// 5) Identify all the files in the bazel package that the glob captures
// 6) Save the structure of the bazel target with the glob into memory
// 7) Create new bazel targets based on what files the glob captured
// 8) Write the bazel targets to file
// 9) Replace the original bazel target with a new target that has deps to all of the new targets
package deglob

// Run performs the main deglob logic, as specified in the package header
func Run(workspacePath string, filter string) {
	// 1) Find every file in a specified path
	allFiles := findAllFilesInDirectory(workspacePath)

	// 2) Filter files by regex pattern
	filteredFiles := filterPathsBasedOnRegexPattern(allFiles, filter)

	// 3) For every filtered file, read the file contents
	for _, filteredFile := range filteredFiles {
		newFileContents := ProcessFile(filteredFile)

		if newFileContents != nil {
			writeContentsToFile(filteredFile, newFileContents)
		}
	}
}

// ProcessFile performs the deglob magic on a file at a specific path.
// It reads the file contents, identifies all globs, and then produces the file
// with the removal of the globs.
func ProcessFile(filteredFile string) []string {
	existingFileContents := LoadFileContentsIntoMemory(filteredFile)
	targetsWithGlob := extractTargetsFromFileContents(existingFileContents, filteredFile)

	if len(targetsWithGlob) > 0 {
		return createNewFileContentsIncludingNewTargets(existingFileContents, targetsWithGlob)
	}

	return nil
}
