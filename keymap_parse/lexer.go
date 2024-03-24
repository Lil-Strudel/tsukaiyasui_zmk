// Thanks Aaron Raff and Thorsten Ball
// https://www.aaronraff.dev/blog/how-to-write-a-lexer-in-go
// https://interpreterbook.com/sample.pdf

package keymap_parse

import (
	"bufio"
	"io"
	"unicode"
)

type Token int

const (
	EOF = iota
	ILLEGAL

	IDENT
	STR
	PACMAN

	SEMI   // ;
	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }

	ASSIGN // =

	COMMENT // //

	INCLUDE
)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",

	// Identifiers
	IDENT:  "IDENT",
	STR:    "STR",
	PACMAN: "PACMAN",

	// Delimiters
	SEMI: ";",

	LPAREN: "(",
	RPAREN: ")",
	LBRACE: "{",
	RBRACE: "}",

	// Operators
	ASSIGN: "=",

	// Comment
	COMMENT: "//",

	// Keywords
	INCLUDE: "INCLUDE",
}

var keywords = map[string]Token{
	"#include": INCLUDE,
}

func LookupIdent(ident string) Token {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func isValidIdentRune(r rune) bool {
	validChars := map[rune]bool{
		'#': true,
		'/': true,
		'-': true,
		'_': true,
		'.': true,
		'&': true,
	}

	if unicode.IsLetter(r) || unicode.IsNumber(r) {
		return true
	}

	if _, ok := validChars[r]; ok {
		return true
	}

	return false
}

func (t Token) String() string {
	return tokens[t]
}

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func isEOF(err error) bool {
	if err != nil {
		if err == io.EOF {
			return true
		}
		// at this point there isn't much we can do, and the compiler
		// should just return the raw error to the user
		panic(err)
	}
	return false
}

func (lexer *Lexer) nextLine() {
	lexer.pos.line++
	lexer.pos.column = 0
}

func (lexer *Lexer) backup() {
	if err := lexer.reader.UnreadRune(); err != nil {
		panic(err)
	}

	lexer.pos.column--
}

// lexIdent scans the input until the end of an identifier and then returns the
// literal.
func (lexer *Lexer) lexIdent() string {
	var lit string
	for {
		rune, _, err := lexer.reader.ReadRune()
		if isEOF(err) {
			return lit
		}

		lexer.pos.column++
		if isValidIdentRune(rune) {
			lit = lit + string(rune)
		} else {
			// scanned something not in the identifier
			lexer.backup()
			return lit
		}
	}
}

func (lexer *Lexer) lexStr() string {
	var lit string
	for {
		rune, _, err := lexer.reader.ReadRune()
		if isEOF(err) || rune == '"' {
			return lit
		}

		lexer.pos.column++

		lit = lit + string(rune)
	}
}

func (lexer *Lexer) lexPacman() string {
	var lit string
	for {
		rune, _, err := lexer.reader.ReadRune()
		if isEOF(err) || rune == '>' {
			return lit
		}

		lexer.pos.column++

		lit = lit + string(rune)
	}
}

// lexComment scans the input until the end of the line and then returns the
// comment literal
func (lexer *Lexer) lexComment() string {
	commentLit := "//"
	for {
		rune, _, err := lexer.reader.ReadRune()
		if isEOF(err) || rune == '\n' {
			return commentLit
		}

		lexer.pos.column++

		commentLit = commentLit + string(rune)
	}
}

// lexComment scans the input until the end comment block marker and then
// returns the comment block literal
func (lexer *Lexer) lexCommentBlock() string {
	commentLit := "/*"
	for {
		rune, _, err := lexer.reader.ReadRune()
		if isEOF(err) {
			return commentLit
		}

		lexer.pos.column++

		if rune == '*' {
			rune, _, err := lexer.reader.ReadRune()
			if isEOF(err) {
				return commentLit
			}
			if rune == '/' {
				return commentLit + "*/"
			}
			lexer.backup()
		}

		commentLit = commentLit + string(rune)
	}
}

// Lex scans the input for the next token. It returns the position of the token,
// the token's type, and the literal value.
func (lexer *Lexer) Lex() (Position, Token, string) {
	// keep looping until we return a token
	for {
		rune, _, err := lexer.reader.ReadRune()
		if isEOF(err) {
			return lexer.pos, EOF, ""
		}

		// update the column to the position of the newly read in rune
		lexer.pos.column++

		switch rune {
		case '\n':
			lexer.nextLine()
		case ';':
			return lexer.pos, SEMI, ";"
		case '=':
			return lexer.pos, ASSIGN, "="
		case '(':
			return lexer.pos, LPAREN, "("
		case ')':
			return lexer.pos, RPAREN, ")"
		case '{':
			return lexer.pos, LBRACE, "{"
		case '}':
			return lexer.pos, RBRACE, "}"
		default:

			// handling comments
			if rune == '/' {
				startPos := lexer.pos
				lexer.pos.column++
				secondRune, _, err := lexer.reader.ReadRune()
				if isEOF(err) {
					// maybe this is not intended behavior
					return startPos, ILLEGAL, string(secondRune)
				}

				if secondRune == '/' {
					comment := lexer.lexComment()
					return startPos, COMMENT, comment
				} else if secondRune == '*' {
					comment := lexer.lexCommentBlock()
					return startPos, COMMENT, comment
				} else {
					lit := lexer.lexIdent()
					return startPos, IDENT, string(rune) + string(secondRune) + lit
				}
			}

			if unicode.IsSpace(rune) {
				continue // nothing to do here, just move on
			} else if rune == '"' {
				startPos := lexer.pos
				str := lexer.lexStr()
				return startPos, STR, str
			} else if rune == '<' {
				startPos := lexer.pos
				pacman := lexer.lexPacman()
				return startPos, PACMAN, pacman
			} else if isValidIdentRune(rune) {
				// backup and let lexIdent rescan the beginning of the ident
				startPos := lexer.pos
				lexer.backup()
				lit := lexer.lexIdent()

				identType := LookupIdent(lit)

				return startPos, identType, lit
			} else {
				return lexer.pos, ILLEGAL, string(rune)
			}
		}

	}
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}
