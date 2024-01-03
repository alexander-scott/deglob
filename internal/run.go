package deglob

// 1) Find every file in a specified path
// 2) Filter files by pattern
// 3) For every filtered file, read the file contents
// 4) Find any globs in the files and the patterns they capture
// 5) Identify all the files in the bazel package that the glob captures
// 6) Save the structure of the bazel target with the glob into memory
// 7) Create new bazel targets based on what files the glob captured
// 8) Write the bazel targets to file
// 9) Replace the original bazel target with a new target that has deps to all of the new targets

func Run(workspace_path string, filter string) {
	// 1) Find every file in a specified path
	var all_files = FindAllFilesInDirectory(workspace_path)

	// 2) Filter files by regex pattern
	var filtered_files = FilterPathsBasedOnRegexPattern(all_files, filter)

	// 3) For every filtered file, read the file contents
	for _, filtered_file := range filtered_files {
		var new_file_contents []string = ProcessFile(filtered_file)

		if new_file_contents != nil {
			WriteContentsToFile(filtered_file, new_file_contents)
		}
	}
}

func ProcessFile(filtered_file string) []string {
	var existing_file_contents []string = LoadFileContentsIntoMemory(filtered_file)
	var targets_with_glob []Target = ExtractTargetsFromFileContents(existing_file_contents, filtered_file)

	if len(targets_with_glob) > 0 {
		return CreateNewFileContentsIncludingNewTargets(existing_file_contents, targets_with_glob)
	}

	return nil
}
