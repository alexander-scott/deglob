package deglob

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func FindAllFilesInDirectory(path string) []string {
	// 1) Find every file in a specified path
	var all_files []string
	var walk = func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			all_files = append(all_files, s)
		}
		return nil
	}
	filepath.WalkDir(path, walk)

	return all_files
}

func LoadFileContentsIntoMemory(file_path string) []string {
	fmt.Println("Reading file: ", file_path)
	file_fs, err := os.OpenFile(file_path, os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file_fs.Close()

	var existing_file_contents []string

	scanner := bufio.NewScanner(file_fs)
	for scanner.Scan() {
		line := scanner.Text() // note, this will return a line of max 64K
		existing_file_contents = append(existing_file_contents, line+"\n")
	}

	return existing_file_contents
}

func WriteContentsToFile(file_path string, file_contents []string) {
	// 8) Write the bazel targets to file
	os.Truncate(file_path, 0) // Clear file contents
	file_fs, err := os.OpenFile(file_path, os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file_fs.Close()

	for _, line := range file_contents {
		file_fs.WriteString(line)
	}
}
