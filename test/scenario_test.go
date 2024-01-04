package test

import (
	"slices"
	"testing"

	deglob "github.com/alexander-scott/deglob/internal"
)

func TestScenarios(t *testing.T) {
	tests := []string{
		"scenario_1",
		"scenario_2",
	}

	const pathToWorkspace = "example_workspace/"
	const pathToBuildFileInScenario = "/BUILD"
	const pathToExpectedBuildFile = "/EXPECTED.BUILD"

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			// Arrange
			expectedFileContents := deglob.LoadFileContentsIntoMemory(pathToWorkspace + test + pathToExpectedBuildFile)

			// Act
			actualFileContents := deglob.ProcessFile(pathToWorkspace + test + pathToBuildFileInScenario)

			// Assert
			if !slices.Equal(expectedFileContents, actualFileContents) {
				t.Errorf("Slices are not equivalent.\nexpectedFileContents: \n%s\nactualFileContents: \n%s\n", difference(expectedFileContents, actualFileContents), difference(actualFileContents, expectedFileContents))
				t.Error(actualFileContents)
			}
		})
	}
}

func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
