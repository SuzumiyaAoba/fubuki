package lambda

import (
	"fmt"

	"github.com/SuzumiyaAoba/fubuki/ast"
)

type Term interface {
	Pretty() string
}

type (
	Var struct {
		Symbol string
		ID     int
	}
	Abs struct {
		Var  *Var
		Body Term
	}
	App struct {
		Lterm Term
		Rterm Term
	}
	Def struct {
		Symbol string
		Bound  Term
	}
)

func (t *Var) Pretty() string {
	return fmt.Sprintf("%s", t.Symbol)
}

func (t *Abs) Pretty() string {
	return fmt.Sprintf("(λ%s.%s)", t.Var.Symbol, t.Body.Pretty())
}

func (t *App) Pretty() string {
	return fmt.Sprintf("(%s %s)", t.Lterm.Pretty(), t.Rterm.Pretty())
}

func (t *Def) Pretty() string {
	return fmt.Sprintf("%s := %s", t.Symbol, t.Bound.Pretty())
}

func ExprToTerm(expr ast.Expr) Term {
	switch e := expr.(type) {
	case *ast.Var:
		return &Var{e.Symbol, 0}
	case *ast.Abs:
		return &Abs{&Var{e.Var.Symbol, 0}, ExprToTerm(e.Body)}
	case *ast.App:
		return &App{ExprToTerm(e.Lexp), ExprToTerm(e.Rexp)}
	case *ast.Def:
		return &Def{e.Symbol, ExprToTerm(e.Bound)}
	}
	panic("unknown expr")
}

func AstToTerms(tree *ast.AST) []Term {
	terms := make([]Term, 0)

	for _, e := range tree.Root {
		terms = append([]Term{ExprToTerm(e)}, terms...)
	}

	for i, j := 0, len(terms)-1; i < j; i, j = i+1, j-1 {
		terms[i], terms[j] = terms[j], terms[i]
	}

	return terms
}

func Readable(term Term) string {
	switch t := term.(type) {
	case *Var:
		return t.Pretty()
	case *Abs:
		if _, ok := t.Body.(*Abs); ok {
			return fmt.Sprintf("(λ%s%s)", t.Var.Pretty(), readable(t.Body))
		}
		return fmt.Sprintf("(λ%s.%s)", t.Var.Pretty(), readable(t.Body))
	case *App:
		return fmt.Sprintf("(%s %s)", Readable(t.Lterm), Readable(t.Rterm))
	case *Def:
		return fmt.Sprintf("%s := %s", t.Symbol, Readable(t.Bound))
	}
	panic("unknown term")
}

func readable(term Term) string {
	switch t := term.(type) {
	case *Abs:
		if _, ok := t.Body.(*Abs); ok {
			return fmt.Sprintf(" %s%s", t.Var.Pretty(), readable(t.Body))
		}
		return fmt.Sprintf(" %s.%s", t.Var.Pretty(), readable(t.Body))
	case *App:
		return fmt.Sprintf("%s %s", Readable(t.Lterm), Readable(t.Rterm))
	}
	return Readable(term)
}
