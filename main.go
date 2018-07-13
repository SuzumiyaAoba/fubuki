package main

import (
	"fmt"

	"github.com/SuzumiyaAoba/fubuki/lambda"
	"github.com/SuzumiyaAoba/fubuki/syntax"

	"github.com/rhysd/locerr"
)

func main() {
	n := "test.fbk"
	s, _ := locerr.NewSourceFromFile(n)
	t, _ := syntax.Parse(s)
	terms := lambda.AstToTerms(t)
	alpha := lambda.Alpha(terms)
	beta := lambda.Beta(alpha)

	for _, term := range beta {
		fmt.Println(lambda.Readable(term))
	}
}
