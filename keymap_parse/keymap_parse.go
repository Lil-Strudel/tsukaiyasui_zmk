package keymap_parse

import (
	"io"
	"strings"
)

type Layer struct {
	Order          int
	DisplayName    string
	Bindings       [][]string
	SensorBindings string
}

type KeymapNode struct {
	Compatable string
	Layers     map[string]Layer
}

type RootNode struct {
	Keymap KeymapNode
}

type Keymap struct {
	Includes []string
	RootNode RootNode
}

func ParseKeymap(file io.Reader) Keymap {
	lexer := NewLexer(file)
	tokens := []LexedToken{}
	for {
		pos, tok, lit := lexer.Lex()
		if tok == EOF {
			break
		}

		tokens = append(tokens, LexedToken{pos, tok, lit})
	}

	keymap := Keymap{}
	keymap.RootNode.Keymap.Layers = make(map[string]Layer)

	depth := make([]string, 0)
	for index := 0; index < len(tokens); index++ {
		token := tokens[index]
		if token.token == INCLUDE {
			nextToken := tokens[index+1]
			if nextToken.token != PACMAN {
				panic("Bad include")
			}
			keymap.Includes = append(keymap.Includes, nextToken.literal)
		}

		if token.token == IDENT {
			nextToken := tokens[index+1]
			if nextToken.token == LBRACE {
				depth = append(depth, token.literal)
			}

			if len(depth) == 2 {
				if token.literal == "compatible" {
					if nextToken.token == ASSIGN {
						nextNextToken := tokens[index+2]
						keymap.RootNode.Keymap.Compatable = nextNextToken.literal
					}
				}
			}

			if len(depth) == 3 {
				layerKey := depth[2]

				layer, ok := keymap.RootNode.Keymap.Layers[layerKey]
				if ok {
					if token.literal == "bindings" {
						if nextToken.token == ASSIGN {
							nextNextToken := tokens[index+2]

							bindings := nextNextToken.literal

							rows := strings.Split(strings.Trim(bindings, "\n "), "\n")
							keymapBindings := make([][]string, len(rows))
							for i := 0; i < len(rows); i++ {
								keymapBindings[i] = strings.Split(strings.Trim(rows[i], " ")[1:], "&")
								for j := 0; j < len(keymapBindings[i]); j++ {
									keymapBindings[i][j] = strings.Trim("&"+keymapBindings[i][j], " ")
								}
							}

							layer.Bindings = keymapBindings
						}
					}
					if token.literal == "sensor-bindings" {
						if nextToken.token == ASSIGN {
							nextNextToken := tokens[index+2]
							layer.SensorBindings = nextNextToken.literal
						}
					}
					if token.literal == "display-name" {
						if nextToken.token == ASSIGN {
							nextNextToken := tokens[index+2]
							layer.DisplayName = nextNextToken.literal
						}
					}

					keymap.RootNode.Keymap.Layers[layerKey] = layer
				} else {
					newLayer := Layer{}
					newLayer.Order = len(keymap.RootNode.Keymap.Layers)

					keymap.RootNode.Keymap.Layers[layerKey] = newLayer
				}

			}

		}

		if token.token == RBRACE {
			depth = depth[:len(depth)-1]
		}
	}

	return keymap
}
