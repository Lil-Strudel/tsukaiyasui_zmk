package tsukaiyasui

import (
	"embed"
	"fmt"
	"maps"
	"slices"
	"strings"
)

//go:embed adapters/*
var fs embed.FS

func GenerateLayout(shield string) {
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

	adaptedKeymap := make([][]string, len(rows))
	for rowI := range rows {
		adaptedKeymap[rowI] = strings.Split(rows[rowI], " ")

		for keyI, part := range adaptedKeymap[rowI] {
			keymap, ok := keymap[part]
			if ok {
				adaptedKeymap[rowI][keyI] = keymap
			} else {
				adaptedKeymap[rowI][keyI] = "&none"
			}
		}

	}

	fmt.Println(adaptedKeymap)
}
