package token

import (
	"fmt"

	"github.com/rhysd/locerr"
)

type Kind int

const (
	ILLEGAL Kind = iota
	BACKSLASH
	LAMBDA
	LPAREN
	RPAREN
	DOT
	SEMICOLON
	COLON_EQUAL
	EOF
)

var tokenTable = [...]string{
	ILLEGAL:     "ILLEGAL",
	BACKSLASH:   "\\",
	LAMBDA:      "λ",
	LPAREN:      "(",
	RPAREN:      ")",
	DOT:         ".",
	SEMICOLON:   ";",
	COLON_EQUAL: ":=",
	EOF:         "EOF",
}

type Token struct {
	Kind  Kind
	Start locerr.Pos
	End   locerr.Pos
	File  *locerr.Source
}

func (tok *Token) String() string {
	return fmt.Sprintf(
		"<%s:%s>(%d:%d:%d-%d:%d:%d)",
		tokenTable[tok.Kind],
		tok.Value(),
		tok.Start.Line, tok.Start.Column, tok.Start.Offset,
		tok.End.Line, tok.End.Column, tok.End.Offset)
}

func (tok *Token) Value() string {
	return string(tok.File.Code[tok.Start.Offset:tok.End.Offset])
}
