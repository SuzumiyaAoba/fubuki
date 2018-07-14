package syntax

import (
	"github.com/SuzumiyaAoba/fubuki/token"

	"github.com/SuzumiyaAoba/fubuki/ast"
	"github.com/rhysd/locerr"
)

type pseudoLexer struct {
	lastToken *token.Token
	tokens    chan token.Token
	err       *locerr.Error
	result    *ast.AST
}

func (l *pseudoLexer) Lex(lval *yySymType) int {
	for {
		select {
		case t := <-l.tokens:
			lval.token = &t
			switch t.Kind {
			case token.EOF, token.Illegal:
				return 0
			}

			l.lastToken = &t

			return int(t.Kind) + Illegal
		}
	}
}

func (l *pseudoLexer) Error(msg string) {
	if l.err == nil {
		if l.lastToken != nil {
			l.err = locerr.ErrorAt(l.lastToken.Start, msg)
		} else {
			l.err = locerr.NewError(msg)
		}
	} else {
		if l.lastToken != nil {
			l.err = l.err.NoteAt(l.lastToken.Start, msg)
		} else {
			l.err = l.err.Note(msg)
		}
	}
}

func Parse(src *locerr.Source) (*ast.AST, error) {
	var lexErr *locerr.Error
	l := NewLexer(src)
	l.Error = func(msg string, pos locerr.Pos) {
		if lexErr == nil {
			lexErr = locerr.ErrorAt(pos, msg)
		} else {
			lexErr = lexErr.NoteAt(pos, msg)
		}
	}
	go l.Lex()

	parsed, err := ParseTokens(l.Tokens)
	if lexErr != nil {
		return nil, lexErr.Note("Lexing source into tokens failed")
	}
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func ParseTokens(tokens chan token.Token) (*ast.AST, error) {
	yyErrorVerbose = true
	yyToknames = token.TokenTable

	l := &pseudoLexer{tokens: tokens}
	ret := yyParse(l)

	if l.err != nil {
		l.Error("Error while parsing")
		return nil, l.err
	}

	root := l.result
	if ret != 0 || root == nil {
		panic("FATAL: Parse failed for unknown reason")
	}

	return root, nil
}
