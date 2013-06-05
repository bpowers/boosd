//line parse.y:6
package boosd

import __yyfmt__ "fmt"

//line parse.y:7
import (
	"fmt"
	"go/token"
)

//line parse.y:18
type boosdSymType struct {
	yys    int
	tok    tok
	strs   []string
	ids    []*Ident
	file   File
	id     *Ident
	str    string
	lit    *BasicLit
	exprs  []Expr
	pexprs []*PairExpr
	expr   Expr
	stmt   Stmt
	tlDecl Decl
	decl   *VarDecl
	decls  []Decl
	block  *BlockStmt
}

const YIMPORT = 57346
const YKIND = 57347
const YKIND_DECL = 57348
const YPACKAGE = 57349
const YSPECIALIZES = 57350
const YINTERFACE = 57351
const YMODEL = 57352
const YIDENT = 57353
const YLITERAL = 57354
const YNUMBER = 57355
const UMINUS = 57356
const FN_CALL = 57357

var boosdToknames = []string{
	"YIMPORT",
	"YKIND",
	"YKIND_DECL",
	"YPACKAGE",
	"YSPECIALIZES",
	"YINTERFACE",
	"YMODEL",
	"YIDENT",
	"YLITERAL",
	"YNUMBER",
	" +",
	" -",
	" *",
	" /",
	" ^",
	"UMINUS",
	"FN_CALL",
}
var boosdStatenames = []string{}

const boosdEofCode = 1
const boosdErrCode = 2
const boosdMaxDepth = 200

//line parse.y:342

/* start of programs */

func Parse(f *token.File, str string) (*File, error) {
	// this is weird, but without passing in a reference to this
	// result object, there isn't another good way to keep the
	// parser and lexer reentrant.
	result := &File{}
	err := boosdParse(newBoosdLex(str, f, result))
	if err != 0 {
		return nil, fmt.Errorf("%d parse errors", err)
	}

	return result, nil
}

//line yacctab:1
var boosdExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const boosdNprod = 56
const boosdPrivate = 57344

var boosdTokenNames []string
var boosdStates []string

const boosdLast = 143

