package ast

import (
	"github.com/SuzumiyaAoba/fubuki/token"

	"github.com/rhysd/locerr"
)

type AST struct {
	Root Expr
}

type Expr interface {
	Pos() locerr.Pos
	End() locerr.Pos
	Name() string
}

type Symbol struct {
	Name string
	Id   int
}

func NewSymbol(name string) *Symbol {
	return &Symbol{name, 0}
}

type (
	Var struct {
		Token  *token.Token
		Symbol *Symbol
	}
	Abs struct {
		StartToken *token.Token
		Symbol     *Symbol
		Body       Expr
	}
	Apply struct {
		StartToken *token.Token
		Lexp       Expr
		Rexp       Expr
	}
	Def struct {
		StartToken *token.Token
		Symbol     *Symbol
		Bound      Expr
	}
)

func (e *Var) Pos() locerr.Pos {
	return e.Token.Start
}

func (e *Var) End() locerr.Pos {
	return e.Token.End
}

func (e *Abs) Pos() locerr.Pos {
	return e.StartToken.Start
}

func (e *Abs) End() locerr.Pos {
	return e.StartToken.End
}

func (e *Apply) Pos() locerr.Pos {
	return e.StartToken.Start
}

func (e *Apply) End() locerr.Pos {
	return e.StartToken.End
}

func (e *Def) Pos() locerr.Pos {
	return e.StartToken.Start
}

func (e *Def) End() locerr.Pos {
	return e.StartToken.End
}
