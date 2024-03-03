package tsukaiyasui

import (
	"slices"
	"strings"
)

func extractYasuiSection(contents, section string) string {
	startTag := "/-" + section
	endTag := "-/"

	startIndex := strings.Index(contents, startTag)
	if startIndex == -1 {
		panic("Improperly formatted Yasui file")
	}

	startIndex += len(startTag) + 1

	endIndex := strings.Index(contents[startIndex:], endTag)
	if endIndex == -1 {
		panic("Improperly formatted Yasui file")
	}

	endIndex += startIndex - 1
	return contents[startIndex:endIndex]
}

func validateShield(shield string) {
	adapters, _ := fs.ReadDir("adapters")
	supportedShields := make([]string, len(adapters))

	for i, adapter := range adapters {
		name, _ := strings.CutSuffix(adapter.Name(), ".yasui")
		supportedShields[i] = name
	}

	found := slices.Contains(supportedShields, shield)
	if !found {
		panic("Unsupported shield")
	}
}
