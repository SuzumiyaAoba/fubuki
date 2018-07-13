package ast

import (
	"fmt"

	"github.com/SuzumiyaAoba/fubuki/token"

	"github.com/rhysd/locerr"
)

type AST struct {
	Root []Expr
}

type Expr interface {
	Pos() locerr.Pos
	End() locerr.Pos
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

func (e *App) Pos() locerr.Pos {
	return e.Lexp.Pos()
}

func (e *App) End() locerr.Pos {
	return e.Rexp.End()
}

func (e *Def) Pos() locerr.Pos {
	return e.StartToken.Start
}

func (e *Def) End() locerr.Pos {
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
