%{
package syntax

import (
        "github.com/SuzumiyaAoba/fubuki/ast"
        "github.com/SuzumiyaAoba/fubuki/token"
)

func makeAbs(idents []*ast.Var, body ast.Expr) ast.Expr {
	if len(idents) == 0 {
		return body
	}
	ident := idents[0]
	return &ast.Abs{ident.Token, ident, makeAbs(idents[1:], body)}
}
%}

%token<token> Illegal
%token<token> Ident
%token<token> Lambda
%token<token> LParen
%token<token> RParen
%token<token> Dot
%token<token> Semicolon
%token<token> ColonEqual

%type<> program
%type<nodes> stmts
%type<node> stmt
%type<node> define
%type<node> expr
%type<node> term
%type<variable> variable
%type<variables> variables

%union {
  token     *token.Token
  node      ast.Expr
  nodes     []ast.Expr
  variable  *ast.Var
  variables []*ast.Var
}

%start program

%%

program: stmts {
           tree := &ast.AST{$1}
           yylex.(*pseudoLexer).result = tree
         }

stmts:                      { $$ = []ast.Expr{} }
     | stmt                 { $$ = []ast.Expr{$1} }
     | stmt Semicolon stmts { $$ = append([]ast.Expr{$1}, $3...) }

stmt: define { $$ = $1 }
    | expr   {$$ = $1; }

define: Ident ColonEqual expr { $$ = &ast.Def{$1, $1.Value(), $3} }

expr: term { $$ = $1 }
    | expr term { $$ = &ast.App{$1, $2} }

term: variable { $$ = $1 }
    | Lambda variables Dot expr {
      $$ = &ast.Abs{$1, $2[0], makeAbs($2[1:], $4)}
    }
    | LParen expr RParen { $$ = $2 }

variables:
  variable {
    $$ = []*ast.Var{$1}
  }
  | variables variable {
    $$ = append($1, $2)
  }

variable: Ident { $$ = &ast.Var{$1, $1.Value()} }

