package ast

import (
	"fmt"

	"github.com/SuzumiyaAoba/fubuki/token"
)

type AST struct {
	Root []Expr
}

type Expr interface {
	Pos() token.Pos
	End() token.Pos
	Name() string
}

type (
	Var struct {
		Token  *token.Token
		Symbol string
	}
	Abs struct {
		StartToken *token.Token
		Var        *Var
		Body       Expr
	}
	App struct {
		Lexp Expr
		Rexp Expr
	}
	Def struct {
		StartToken *token.Token
		Symbol     string
		Bound      Expr
	}
)

func (e *Var) Pos() token.Pos {
	return e.Token.Start
}

func (e *Var) End() token.Pos {
	return e.Token.End
}

func (e *Abs) Pos() token.Pos {
	return e.StartToken.Start
}

func (e *Abs) End() token.Pos {
	return e.StartToken.End
}

func (e *App) Pos() token.Pos {
	return e.Lexp.Pos()
}

func (e *App) End() token.Pos {
	return e.Rexp.End()
}

func (e *Def) Pos() token.Pos {
	return e.StartToken.Start
}

func (e *Def) End() token.Pos {
	return e.StartToken.End
}

func (e *Var) Name() string {
	return fmt.Sprintf("Var (%s)", e.Symbol)
}

func (e *Abs) Name() string {
	return fmt.Sprintf("Abs (%s, %s)", e.Var.Name(), e.Body.Name())
}

func (e *App) Name() string {
	return fmt.Sprintf("App (%s %s)", e.Lexp.Name(), e.Rexp.Name())
}

func (e *Def) Name() string {
	return fmt.Sprintf("Def (%s, %s)", e.Symbol, e.Bound.Name())
}
