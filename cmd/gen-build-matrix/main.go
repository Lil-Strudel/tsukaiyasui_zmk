package main

import (
	"flag"
	"fmt"

	"github.com/Lil-Strudel/tsukaiyasui_zmk/tsukaiyasui"
)

var shield, board string

func main() {
	flag.StringVar(&shield, "shield", shield, "the shield to build the build matrix for")
	flag.StringVar(&board, "board", board, "the board to build the build matrix for")
	flag.Parse()

	buildMatrix := tsukaiyasui.GenerateBuildMatrix(shield, board)
	fmt.Println(buildMatrix)
}
