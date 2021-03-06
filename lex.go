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
	if l.isEof() {
		return -1
	}

	l.skipSpace()

	if l.consumeIf('#') {
		l.skipComment()
	}

	if isDigit(l.peek()) {
		return l.number(lval)
	}

	if l.consumeIf(';') {
		return LF
	}

	if l.consumeIf('\n') {
		return LF
	}

	if l.consumeIf('\r') {
		if l.consumeIf('\n') {
			return LF
		}
	}

	if l.consumeIf('[') {
		return '['
	}

	if l.consumeIf(']') {
		return ']'
	}

	if l.consumeIf('{') {
		return '{'
	}

	if l.consumeIf('}') {
		return '}'
	}

	if l.consumeIf('+') {
		if l.consumeIf('=') {
			return ADD_ASSIGN
		}
		if l.consumeIf('+') {
			return INC
		}
		return '+'
	}

	if l.consumeIf('-') {
		if l.consumeIf('=') {
			return SUB_ASSIGN
		}
		if l.consumeIf('-') {
			return DEC
		}
		return '-'
	}

	if l.consumeIf('*') {
		if l.consumeIf('*') {
			if l.consumeIf('=') {
				return POW_ASSIGN
			}
			return POW
		}
		if l.consumeIf('=') {
			return MUL_ASSIGN
		}
		return '*'
	}

	if l.consumeIf('/') {
		if l.consumeIf('=') {
			return DIV_ASSIGN
		}
		return '/'
	}

	if l.consumeIf('%') {
		if l.consumeIf('=') {
			return MOD_ASSIGN
		}
		return '%'
	}

	if l.consumeIf('=') {
		if l.consumeIf('=') {
			return NUMEQ
		}
		return '='
	}

	if l.consumeIf('!') {
		if l.consumeIf('=') {
			return NUMNE
		}
		if l.consumeIf('~') {
			return NOT_MATCH
		}
		return '!'
	}

	if l.consumeIf('<') {
		if l.consumeIf('=') {
			return NUMLE
		}
		return '<'
	}

	if l.consumeIf('>') {
		if l.consumeIf('=') {
			return NUMGE
		}
		return '>'
	}

	if l.consumeIf('&') {
		if l.consumeIf('&') {
			return AND
		}
		panic("invalid char &")
	}

	if l.consumeIf('|') {
		if l.consumeIf('|') {
			return OR
		}
		panic("invalid char |")
	}

	if l.consumeIf('~') {
		return '~'
	}

	if l.consumeIf('.') {
		if l.consumeIf('=') {
			return CONCAT_ASSIGN
		}
		return '.'
	}

	if l.consumeIf('(') {
		return '('
	}

	if l.consumeIf(')') {
		return ')'
	}

	if l.consumeIf(',') {
		return ','
	}

	if l.peek() == '"' {
		return l.doubleQuoteStr(lval)
	}
	if l.peek() == '\'' {
		return l.singleQuoteStr(lval)
	}

	if isIdent(l.peek()) {
		return l.word(lval)
	}

	fatalError("syntax error: %c", l.peek())

	return -1
}

func (l *Lexer) Error(e string) {
	fmt.Fprintf(execContext.errout, "[error] %s\n", e)
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

func (l *Lexer) skipComment() {
	for l.peek() != '\n' {
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

func (l *Lexer) consumeIf(r rune) bool {
	if l.peek() == r {
		l.consume()
		return true
	}
	return false
}

func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

func isIdent(c rune) bool {
	return unicode.IsLetter(c) || c == '_' || unicode.IsDigit(c) || c == '@' || c == '$'
}

func (l *Lexer) word(lval *yySymType) int {
	s := string(l.consume())

	for isIdent(l.peek()) {
		s += string(l.consume())
	}

	if s == "if" {
		return IF
	}

	if s == "else" {
		return ELSE
	}

	if s == "eq" {
		return STREQ
	}

	if s == "ne" {
		return STRNE
	}

	if s == "lt" {
		return COLLT
	}
	if s == "le" {
		return COLLE
	}

	if s == "gt" {
		return COLGT
	}
	if s == "ge" {
		return COLGE
	}

	if s == "while" {
		return WHILE
	}

	if s == "do" {
		return DO
	}

	if s == "for" {
		return FOR
	}

	if s == "break" {
		return BREAK
	}

	if s == "continue" {
		return CONTINUE
	}

	if s == "function" {
		return FUNCTION
	}

	if s == "return" {
		return RETURN
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

func (l *Lexer) doubleQuoteStr(lval *yySymType) int {
	l.consume()
	s := ""

	for c := l.consume(); c != '"' && !l.isEof(); c = l.consume() {
		if c == '\\' {
			c = l.consumeEscapeChar()
		}
		s += string(c)
	}
	lval.str = s

	return STRING
}

func (l *Lexer) singleQuoteStr(lval *yySymType) int {
	l.consume()
	s := ""

	for c := l.consume(); c != '\'' && !l.isEof(); c = l.consume() {
		if c == '\\' {
			c = l.consumeEscapeChar()
		}
		s += string(c)
	}
	lval.str = s

	return STRING
}

func (l *Lexer) consumeEscapeChar() rune {
	c := l.consume()
	if c == 'a' {
		return '\a'
	}
	if c == 'b' {
		return '\b'
	}
	if c == 'f' {
		return '\f'
	}
	if c == 'n' {
		return '\n'
	}
	if c == 'r' {
		return '\r'
	}
	if c == 't' {
		return '\t'
	}
	if c == 'v' {
		return '\v'
	}
	if c == '\\' {
		return '\\'
	}
	if c == '"' {
		return '"'
	}
	if c == '\'' {
		return '\''
	}

	panic("unknown escape")
}
