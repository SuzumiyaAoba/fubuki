package main

import (
	"github.com/SuzumiyaAoba/fubuki/ast"
	"github.com/SuzumiyaAoba/fubuki/syntax"

	"github.com/rhysd/locerr"
)

func main() {
	n := "test.fbk"
	s, _ := locerr.NewSourceFromFile(n)
	t, _ := syntax.Parse(s)
	ast.Println(t)
}
