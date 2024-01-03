package test

import (
	deglob "alexander-scott/deglob/internal"
	"slices"
	"testing"
)

func TestScenario1(t *testing.T) {
	// Arrange
	path_to_workspace := "example_workspace/"
	path_to_build_file := path_to_workspace + "scenario_1/BUILD"
	path_to_expected_build_file := path_to_workspace + "scenario_1/BUILD_UPDATED"
	var expected_file_contents []string = deglob.LoadFileContentsIntoMemory(path_to_expected_build_file)

	// Act
	updated_file_contents := deglob.ProcessFile(path_to_build_file)

	// Assert
	if !slices.Equal(expected_file_contents, updated_file_contents) {
		t.Fatalf("Are not equivalent %s", updated_file_contents)
	}
}