var boosdAct = []int{

	46, 51, 44, 70, 66, 20, 9, 68, 57, 55,
	58, 12, 101, 15, 60, 61, 62, 63, 64, 87,
	13, 71, 53, 23, 48, 56, 25, 86, 92, 57,
	94, 58, 29, 89, 91, 34, 47, 36, 52, 40,
	39, 38, 28, 43, 96, 45, 54, 100, 65, 67,
	41, 35, 59, 60, 61, 62, 63, 64, 77, 78,
	74, 79, 80, 81, 82, 83, 75, 84, 13, 85,
	22, 24, 16, 88, 60, 61, 62, 63, 64, 64,
	74, 90, 53, 13, 10, 53, 21, 48, 10, 13,
	93, 95, 97, 98, 13, 42, 13, 27, 99, 47,
	22, 52, 72, 62, 63, 64, 22, 31, 60, 61,
	62, 63, 64, 60, 61, 62, 63, 64, 22, 19,
	18, 8, 5, 13, 7, 3, 6, 11, 33, 69,
	76, 50, 37, 73, 49, 32, 30, 17, 26, 4,
	1, 14, 2,
}
var boosdPact = []int{

	-1000, -1000, 118, 116, -1000, 76, 85, -1000, 85, 51,
	-1000, -1000, 110, -1000, 64, -1000, -1000, 100, -1000, -1000,
	50, 85, -1000, 89, -1000, -1000, 19, 85, -1000, -1000,
	83, 30, -1000, 16, 112, -1000, -1000, 29, 72, -1000,
	100, -1000, -1000, 2, -1000, -1000, 94, 9, 9, -22,
	-1000, -1000, -6, -1000, -1000, 78, -1000, 9, 9, -1000,
	9, 9, 9, 9, 9, 39, -19, -1000, 9, -3,
	-1000, 69, -1000, -1000, 7, 57, 6, 99, 60, 87,
	87, 61, 61, -1000, -1000, 0, -1000, -6, 22, 9,
	-1000, -1000, 9, -1000, -1000, -1000, 69, 26, 99, -16,
	-1000, -1000,
}
var boosdPgo = []int{

	0, 142, 141, 140, 139, 4, 138, 137, 136, 135,
	0, 1, 3, 134, 2, 5, 133, 132, 131, 130,
	9, 129, 128, 127, 126, 6, 125, 124,
}
var boosdR1 = []int{

	0, 3, 1, 1, 4, 26, 26, 27, 15, 15,
	2, 2, 24, 24, 23, 7, 7, 6, 6, 8,
	8, 9, 9, 22, 22, 17, 17, 17, 17, 20,
	20, 16, 14, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 18, 5, 25, 11,
	19, 19, 13, 21, 21, 12,
}
var boosdR2 = []int{

	0, 3, 0, 2, 3, 0, 2, 4, 0, 1,
	1, 3, 0, 2, 8, 1, 1, 0, 2, 0,
	2, 2, 3, 2, 3, 4, 5, 2, 2, 0,
	2, 4, 2, 3, 3, 3, 3, 3, 3, 2,
	4, 4, 4, 1, 1, 1, 1, 1, 1, 1,
	1, 3, 3, 1, 3, 5,
}
var boosdChk = []int{

	-1000, -3, -1, -26, -4, 4, -24, -27, 5, -25,
	12, -23, -5, 11, -2, -5, 21, -7, 10, 9,
	-15, 22, 6, -15, 21, -5, -6, 8, 23, -5,
	-8, 24, -9, -22, -5, 21, 21, -17, 25, -15,
	-5, 21, 23, -5, -14, -25, -10, 27, 15, -13,
	-18, -11, 29, 13, -15, -20, 23, 27, 29, -15,
	14, 15, 16, 17, 18, -10, -5, -10, 29, -21,
	-12, 27, 24, -16, -5, -20, -19, -10, -10, -10,
	-10, -10, -10, -10, 28, -10, 30, 22, -11, 26,
	24, 28, 22, 30, 30, -12, 22, -14, -10, -11,
	21, 28,
}
var boosdDef = []int{

	2, -2, 5, 12, 3, 0, 1, 6, 0, 0,
	48, 13, 0, 47, 8, 10, 4, 8, 15, 16,
	0, 0, 9, 17, 7, 11, 0, 0, 19, 18,
	0, 0, 20, 0, 8, 14, 21, 0, 0, 23,
	8, 22, 29, 46, 27, 28, 8, 0, 0, 43,
	44, 45, 0, 49, 24, 0, 29, 0, 0, 32,
	0, 0, 0, 0, 0, 0, 46, 39, 0, 0,
	53, 0, 25, 30, 0, 0, 0, 50, 0, 34,
	35, 36, 37, 38, 33, 0, 52, 0, 0, 0,
	26, 40, 0, 42, 41, 54, 0, 0, 51, 0,
	31, 55,
}
var boosdTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	27, 28, 16, 14, 22, 15, 3, 17, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 26, 21,
	3, 25, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 29, 3, 30, 18, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 23, 3, 24,
}
var boosdTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 19, 20,
}
var boosdTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var boosdDebug = 0

type boosdLexer interface {
	Lex(lval *boosdSymType) int
	Error(s string)
}

const boosdFlag = -1000

func boosdTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(boosdToknames) {
		if boosdToknames[c-4] != "" {
			return boosdToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func boosdStatname(s int) string {
	if s >= 0 && s < len(boosdStatenames) {
		if boosdStatenames[s] != "" {
			return boosdStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func boosdlex1(lex boosdLexer, lval *boosdSymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = boosdTok1[0]
		goto out
	}
	if char < len(boosdTok1) {
		c = boosdTok1[char]
		goto out
	}
	if char >= boosdPrivate {
		if char < boosdPrivate+len(boosdTok2) {
			c = boosdTok2[char-boosdPrivate]
			goto out
		}
	}
	for i := 0; i < len(boosdTok3); i += 2 {
		c = boosdTok3[i+0]
		if c == char {
			c = boosdTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = boosdTok2[1] /* unknown char */
	}
	if boosdDebug >= 3 {
		__yyfmt__.Printf("lex %U %s\n", uint(char), boosdTokname(c))
	}
	return c
}

func boosdParse(boosdlex boosdLexer) int {
	var boosdn int
	var boosdlval boosdSymType
	var boosdVAL boosdSymType
	boosdS := make([]boosdSymType, boosdMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	boosdstate := 0
	boosdchar := -1
	boosdp := -1
	goto boosdstack

ret0:
	return 0

ret1:
	return 1

boosdstack:
	/* put a state and value onto the stack */
	if boosdDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", boosdTokname(boosdchar), boosdStatname(boosdstate))
	}

	boosdp++
	if boosdp >= len(boosdS) {
		nyys := make([]boosdSymType, len(boosdS)*2)
		copy(nyys, boosdS)
		boosdS = nyys
	}
	boosdS[boosdp] = boosdVAL
	boosdS[boosdp].yys = boosdstate

boosdnewstate:
	boosdn = boosdPact[boosdstate]
	if boosdn <= boosdFlag {
		goto boosddefault /* simple state */
	}
	if boosdchar < 0 {
		boosdchar = boosdlex1(boosdlex, &boosdlval)
	}
	boosdn += boosdchar
	if boosdn < 0 || boosdn >= boosdLast {
		goto boosddefault
	}
	boosdn = boosdAct[boosdn]
	if boosdChk[boosdn] == boosdchar { /* valid shift */
		boosdchar = -1
		boosdVAL = boosdlval
		boosdstate = boosdn
		if Errflag > 0 {
			Errflag--
		}
		goto boosdstack
	}

boosddefault:
	/* default state action */
	boosdn = boosdDef[boosdstate]
	if boosdn == -2 {
		if boosdchar < 0 {
			boosdchar = boosdlex1(boosdlex, &boosdlval)
		}

		/* look through exception table */
		xi := 0
		for {
			if boosdExca[xi+0] == -1 && boosdExca[xi+1] == boosdstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			boosdn = boosdExca[xi+0]
			if boosdn < 0 || boosdn == boosdchar {
				break
			}
		}
		boosdn = boosdExca[xi+1]
		if boosdn < 0 {
			goto ret0
		}
	}
	if boosdn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			boosdlex.Error("syntax error")
			Nerrs++
			if boosdDebug >= 1 {
				__yyfmt__.Printf("%s", boosdStatname(boosdstate))
				__yyfmt__.Printf("saw %s\n", boosdTokname(boosdchar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for boosdp >= 0 {
				boosdn = boosdPact[boosdS[boosdp].yys] + boosdErrCode
				if boosdn >= 0 && boosdn < boosdLast {
					boosdstate = boosdAct[boosdn] /* simulate a shift of "error" */
					if boosdChk[boosdstate] == boosdErrCode {
						goto boosdstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if boosdDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", boosdS[boosdp].yys)
				}
				boosdp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if boosdDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", boosdTokname(boosdchar))
			}
			if boosdchar == boosdEofCode {
				goto ret1
			}
			boosdchar = -1
			goto boosdnewstate /* try again in the same state */
		}
	}

	/* reduction by production boosdn */
	if boosdDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", boosdn, boosdStatname(boosdstate))
	}

	boosdnt := boosdn
	boosdpt := boosdp
	_ = boosdpt // guard against "declared and not used"

	boosdp -= boosdR2[boosdn]
	boosdVAL = boosdS[boosdp+1]

	/* consult goto table to find next state */
	boosdn = boosdR1[boosdn]
	boosdg := boosdPgo[boosdn]
	boosdj := boosdg + boosdS[boosdp].yys + 1

	if boosdj >= boosdLast {
		boosdstate = boosdAct[boosdg]
	} else {
		boosdstate = boosdAct[boosdj]
		if boosdChk[boosdstate] != -boosdn {
			boosdstate = boosdAct[boosdg]
		}
	}
	// dummy call; replaced with literal code
	switch boosdnt {

	case 1:
		//line parse.y:70
		{
			boosdVAL.file.Decls = boosdS[boosdpt-0].decls
			*boosdlex.(*boosdLex).file = boosdVAL.file
		}
	case 2:
		//line parse.y:76
		{
		}
	case 3:
		//line parse.y:78
		{
			boosdVAL.strs = append(boosdS[boosdpt-1].strs, boosdS[boosdpt-0].str)
		}
	case 4:
		//line parse.y:84
		{
		}
	case 5:
		//line parse.y:88
		{
		}
	case 6:
		//line parse.y:90
		{
		}
	case 7:
		//line parse.y:95
		{
		}
	case 8:
		//line parse.y:99
		{
			boosdVAL.expr = nil
		}
	case 9:
		//line parse.y:103
		{
			boosdVAL.expr = &BasicLit{Kind: token.STRING, Value: boosdS[boosdpt-0].tok.val}
		}
	case 10:
		//line parse.y:109
		{
			boosdVAL.ids = []*Ident{boosdS[boosdpt-0].id}
		}
	case 11:
		//line parse.y:113
		{
			boosdVAL.ids = append(boosdS[boosdpt-2].ids, boosdS[boosdpt-0].id)
		}
	case 12:
		//line parse.y:118
		{
		}
	case 13:
		//line parse.y:120
		{
			boosdVAL.decls = append(boosdS[boosdpt-1].decls, boosdS[boosdpt-0].tlDecl)
		}
	case 14:
		//line parse.y:126
		{
			if boosdS[boosdpt-6].tok.val == "model" {
				boosdVAL.tlDecl = &ModelDecl{Name: boosdS[boosdpt-7].id, Body: boosdS[boosdpt-2].block}
			} else {
				boosdVAL.tlDecl = &InterfaceDecl{Name: boosdS[boosdpt-7].id, Body: boosdS[boosdpt-2].block}
			}
		}
	case 15:
		//line parse.y:136
		{
			boosdVAL.tok = boosdS[boosdpt-0].tok
		}
	case 16:
		//line parse.y:140
		{
			boosdVAL.tok = boosdS[boosdpt-0].tok
		}
	case 17:
		//line parse.y:145
		{
		}
	case 18:
		//line parse.y:147
		{
			boosdVAL.id = boosdS[boosdpt-0].id
		}
	case 19:
		//line parse.y:153
		{
			boosdVAL.block = &BlockStmt{List: []Stmt{}}
		}
	case 20:
		//line parse.y:157
		{
			boosdVAL.block = boosdS[boosdpt-1].block
			boosdVAL.block.List = append(boosdS[boosdpt-1].block.List, boosdS[boosdpt-0].stmt)
		}
	case 21:
		//line parse.y:164
		{
			boosdVAL.stmt = &DeclStmt{boosdS[boosdpt-1].decl}
		}
	case 22:
		//line parse.y:168
		{
			boosdVAL.stmt = &AssignStmt{Lhs: boosdS[boosdpt-2].decl, Rhs: boosdS[boosdpt-1].expr}
		}
	case 23:
		//line parse.y:175
		{
			boosdVAL.decl = &VarDecl{Name: boosdS[boosdpt-1].id, Units: boosdS[boosdpt-0].expr}
		}
	case 24:
		//line parse.y:179
		{
			boosdVAL.decl = &VarDecl{Name: boosdS[boosdpt-2].id, Type: boosdS[boosdpt-1].id, Units: boosdS[boosdpt-0].expr}
		}
	case 25:
		//line parse.y:185
		{
			boosdVAL.expr = &CompositeLit{Type: NewIdent("stock"), Elts: boosdS[boosdpt-1].exprs}
		}
	case 26:
		//line parse.y:189
		{
			boosdVAL.expr = &CompositeLit{Type: boosdS[boosdpt-3].id, Elts: boosdS[boosdpt-1].exprs}
		}
	case 27:
		//line parse.y:193
		{
			boosdVAL.expr = boosdS[boosdpt-0].expr
		}
	case 28:
		//line parse.y:197
		{
			boosdVAL.expr = boosdS[boosdpt-0].lit
		}
	case 29:
		//line parse.y:202
		{
			boosdVAL.exprs = []Expr{}
		}
	case 30:
		//line parse.y:206
		{
			boosdVAL.exprs = append(boosdS[boosdpt-1].exprs, boosdS[boosdpt-0].expr)
		}
	case 31:
		//line parse.y:212
		{
			boosdVAL.expr = &KeyValueExpr{Key: boosdS[boosdpt-3].id, Value: boosdS[boosdpt-1].expr}
		}
	case 32:
		//line parse.y:218
		{
			boosdVAL.expr = &UnitExpr{boosdS[boosdpt-1].expr, boosdS[boosdpt-0].expr}
		}
	case 33:
		//line parse.y:224
		{
			boosdVAL.expr = boosdS[boosdpt-1].expr
		}
	case 34:
		//line parse.y:228
		{
			boosdVAL.expr = &BinaryExpr{X: boosdS[boosdpt-2].expr, Y: boosdS[boosdpt-0].expr, Op: token.ADD}
		}
	case 35:
		//line parse.y:232
		{
			boosdVAL.expr = &BinaryExpr{X: boosdS[boosdpt-2].expr, Y: boosdS[boosdpt-0].expr, Op: token.SUB}
		}
	case 36:
		//line parse.y:236
		{
			boosdVAL.expr = &BinaryExpr{X: boosdS[boosdpt-2].expr, Y: boosdS[boosdpt-0].expr, Op: token.MUL}
		}
	case 37:
		//line parse.y:240
		{
			boosdVAL.expr = &BinaryExpr{X: boosdS[boosdpt-2].expr, Y: boosdS[boosdpt-0].expr, Op: token.QUO}
		}
	case 38:
		//line parse.y:244
		{
			boosdVAL.expr = &BinaryExpr{X: boosdS[boosdpt-2].expr, Y: boosdS[boosdpt-0].expr, Op: token.XOR}
		}
	case 39:
		//line parse.y:248
		{
			boosdVAL.expr = &UnaryExpr{X: boosdS[boosdpt-0].expr, Op: token.SUB}
		}
	case 40:
		//line parse.y:252
		{
			boosdVAL.expr = &CallExpr{Fun: boosdS[boosdpt-3].id, Args: boosdS[boosdpt-1].exprs}
		}
	case 41:
		//line parse.y:256
		{
			boosdVAL.expr = &IndexExpr{X: boosdS[boosdpt-3].expr, Index: boosdS[boosdpt-1].expr}
		}
	case 42:
		//line parse.y:260
		{
			boosdVAL.expr = &IndexExpr{X: boosdS[boosdpt-3].id, Index: boosdS[boosdpt-1].expr}
		}
	case 43:
		//line parse.y:264
		{
			boosdVAL.expr = boosdS[boosdpt-0].expr
		}
	case 44:
		//line parse.y:268
		{
			boosdVAL.expr = boosdS[boosdpt-0].expr
		}
	case 45:
		//line parse.y:272
		{
			boosdVAL.expr = boosdS[boosdpt-0].expr
		}
	case 46:
		//line parse.y:278
		{
			boosdVAL.expr = &RefExpr{*boosdS[boosdpt-0].id}
		}
	case 47:
		//line parse.y:283
		{
			boosdVAL.id = &Ident{Name: boosdS[boosdpt-0].tok.val}
		}
	case 48:
		//line parse.y:289
		{
			boosdVAL.lit = &BasicLit{Kind: token.STRING, Value: boosdS[boosdpt-0].tok.val}
		}
	case 49:
		//line parse.y:295
		{
			boosdVAL.expr = &BasicLit{Kind: token.FLOAT, Value: boosdS[boosdpt-0].tok.val}
		}
	case 50:
		//line parse.y:301
		{
			boosdVAL.exprs = make([]Expr, 1, 16)
			boosdVAL.exprs[0] = boosdS[boosdpt-0].expr
		}
	case 51:
		//line parse.y:306
		{
			boosdVAL.exprs = append(boosdS[boosdpt-2].exprs, boosdS[boosdpt-0].expr)
		}
	case 52:
		//line parse.y:312
		{
			boosdVAL.expr = &TableExpr{Pairs: boosdS[boosdpt-1].pexprs}
		}
	case 53:
		//line parse.y:318
		{
			boosdVAL.pexprs = make([]*PairExpr, 1, 8)
			pe, ok := boosdS[boosdpt-0].expr.(*PairExpr)
			if !ok {
				panic(fmt.Sprintf("not PairExpr 1: %#v", boosdS[boosdpt-0].expr))
			}
			boosdVAL.pexprs[0] = pe
		}
	case 54:
		//line parse.y:327
		{
			pe, ok := boosdS[boosdpt-0].expr.(*PairExpr)
			if !ok {
				panic(fmt.Sprintf("not PairExpr 1: %#v", boosdS[boosdpt-0].expr))
			}
			boosdVAL.pexprs = append(boosdS[boosdpt-2].pexprs, pe)
		}
	case 55:
		//line parse.y:337
		{
			boosdVAL.expr = &PairExpr{boosdS[boosdpt-3].expr, boosdS[boosdpt-1].expr}
		}
	}
	goto boosdstack /* stack new state and value */
}
