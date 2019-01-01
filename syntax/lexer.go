package syntax

import (
	"bytes"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"

	"github.com/SuzumiyaAoba/fubuki/token"
	"github.com/SuzumiyaAoba/fubuki/util"
)

type stateFn func(*Lexer) stateFn

const eof = -1

type Lexer struct {
	state   stateFn
	start   token.Pos
	current token.Pos
	src     token.Source
	input   *bytes.Reader
	Tokens  chan token.Token
	top     rune
	eof     bool
	Error   func(msg string, pos token.Pos)
}

func NewLexer(src token.Source) *Lexer {
	start := token.Pos{
		Offset: 0,
		Line:   1,
		Column: 1,
		Src:    src,
	}

	return &Lexer{
		state:   lex,
		start:   start,
		current: start,
		input:   bytes.NewReader(src.Code()),
		src:     src,
		Tokens:  make(chan token.Token),
		Error:   nil,
	}
}

func (l *Lexer) Lex() {
	l.forward()
	for l.state != nil {
		l.state = l.state(l)
	}
}

func (l *Lexer) emit(kind token.Kind) {
	l.Tokens <- token.Token{
		kind,
		l.start,
		l.current,
		l.src,
	}
	l.start = l.current
}

func (l *Lexer) emitIllegal(reason string) {
	l.errmsg(reason)
	t := token.Token{
		token.Illegal,
		l.start,
		l.current,
		l.src,
	}
	l.Tokens <- t
	l.start = l.current
}

func (l *Lexer) expected(s string, actual rune) {
	l.emitIllegal(fmt.Sprintf("Expected %s but got '%c'(%d)", s, actual, actual))
}

func (l *Lexer) forward() {
	r, _, err := l.input.ReadRune()
	if err == io.EOF {
		l.top = 0
		l.eof = true
		return
	}

	if err != nil {
		panic(err)
	}

	if !utf8.ValidRune(r) {
		l.emitIllegal(fmt.Sprintf("Invalid UTF-8 character '%c' (%d)", r, r))
		l.eof = true
		return
	}

	l.top = r
	l.eof = false
}

func (l *Lexer) eat() {
	size := utf8.RuneLen(l.top)
	l.current.Offset += size

	if l.top == '\n' {
		l.current.Line++
		l.current.Column = 1
		l.forward()
		return
	}

	if l.top == '\r' {
		l.forward()
		size = utf8.RuneLen(l.top)
		l.current.Offset += size
		if l.top == '\n' {
			l.forward()
		}
		l.current.Line++
		l.current.Column = 1
		return
	}

	l.current.Column += size
	l.forward()
}

func (l *Lexer) consume() {
	if l.eof {
		return
	}
	l.eat()
	l.start = l.current
}

func (l *Lexer) errmsg(msg string) {
	if l.Error == nil {
		return
	}
	l.Error(msg, l.current)
}

func eatIdent(l *Lexer) bool {
	if !(util.IsLetter(l.top) || util.IsDigit(l.top) || l.top == '#') {
		l.expected("letter for head character of identifer", l.top)
		return false
	}
	l.eat()

	for util.IsLetter(l.top) || util.IsDigit(l.top) {
		l.eat()
	}

	return true
}

func lexIdent(l *Lexer) stateFn {
	if !eatIdent(l) {
		return nil
	}

	l.emit(token.Ident)

	return lex
}

func lexComment(l *Lexer) stateFn {
	for {
		if l.eof {
			return lex
		}
		line := l.current.Line
		l.eat()
		if l.current.Line > line {
			return lex
		}
		l.eat()
	}
}

func lexHyphenMinus(l *Lexer) stateFn {
	prev := l.top
	l.eat()

	if prev != l.top {
		l.expected("comment --", l.top)
		return nil
	}

	for l.top != '\n' && !l.eof {
		l.eat()
	}

	return lex
}

func lexColon(l *Lexer) stateFn {
	l.eat()

	if l.top != '=' {
		l.expected("define :=", l.top)
		return nil
	}

	l.eat()
	l.emit(token.ColonEqual)

	return lex

}

func lex(l *Lexer) stateFn {
	for {
		if l.eof {
			l.emit(token.EOF)
			return nil
		}
		switch l.top {
		case 'λ', '\\':
			l.eat()
			l.emit(token.Lambda)
		case '.':
			l.eat()
			l.emit(token.Dot)
		case '(':
			l.eat()
			l.emit(token.LParen)
		case ')':
			l.eat()
			l.emit(token.RParen)
		case ';':
			l.eat()
			l.emit(token.Semicolon)
		case ':':
			return lexColon
		case '-':
			return lexHyphenMinus
		default:
			switch {
			case unicode.IsSpace(l.top):
				l.consume()
			default:
				return lexIdent
			}
		}
	}
}
