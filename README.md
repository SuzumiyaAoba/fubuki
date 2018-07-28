# Fubuki :cat:

Fubuki is an implementation of λ-calculus interpreter in Golang.

## Installation

```
$ go get github.com/SuzumiyaAoba/fubuki
```

## Grammar

The grammar of λ-expressions in Fubuki is following (not correct, but enough).

```
<expr>  := <term> | <def>
<term>  := <var> | <abs> | <app> | '('<term>')'
<ident> := [#0-9a-zA-Z_][0-9a-zA-Z_]*
<var>   := <ident>
<vars>  := <var> <vars>*
<abs>   := '\'<vars>'.'<term>
<app>   := <term> <term>
<def>   := <ident> ":=" <term>
```

## Evaluation strategy
rightmost-innermost.

## Usage

```
$ fubuki
```

### Command

Following commands can be used in REPL.

```
:exit                 exit REPL
:load [<path>...]     load file              (short :l)
:env [asc] [desc] [#] show environment
:show [<name>]        show lambda expression (short :s)
:help                 show help              (short :h)
```

## Example

```
$ cat sample/nat.fbk
0 := \f x. x;
1 := \f x. f x;
2 := \f x. f (f x);
3 := \f x. f (f (f x));

succ := \n f x. f (n f x);
plus := \m n f x. m f (n f x);
mult := \m n. m (plus n) 0;
pred := \n f x. n (\g h. h (g f)) (\u. x) (\u. u);
sub := \m n. n pred m;

$ cat sample/bool.fbk
true := \x y. x;
false := \x y. y;

if := \p x y. p x y;

and := \p q. p q false;
or := \p q. p true q;
not := \p. p false true;

$ fubuki
Welcome to Fubuki 0.0.1
see https://github.com/SuzumiyaAoba/fubuki :help for help

fubuki> \x y. x;
#0: (λx y.x)

fubuki> #0 a
#1: (λy.a)

fubuki> #0 a b
#2: a

fubuki> :load sample/nat.fbk sample/bool.fbk
Success: load sample/nat.fbk
Success: load sample/bool.fbk

fubuki> :env
pred := (λn f x.((n (λg h.h (g f))) (λu.x)) (λu.u))
or := (λp q.(p true) q)
not := (λp.(p false) true)
0 := (λf x.x)
succ := (λn f x.f ((n f) x))
3 := (λf x.f (f (f x)))
mult := (λm n.(m (plus n)) 0)
and := (λp q.(p q) false)
1 := (λf x.f x)
plus := (λm n f x.(m f) ((n f) x))
true := (λx y.x)
2 := (λf x.f (f x))
false := (λx y.y)
if := (λp x y.(p x) y)
sub := (λm n.(n pred) m)

fubuki> :env asc
0 := (λf x.x)
1 := (λf x.f x)
2 := (λf x.f (f x))
3 := (λf x.f (f (f x)))
and := (λp q.(p q) false)
false := (λx y.y)
if := (λp x y.(p x) y)
mult := (λm n.(m (plus n)) 0)
not := (λp.(p false) true)
or := (λp q.(p true) q)
plus := (λm n f x.(m f) ((n f) x))
pred := (λn f x.((n (λg h.h (g f))) (λu.x)) (λu.u))
sub := (λm n.(n pred) m)
succ := (λn f x.f ((n f) x))
true := (λx y.x)

fubuki> :env # asc
#0 := (λx y.x)
#1 := (λy.a)
#2 := a
0 := (λf x.x)
1 := (λf x.f x)
2 := (λf x.f (f x))
3 := (λf x.f (f (f x)))
and := (λp q.(p q) false)
false := (λx y.y)
if := (λp x y.(p x) y)
mult := (λm n.(m (plus n)) 0)
not := (λp.(p false) true)
or := (λp q.(p true) q)
plus := (λm n f x.(m f) ((n f) x))
pred := (λn f x.((n (λg h.h (g f))) (λu.x)) (λu.u))
sub := (λm n.(n pred) m)
succ := (λn f x.f ((n f) x))
true := (λx y.x)

fubuki> :show hoge
Not found: hoge

fubuki> :show true
Exists: true := (λx y.x)
```
