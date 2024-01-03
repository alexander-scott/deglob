package main

// 1) Find every file in a specified path
// 2) Filter files by pattern
// 3) For every filtered file, read the file contents
// 4) Find any globs in the files and the patterns they capture
// 5) Identify all the files in the bazel package that the glob captures
// 6) Save the structure of the bazel target with the glob into memory
// 7) Create new bazel targets based on what files the glob captured
// 8) Write the bazel targets to file
// 9) Replace the original bazel target with a new target that has deps to all of the new targets

import (
	deglob "alexander-scott/deglob/internal"
	"flag"
	"os"
)

func main() {
	// 0) Arg parse
	var workspace_path = flag.String("workspace_path", "", "REQUIRED: Path to workspace")
	var filter = flag.String("filter", ".*BUILD$", "File regex pattern")

	flag.Parse()

	if *workspace_path == "" {
		flag.Usage()
		os.Exit(1)
	}

	deglob.Run(*workspace_path, *filter)
}
