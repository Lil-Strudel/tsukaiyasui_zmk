package tsukaiyasui

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	keymap := GenerateKeymap("corne")
	fmt.Println(keymap)
	GenerateBuildMatrix("corne", "nice_nano_v2")
}
