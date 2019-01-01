package token

import (
	"fmt"
)

// Kind は字句の種類
type Kind int

const (
	// Illegal は不正な字句
	Illegal Kind = iota
	// Ident は識別子
	Ident
	// Lambda は `λ`
	Lambda
	// LParen は `(`
	LParen
	// RParen は `)`
	RParen
	// Dot は `.`
	Dot
	// Semicolon は `;`
	Semicolon
	// ColonEqual は `:=`
	ColonEqual
	// EOF は終端文字
	EOF
)

// TokenTable は各字句とそれが表す文字列の対応表
var TokenTable = [...]string{
	"illegal token",
	"ident",
	"λ",
	"(",
	")",
	".",
	";",
	":=",
	"EOF",
	"Error",
	"Unknown",
}

type Token struct {
	Kind  Kind
	Start Pos
	End   Pos
	Src   Source
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
	code := tok.Src.Code()
	return string(code[tok.Start.Offset:tok.End.Offset])
}
