package test

import (
	deglob "alexander-scott/deglob/internal"
	"slices"
	"testing"
)

func TestScenario1(t *testing.T) {
	// Arrange
	pathToWorkspace := "example_workspace/"
	pathToBuildFile := pathToWorkspace + "scenario_1/BUILD"
	pathToExpectedBuildFile := pathToWorkspace + "scenario_1/BUILD_UPDATED"
	expectedFileContents := deglob.LoadFileContentsIntoMemory(pathToExpectedBuildFile)

	// Act
	updatedFileContents := deglob.ProcessFile(pathToBuildFile)

	// Assert
	if !slices.Equal(expectedFileContents, updatedFileContents) {
		t.Fatalf("Are not equivalent %s", updatedFileContents)
	}
}
