package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"utf8"
)

type itemType int

const eof = 0

const (
	itemEOF        itemType = iota
	itemIdentifier
	itemNumber
)

type stateFn func (*calcLex) stateFn

type calcLex struct {
	s     string           // the string being scanned
	pos   int              // current position in the input
	start int              // start of this token
	width int              // width of the last rune
	items chan calcSymType // channel of scanned items
	state stateFn
}

func (l *calcLex) Lex(lval *calcSymType) int {
	for {
		select {
		case item := <-l.items:
			*lval = item
			return item.kind
		default:
			l.state = l.state(l)
		}
	}
	panic("unreachable")
}

func newCalcLex(input string) calcLexer {
	return &calcLex{
		s:     input,
		items: make(chan calcSymType, 2),
		state: lexStatement,
	}
}

func (l *calcLex) Error(s string) {
	panic(s)
	fmt.Printf("syntax error: %s\n", s)
}

func (l *calcLex) next() int {
	if l.pos >= len(l.s) {
		return 0
	}
	rune, width := utf8.DecodeRuneInString(l.s[l.pos:])
	l.pos += width
	l.width = width
	return rune
}

func (l *calcLex) backup() {
	l.pos -= l.width
}

func (l *calcLex) peek() int {
	peek := l.next()
	l.backup()
	return peek
}

func (l *calcLex) ignore() {
	l.start = l.pos
}

func (l *calcLex) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *calcLex) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {}
	l.backup()
}

func (l *calcLex) emit(t int) {
	item := calcSymType{kind: t, id:l.s[l.start:l.pos]}
	if t == NUMBER {
		item.val, _ = strconv.Atoi(item.id)
	}
	l.items <- item
	l.ignore()
}

func (l *calcLex) errorf(format string, args ...interface{}) stateFn {
	fmt.Printf(format, args...)
	l.emit(0)
	return nil
}

func lexStatement(l *calcLex) stateFn {
	switch r := l.next(); {
	case r == eof:
		l.emit(eof)
	case r == '\n':
		l.emit('\n')
	case unicode.IsSpace(r):
		l.ignore()
	case unicode.IsDigit(r):
		l.backup()
		return lexNumber
	case unicode.IsLower(r):
		l.backup()
		return lexIdentifier
	case isOperator(r):
		l.emit(r)
	default:
		return l.errorf("unrecognized char: %#U\n", r)
	}
	return lexStatement
}

func lexNumber(l *calcLex) stateFn {
	l.acceptRun("0123456789")
	l.emit(NUMBER)
	return lexStatement
}

func lexIdentifier(l *calcLex) stateFn {
	for isAlphaNumeric(l.next()) {}
	l.backup()
	fmt.Printf("ident: %s\n", l.s[l.start:l.pos])
	l.emit(IDENTIFIER)
	return lexStatement
}

func isOperator(rune int) bool {
	return bytes.IndexRune([]byte("+-*/|&="), rune) > -1
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r int) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
