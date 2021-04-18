package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

type Lexer struct {
	src     []rune
	current int
	ast     Node
}

func NewLexer(code string) *Lexer {
	return &Lexer{
		src:     []rune(code + "\n"),
		current: 0,
	}
}

func (l *Lexer) Lex(lval *yySymType) int {
	for !l.isEof() {
		l.skipSpace()

		if isDigit(l.peek()) {
			return l.number(lval)
		}

		if l.peek() == ';' {
			l.consume()
			return LF
		}

		if l.peek() == '\n' {
			l.consume()
			return LF
		}

		if l.peek() == '[' {
			l.consume()
			return '['
		}

		if l.peek() == ']' {
			l.consume()
			return ']'
		}

		if l.peek() == '=' {
			l.consume()
			return '='
		}

		if l.peek() == '(' {
			l.consume()
			return '('
		}

		if l.peek() == ')' {
			l.consume()
			return ')'
		}

		if l.peek() == ',' {
			l.consume()
			return ','
		}

		if l.peek() == '"' || l.peek() == '\'' {
			return l.str(lval)
		}

		if isIdent(l.peek()) {
			return l.ident(lval)
		}
	}
	return -1
}

func (l *Lexer) Error(e string) {
	fmt.Println("[error] " + e)
	os.Exit(1)
}

func (l *Lexer) isEof() bool {
	return l.current >= len(l.src)
}

func (l *Lexer) skipSpace() {
	for l.peek() == ' ' || l.peek() == '\t' {
		l.consume()
	}
}

func (l *Lexer) peek() rune {
	return l.src[l.current]
}

func (l *Lexer) consume() rune {
	c := l.src[l.current]
	l.current++
	return c
}

func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

func isIdent(c rune) bool {
	return unicode.IsLetter(c) || c == '_' || unicode.IsDigit(c)
}

func (l *Lexer) ident(lval *yySymType) int {
	s := string(l.consume())

	for isIdent(l.peek()) {
		s += string(l.consume())
	}

	lval.ident = s
	return IDENT
}

func (l *Lexer) number(lval *yySymType) int {
	s := ""
	dotAppeared := false

	for !l.isEof() && isDigit(l.peek()) {
		c := l.consume()
		s += string(c)
		if l.peek() == '.' {
			if dotAppeared == true {
				panic("invalid number")
			}
			s += string(l.consume())
			dotAppeared = true
		}
	}
	lval.num, _ = strconv.ParseFloat(s, 64)
	return NUMBER
}

func (l *Lexer) str(lval *yySymType) int {
	l.consume()
	s := ""

	for c := l.consume(); c != '"' && c != '\''; c = l.consume() {
		if c == '\\' {
			c = l.consume()
		}
		s += string(c)
	}
	lval.str = s

	return STRING
}
