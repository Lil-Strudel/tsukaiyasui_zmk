package tsukaiyasui

import (
	// "fmt"
	"embed"
	"maps"
	"slices"
	"strings"
)

//go:embed adapters/* zmk/*
var fs embed.FS

func GenerateKeymap(shield string) string {
	supportedShields := []string{"corne", "lily58"}

	found := slices.Contains(supportedShields, shield)
	if !found {
		panic("Unsupported shield")
	}

	keymap := make(map[string]string)

	maps.Copy(keymap, CoreKeymap)
	maps.Copy(keymap, leftKeymap)
	maps.Copy(keymap, rightKeymap)
	maps.Copy(keymap, thumbKeymap)
	maps.Copy(keymap, extraThumbKeymap)
	maps.Copy(keymap, specialKeymap)
	maps.Copy(keymap, numberRowKeymap)

	adapterBytes, _ := fs.ReadFile("adapters/" + shield + ".yasui")
	adapter := string(adapterBytes)

	rows := strings.Split(adapter, "\n")
	rows = rows[:len(rows)-1]

	longestKeymap := make(map[int]int)
	adaptedKeymap := make([][]string, len(rows))
	for rowI := range adaptedKeymap {
		adaptedKeymap[rowI] = strings.Split(rows[rowI], " ")
		for keyI, keyCode := range adaptedKeymap[rowI] {
			var thisKeymap string
			if keyCode == "___" {
				thisKeymap = ""
			} else {
				keymap, ok := keymap[keyCode]
				if ok {
					thisKeymap = keymap
				} else {
					thisKeymap = "&none"
				}
			}

			adaptedKeymap[rowI][keyI] = thisKeymap

			if longestKeymap[keyI] < len(thisKeymap) {
				longestKeymap[keyI] = len(thisKeymap)
			}
		}
	}

	joinedAdaptedKeymap := make([]string, len(adaptedKeymap))
	for rowI := range adaptedKeymap {
		for keyI, key := range adaptedKeymap[rowI] {
			spacesToAdd := longestKeymap[keyI] - len(key)
			if spacesToAdd > 0 {
				adaptedKeymap[rowI][keyI] = key + strings.Repeat(" ", spacesToAdd)
			}
		}

		joinedAdaptedKeymap[rowI] = strings.Join(adaptedKeymap[rowI], "  ")
	}

	keymapString := strings.Join(joinedAdaptedKeymap, "\n")

	templateKeymapBytes, _ := fs.ReadFile("zmk/template.keymap")
	templateKeymap := string(templateKeymapBytes)

	exportKeymap := strings.Replace(templateKeymap, "**yasui base**", keymapString, 1)

	return exportKeymap
}
