package keymap_parse

import (
	"io"
)

type Layer struct {
	Order          int
	Label          string
	Bindings       string
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
							layer.Bindings = nextNextToken.literal
						}
					}
					if token.literal == "sensor-bindings" {
						if nextToken.token == ASSIGN {
							nextNextToken := tokens[index+2]
							layer.SensorBindings = nextNextToken.literal
						}
					}
					if token.literal == "label" {
						if nextToken.token == ASSIGN {
							nextNextToken := tokens[index+2]
							layer.Label = nextNextToken.literal
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
