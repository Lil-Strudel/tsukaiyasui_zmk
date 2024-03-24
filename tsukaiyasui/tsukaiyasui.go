package tsukaiyasui

import (
	"embed"
	"encoding/json"
	"strings"
)

//go:embed adapters/* zmk/*
var fs embed.FS

func GenerateKeymap(shield string) string {
	validateShield(shield)
	return "123"
}

type build struct {
	Board  string `json:"board"`
	Shield string `json:"shield"`
}

type buildMatrix struct {
	Include []build `json:"include"`
}

func GenerateBuildMatrix(shield, board string) string {
	validateShield(shield)

	adapterBytes, _ := fs.ReadFile("adapters/" + shield + ".yasui")
	adapterShields := extractYasuiSection(string(adapterBytes), "shields")

	shields := strings.Split(adapterShields, "\n")

	layout := make([]build, len(shields))
	for i, shield := range shields {
		layout[i] = build{board, shield}
	}

	buildMatrix := buildMatrix{layout}

	jsonBuildMatrix, _ := json.Marshal(buildMatrix)

	return (string(jsonBuildMatrix))
}
