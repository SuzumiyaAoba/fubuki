package main

import (
	"fmt"

	"github.com/SuzumiyaAoba/fubuki/syntax"
	"github.com/SuzumiyaAoba/fubuki/token"
)

func main() {
	src, _ := token.ReadFile("sample/bool.fbk")
	lexer := syntax.NewLexer(src)
	go lexer.Lex()

loop:
	for {
		select {
		case t := <-lexer.Tokens:
			switch t.Kind {
			case token.EOF:
				break loop
			default:
				fmt.Println(t.String())
			}
		}
	}

	fmt.Println(string(src.Code()))
}
