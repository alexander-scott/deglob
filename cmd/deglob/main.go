// Package main is the entry point for the deglob package
package main

import (
	deglob "alexander-scott/deglob/internal"
	"flag"
	"os"
)

func main() {
	// 0) Arg parse
	workspacePath := flag.String("workspace_path", "", "REQUIRED: Path to workspace")
	filter := flag.String("filter", ".*BUILD$", "File regex pattern")

	flag.Parse()

	if *workspacePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	deglob.Run(*workspacePath, *filter)
}
