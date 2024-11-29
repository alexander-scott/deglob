package deglob

import "strings"

func returnTargetNameFromLine(line string) string {
	matches := targetNamePattern.FindStringSubmatch(line)
	nameIndex := targetNamePattern.SubexpIndex("name")
	return matches[nameIndex]
}

func createListOfNewTargetNamesFromTarget(t target) string {
	var newTargetNames []string
	for _, targetGlobbedFile := range t.globbedFiles {
		for _, targetContentLine := range t.content {
			if targetNamePattern.MatchString(targetContentLine) {
				newTargetName := generateNewTargetNameForGlobbedFile(t.name, targetGlobbedFile, true, true)
				newTargetNames = append(newTargetNames, newTargetName)
			}
		}
	}
	return strings.Join(newTargetNames, ", ")
}

func generateNewTargetNameForGlobbedFile(targetName string, globbedFileName string, asLabel bool, wrapWithQuotes bool) string {
	newNameSuffix := strings.Split(globbedFileName, ".")[0]
	newNameSuffix = strings.ReplaceAll(newNameSuffix, "/", "_")
	newTargetName := targetName + "_" + newNameSuffix
	if asLabel {
		newTargetName = ":" + newTargetName
	}
	if wrapWithQuotes {
		newTargetName = "\"" + newTargetName + "\""
	}
	return newTargetName
}
