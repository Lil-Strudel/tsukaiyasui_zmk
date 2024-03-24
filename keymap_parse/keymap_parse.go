package keymap_parse

import (
	"fmt"
	"os"
)

func ParseKeymap() {
	file, err := os.Open("zmk/corne.keymap")
	if err != nil {
		panic(err)
	}

	lexer := NewLexer(file)
	for {
		pos, tok, lit := lexer.Lex()
		if tok == EOF {
			break
		}

		fmt.Printf("%d:%d\t%s\t%s\n", pos.line, pos.column, tok, lit)
	}
}
