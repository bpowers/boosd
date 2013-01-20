package parser

import (
	"bytes"
	"fmt"
	"go/token"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"
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

type tok struct {
	pos    token.Pos
	val    string
	kind   itemType
	yyKind int
}

type stateFn func(*boosdLex) stateFn

type boosdLex struct {
	f     *token.File
	s     string // the string being scanned
	pos   int    // current position in the input
	start int    // start of this token
	width int    // width of the last rune
	last  tok
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

func newBoosdLex(input string, file *token.File, result *File) *boosdLex {
	return &boosdLex{
		f:     file,
		s:     input,
		items: make(chan tok, 2),
		state: lexStatement,
		file:  result,
	}
}

func (l *boosdLex) getLine(pos token.Position) string {
	result := l.s[pos.Offset-pos.Column:]
	if newline := strings.IndexRune(result, '\n'); newline != -1 {
		result = result[:newline]
	}
	return result
}

func (l *boosdLex) Error(s string) {
	pos := l.f.Position(l.last.pos)
	line := l.getLine(pos)
	// we want the number of spaces (taking into account tabs)
	// before the problematic token
	prefixLen := pos.Column + strings.Count(line[:pos.Column], "\t")*7 - 1
	prefix := strings.Repeat(" ", prefixLen)

	line = strings.Replace(line, "\t", "        ", -1)

	fmt.Printf("%s:%d:%d: error: %s\n", pos.Filename,
		pos.Line, pos.Column, s)
	fmt.Printf("%s\n", line)
	fmt.Printf("%s^\n", prefix)
}

func (l *boosdLex) next() rune {
	if l.pos >= len(l.s) {
		return 0
	}
	r, width := utf8.DecodeRuneInString(l.s[l.pos:])
	l.pos += width
	l.width = width

	if r == '\n' {
		l.f.AddLine(l.pos + 1)
	}
	return r
}

func (l *boosdLex) backup() {
	l.pos -= l.width
}

func (l *boosdLex) peek() rune {
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

func (l *boosdLex) emit(yyTy rune, ty itemType) {
	t := tok{
		pos:    l.f.Pos(l.pos),
		val:    l.s[l.start:l.pos],
		yyKind: int(yyTy),
		kind:   ty,
	}
	//log.Printf("t: %#v\n", t)
	l.last = t
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
		if l.peek() == '*' {
			l.next()
			return lexMultiComment
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

func lexMultiComment(l *boosdLex) stateFn {
	// skip everything until the end of the line, or the end of
	// the file, whichever is first
	for r := l.next(); ; r = l.next() {
		if r == eof {
			l.backup()
			break
		}
		if r != '*' {
			continue
		}
		if l.peek() == '/' {
			l.next()
			break
		}
	}
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
	l.emit(YKIND_DECL, itemKindDecl)
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
	l.emit(YNUMBER, itemNumber)
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
	l.emit(YLITERAL, itemLiteral)
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
		l.emit(YKIND, itemKeyword)
	case id == "import":
		l.emit(YIMPORT, itemKeyword)
	case id == "package":
		l.emit(YPACKAGE, itemKeyword)
	case id == "model":
		l.emit(YMODEL, itemKeyword)
	case id == "interface":
		l.emit(YINTERFACE, itemKeyword)
	case id == "specializes":
		l.emit(YSPECIALIZES, itemKeyword)
	default:
		l.emit(YIDENT, itemIdentifier)
	}
	return lexStatement
}

func isLiteralStart(r rune) bool {
	return r == '"'
}

func isOperator(r rune) bool {
	return bytes.IndexRune([]byte(",+-*/|&=(){}[]:"), r) > -1
}

func isIdentifierStart(r rune) bool {
	return !(unicode.IsDigit(r) || unicode.IsSpace(r) || isOperator(r))
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return !(unicode.IsSpace(r) || isOperator(r) || r == ';')
}
