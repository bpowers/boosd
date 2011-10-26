package main

import (
	"bytes"
	"fmt"
	"log"
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
	itemSemi
	itemOperator
	itemKindDecl
	itemKeyword
)


type kind struct {
	names []string
	def   string
}

type tok struct {
	val    string
	line   int
	off    int
	yyKind int
	kind   itemType
}

type stateFn func (*calcLex) stateFn

type calcLex struct {
	s     string           // the string being scanned
	pos   int              // current position in the input
	start int              // start of this token
	width int              // width of the last rune
	items chan tok         // channel of scanned items
	state stateFn
	semi  bool
}

func (l *calcLex) Lex(lval *calcSymType) int {
	for {
		select {
		case item := <-l.items:
			lval.tok = item
			return item.yyKind
		default:
			l.state = l.state(l)
		}
	}
	panic("unreachable")
}

func newCalcLex(input string) calcLexer {
	return &calcLex{
		s:     input,
		items: make(chan tok, 2),
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

func (l *calcLex) emit(yyTy int, ty itemType) {
	t := tok{val:l.s[l.start:l.pos], yyKind: yyTy, kind: ty}
	log.Printf("t: %#v\n", t)
	l.items <- t
	l.ignore()

	switch {
	case ty == itemIdentifier || ty == itemNumber || ty == itemKindDecl:
		l.semi = true
	default:
		l.semi = false
	}
}

func (l *calcLex) errorf(format string, args ...interface{}) stateFn {
	fmt.Printf(format, args...)
	log.Printf("a")
	l.emit(eof, itemEOF)
	return nil
}

func lexStatement(l *calcLex) stateFn {
	switch r := l.next(); {
	case r == eof:
		log.Printf("b")
		l.emit(eof, itemEOF)
	case r == '/':
		if l.peek() == '/' {
			l.next()
			return lexComment
		}
		log.Printf("d")
		l.emit(r, itemOperator)
	case r == '`':
		return lexType
	case r == ';':
		l.emit(r, itemSemi)
	case unicode.IsSpace(r):
		if r == '\n' && l.semi {
			l.emit(';', itemSemi)
		}
//		log.Print("1 ignoring:", l.s[l.start:l.pos])
		l.ignore()
	case unicode.IsDigit(r):
		l.backup()
		return lexNumber
	case isIdentifierStart(r):
		l.backup()
		return lexIdentifier
	case isOperator(r):
		log.Printf("e")
		l.emit(r, itemOperator)
	default:
		return l.errorf("unrecognized char: %#U\n", r)
	}
	return lexStatement
}

func lexComment(l *calcLex) stateFn {
	// skip everything until the end of the line, or the end of
	// the file, whichever is first
	for r := l.next(); r != '\n' && r != eof; r = l.next() {}
	l.backup()
//	log.Print("2 ignoring:", l.s[l.start:l.pos])
	l.ignore()
	return lexStatement
}

func lexType(l *calcLex) stateFn {
	log.Print("3 ignoring:", l.s[l.start:l.pos])
	l.ignore()
	for r := l.next(); r != '`' && r != eof; r = l.next() {}
	l.backup()

	if l.peek() != '`' {
		return l.errorf("unexpected EOF")
	}
	l.emit(KIND_DECL, itemKindDecl)
	l.next()
	log.Print("4 ignoring:", l.s[l.start:l.pos])
	l.ignore()
	return lexStatement
}

func lexNumber(l *calcLex) stateFn {
	l.acceptRun("0123456789")
	log.Printf("g")
	l.emit(NUMBER, itemNumber)
	return lexStatement
}

func lexIdentifier(l *calcLex) stateFn {
	for isAlphaNumeric(l.next()) {}
	l.backup()
	switch id := l.s[l.start:l.pos]; {
	case id == "kind":
		l.emit(KIND, itemKeyword)
	default:
		l.emit(ID, itemIdentifier)
	}
	return lexStatement
}

func isOperator(rune int) bool {
	return bytes.IndexRune([]byte(",+-*/|&="), rune) > -1
}

func isIdentifierStart(r int) bool {
	return !(unicode.IsDigit(r) || unicode.IsSpace(r) || isOperator(r))
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r int) bool {
	return !(unicode.IsSpace(r) || isOperator(r))
}
