package main

import (
	"flag"
	"fmt"

	"github.com/Lil-Strudel/tsukaiyasui_zmk/tsukaiyasui"
)

var shield string

func main() {
	flag.StringVar(&shield, "shield", shield, "the shield to build the layout for")
	flag.Parse()

	keymap := tsukaiyasui.GenerateKeymap(shield)
	fmt.Println(keymap)
}
