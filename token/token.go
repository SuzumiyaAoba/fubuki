package token

import (
	"fmt"

	"github.com/rhysd/locerr"
)

type Kind int

const (
	Illegal Kind = iota
	Ident
	Lambda
	LParen
	RParen
	Dot
	Semicolon
	ColonEqual
	EOF
)

var TokenTable = [...]string{
	"EOF",
	"Error",
	"Unknown",
	"illegal token",
	"ident",
	"λ",
	"(",
	")",
	".",
	";",
	":=",
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
		TokenTable[tok.Kind],
		tok.Value(),
		tok.Start.Line, tok.Start.Column, tok.Start.Offset,
		tok.End.Line, tok.End.Column, tok.End.Offset)
}

func (tok *Token) Value() string {
	return string(tok.File.Code[tok.Start.Offset:tok.End.Offset])
}
