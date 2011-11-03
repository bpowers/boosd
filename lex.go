package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"unicode"
	"utf8"
)

const eof = 0

type itemType int

const (
	itemEOF itemType = iota
	itemIdentifier
	itemNumber
	itemSemi
	itemOperator
	itemKindDecl
	itemKeyword
	itemLiteral
	itemLBracket
	itemRBracket
	itemLParen
	itemRParen
	itemLSquare
	itemRSquare
)

type astType int

const (
	astRef astType = iota
	astNumber
	astOp
	astFn
)

// ast node
type node struct {

}

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

type mdl struct {
	name string
	sig  []string
}

type File struct {
	pkgName string
	imports []string
	kinds   []kind
	models  []mdl
}

type stateFn func(*boosdLex) stateFn

type boosdLex struct {
	s     string   // the string being scanned
	pos   int      // current position in the input
	start int      // start of this token
	width int      // width of the last rune
	items chan tok // channel of scanned items
	state stateFn
	semi  bool

	file *File
}

func (l *boosdLex) Lex(lval *boosdSymType) int {
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

func newBoosdLex(input string, file *File) *boosdLex {
	return &boosdLex{
		s:     input,
		items: make(chan tok, 2),
		state: lexStatement,
		file:  file,
	}
}

func (l *boosdLex) Error(s string) {
	fmt.Printf("syntax error: %s\n", s)
}

func (l *boosdLex) next() int {
	if l.pos >= len(l.s) {
		return 0
	}
	rune, width := utf8.DecodeRuneInString(l.s[l.pos:])
	l.pos += width
	l.width = width
	return rune
}

func (l *boosdLex) backup() {
	l.pos -= l.width
}

func (l *boosdLex) peek() int {
	peek := l.next()
	l.backup()
	return peek
}

func (l *boosdLex) ignore() {
	l.start = l.pos
}

func (l *boosdLex) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *boosdLex) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *boosdLex) emit(yyTy int, ty itemType) {
	t := tok{val: l.s[l.start:l.pos], yyKind: yyTy, kind: ty}
	//	log.Printf("t: %#v\n", t)
	l.items <- t
	l.ignore()

	switch {
	case ty == itemRBracket || ty == itemRParen || ty == itemRSquare:
		fallthrough
	case ty == itemIdentifier || ty == itemNumber || ty == itemKindDecl || ty == itemLiteral:
		l.semi = true
	default:
		l.semi = false
	}
}

func (l *boosdLex) errorf(format string, args ...interface{}) stateFn {
	log.Printf(format, args...)
	l.emit(eof, itemEOF)
	return nil
}

func lexStatement(l *boosdLex) stateFn {
	switch r := l.next(); {
	case r == eof:
		if l.semi {
			l.emit(';', itemSemi)
		}
		l.emit(eof, itemEOF)
	case r == '/':
		if l.peek() == '/' {
			l.next()
			return lexComment
		}
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
	case unicode.IsDigit(r) || r == '.':
		l.backup()
		return lexNumber
	case isLiteralStart(r):
		l.backup()
		return lexLiteral
	case isIdentifierStart(r):
		l.backup()
		return lexIdentifier
	case isOperator(r):
		l.backup()
		return lexOperator
	default:
		return l.errorf("unrecognized char: %#U\n", r)
	}
	return lexStatement
}

func lexOperator(l *boosdLex) stateFn {
	ty := itemOperator
	r := l.next()
	switch {
	case r == '{':
		ty = itemLBracket
	case r == '}':
		ty = itemRBracket
	case r == '(':
		ty = itemLParen
	case r == ')':
		ty = itemRParen
	case r == '[':
		ty = itemLSquare
	case r == ']':
		ty = itemRSquare
	}
	l.emit(r, ty)
	return lexStatement
}

func lexComment(l *boosdLex) stateFn {
	// skip everything until the end of the line, or the end of
	// the file, whichever is first
	for r := l.next(); r != '\n' && r != eof; r = l.next() {
	}
	l.backup()
	//	log.Print("2 ignoring:", l.s[l.start:l.pos])
	l.ignore()
	return lexStatement
}

func lexType(l *boosdLex) stateFn {
	l.ignore()
	for r := l.next(); r != '`' && r != eof; r = l.next() {
	}
	l.backup()

	if l.peek() != '`' {
		return l.errorf("unexpected EOF")
	}
	l.emit(KIND_DECL, itemKindDecl)
	l.next()
	l.ignore()
	return lexStatement
}

func lexNumber(l *boosdLex) stateFn {
	l.acceptRun("0123456789")
	l.accept(".")
	l.acceptRun("0123456789")
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	l.emit(NUMBER, itemNumber)
	return lexStatement
}

func lexLiteral(l *boosdLex) stateFn {
	delim := l.next()
	l.ignore()
	for r := l.next(); r != delim && r != eof; r = l.next() {
	}
	l.backup()

	if l.peek() != delim {
		return l.errorf("unexpected EOF")
	}
	l.emit(LITERAL, itemLiteral)
	l.next()
	l.ignore()
	return lexStatement
}

func lexIdentifier(l *boosdLex) stateFn {
	for isAlphaNumeric(l.next()) {
	}
	l.backup()
	switch id := l.s[l.start:l.pos]; {
	case id == "kind":
		l.emit(KIND, itemKeyword)
	case id == "import":
		l.emit(IMPORT, itemKeyword)
	case id == "package":
		l.emit(PACKAGE, itemKeyword)
	case id == "model":
		l.emit(MODEL, itemKeyword)
	case id == "interface":
		l.emit(INTERFACE, itemKeyword)
	case id == "callable":
		l.emit(CALLABLE, itemKeyword)
	case id == "specializes":
		l.emit(SPECIALIZES, itemKeyword)
	default:
		l.emit(IDENT, itemIdentifier)
	}
	return lexStatement
}

func isLiteralStart(r int) bool {
	return r == '"'
}

func isOperator(rune int) bool {
	return bytes.IndexRune([]byte(",+-*/|&=(){}[]:"), rune) > -1
}

func isIdentifierStart(r int) bool {
	return !(unicode.IsDigit(r) || unicode.IsSpace(r) || isOperator(r))
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r int) bool {
	return !(unicode.IsSpace(r) || isOperator(r))
}
