package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
			assert.Equal(t, expectedFileContents, actualFileContents)
		})
	}
}
