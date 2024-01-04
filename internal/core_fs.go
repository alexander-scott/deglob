package deglob

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// FindAllFilesInDirectory recursively finds all files within a directory
func FindAllFilesInDirectory(path string) []string {
	var allFiles []string
	walk := func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			allFiles = append(allFiles, s)
		}
		return nil
	}
	err := filepath.WalkDir(path, walk)
	if err != nil {
		panic(err)
	}

	return allFiles
}

// LoadFileContentsIntoMemory reads all lines from a specified file path into memory and
// returns it as a slice.
func LoadFileContentsIntoMemory(filePath string) []string {
	fmt.Println("Reading file: ", filePath)
	fileFs, err := os.OpenFile(filepath.Clean(filePath), os.O_RDONLY, 0o600)
	if err != nil {
		panic(err)
	}

	var existingFileContents []string

	scanner := bufio.NewScanner(fileFs)
	for scanner.Scan() {
		line := scanner.Text() // note, this will return a line of max 64K
		existingFileContents = append(existingFileContents, line+"\n")
	}

	err = fileFs.Close()
	if err != nil {
		panic(err)
	}

	return existingFileContents
}

// WriteContentsToFile writes a slice of strings into a specific filePath
func WriteContentsToFile(filePath string, fileContents []string) {
	err := os.Truncate(filePath, 0) // Clear existing file contents
	if err != nil {
		panic(err)
	}

	fileFs, err := os.OpenFile(filepath.Clean(filePath), os.O_WRONLY, 0o600)
	if err != nil {
		panic(err)
	}

	for _, line := range fileContents {
		_, writeErr := fileFs.WriteString(line)
		if writeErr != nil {
			panic(writeErr)
		}
	}

	err = fileFs.Close()
	if err != nil {
		panic(err)
	}
}
