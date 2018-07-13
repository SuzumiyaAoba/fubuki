package ast

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Printer struct {
	indent int
	out    io.Writer
}

func (p Printer) VisitTopdown(e Expr) Visitor {
	fmt.Fprintf(p.out, "\n%s- %s (%d:%d-%d:%d)", strings.Repeat("  ", p.indent-1), e.Name(), e.Pos().Line, e.Pos().Column, e.End().Line, e.End().Column)
	return Printer{p.indent + 1, p.out}
}

func (p Printer) VisitBottomup(e Expr) {
	return
}

func Fprint(out io.Writer, a *AST) {
	p := Printer{1, out}
	for _, e := range a.Root {
		Visit(p, e)
	}
}

func Print(a *AST) {
	Fprint(os.Stdout, a)
}

func Println(a *AST) {
	Print(a)
	fmt.Println()
}
