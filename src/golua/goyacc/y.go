//line lua.y:1
package goyacc

import __yyfmt__ "fmt"

//line lua.y:3
const AND = 57346
const BREAK = 57347
const DO = 57348
const ELSEIF = 57349
const ELSE = 57350
const END = 57351
const FALSE = 57352
const FOR = 57353
const FUNCTION = 57354
const GOTO = 57355
const IF = 57356
const IN = 57357
const LOCAL = 57358
const NIL = 57359
const NOT = 57360
const OR = 57361
const RETURN = 57362
const REPEAT = 57363
const THEN = 57364
const TRUE = 57365
const UNTIL = 57366
const WHILE = 57367
const CLOSURE = 57368
const NUMBER = 57369
const STRING = 57370
const SLTEQ = 57371
const SGTEQ = 57372
const SEQ = 57373
const SNOTEQ = 57374
const NAME = 57375
const MORE = 57376
const STRADD = 57377

var yyToknames = []string{
	"AND",
	"BREAK",
	"DO",
	"ELSEIF",
	"ELSE",
	"END",
	"FALSE",
	"FOR",
	"FUNCTION",
	"GOTO",
	"IF",
	"IN",
	"LOCAL",
	"NIL",
	"NOT",
	"OR",
	"RETURN",
	"REPEAT",
	"THEN",
	"TRUE",
	"UNTIL",
	"WHILE",
	"CLOSURE",
	"NUMBER",
	"STRING",
	"SLTEQ",
	"SGTEQ",
	"SEQ",
	"SNOTEQ",
	"NAME",
	"MORE",
	"STRADD",
	" <",
	" >",
	" +",
	" -",
	" *",
	" /",
	" %",
	" ~",
	" ^",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line lua.y:231

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 13,
	46, 70,
	49, 70,
	50, 70,
	52, 70,
	-2, 16,
	-1, 17,
	45, 31,
	48, 31,
	-2, 69,
	-1, 110,
	45, 32,
	48, 32,
	-2, 69,
}

const yyNprod = 95
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 454

var yyAct = []int{

	68, 95, 67, 90, 64, 49, 96, 56, 25, 122,
	96, 50, 140, 122, 143, 116, 41, 13, 164, 138,
	13, 107, 108, 13, 61, 97, 13, 66, 93, 97,
	65, 63, 163, 62, 159, 122, 66, 87, 88, 89,
	33, 18, 154, 112, 18, 92, 99, 18, 112, 111,
	18, 54, 112, 59, 55, 106, 169, 109, 2, 141,
	115, 13, 79, 117, 152, 24, 150, 114, 46, 121,
	152, 168, 124, 125, 126, 127, 128, 129, 130, 131,
	132, 133, 40, 17, 149, 18, 17, 13, 57, 17,
	58, 137, 17, 77, 78, 80, 18, 79, 142, 119,
	136, 118, 145, 100, 58, 50, 113, 19, 13, 58,
	147, 18, 60, 53, 151, 101, 162, 153, 155, 156,
	13, 161, 13, 160, 158, 144, 102, 17, 69, 123,
	74, 44, 18, 76, 75, 77, 78, 80, 110, 79,
	139, 165, 166, 94, 18, 4, 18, 104, 103, 21,
	134, 135, 91, 17, 73, 35, 34, 32, 52, 48,
	12, 51, 146, 47, 148, 72, 8, 5, 45, 20,
	170, 3, 1, 0, 17, 0, 0, 0, 0, 0,
	71, 0, 0, 0, 0, 0, 17, 0, 17, 0,
	83, 84, 85, 86, 72, 0, 74, 81, 82, 76,
	75, 77, 78, 80, 0, 79, 0, 0, 0, 71,
	0, 0, 0, 0, 167, 0, 0, 0, 0, 83,
	84, 85, 86, 0, 0, 74, 81, 82, 76, 75,
	77, 78, 80, 28, 79, 39, 0, 0, 22, 6,
	26, 36, 0, 157, 16, 11, 27, 10, 0, 14,
	29, 30, 0, 23, 9, 0, 19, 31, 7, 15,
	0, 0, 38, 0, 0, 28, 19, 39, 0, 0,
	120, 0, 26, 36, 37, 43, 0, 42, 27, 0,
	0, 0, 29, 30, 0, 0, 0, 0, 19, 31,
	0, 0, 0, 0, 38, 0, 0, 28, 0, 39,
	0, 0, 0, 0, 26, 36, 37, 43, 98, 42,
	27, 0, 0, 0, 29, 30, 0, 72, 0, 0,
	19, 31, 0, 0, 0, 0, 38, 0, 0, 0,
	0, 0, 71, 0, 0, 105, 0, 0, 37, 43,
	0, 42, 83, 84, 85, 86, 0, 0, 74, 81,
	82, 76, 75, 77, 78, 80, 72, 79, 70, 6,
	0, 0, 0, 0, 16, 11, 0, 10, 0, 14,
	0, 71, 0, 0, 9, 0, 0, 0, 7, 15,
	72, 83, 84, 85, 86, 0, 19, 74, 81, 82,
	76, 75, 77, 78, 80, 71, 79, 72, 0, 0,
	0, 0, 0, 0, 0, 83, 84, 85, 86, 0,
	0, 74, 81, 82, 76, 75, 77, 78, 80, 0,
	79, 0, 83, 84, 85, 86, 0, 0, 74, 81,
	82, 76, 75, 77, 78, 80, 0, 79, 83, 84,
	85, 86, 0, 0, 74, 81, 82, 76, 75, 77,
	78, 80, 0, 79,
}
var yyPact = []int{

	353, -1000, -1000, 233, -1000, -1000, 353, 287, 125, 353,
	287, 80, 6, -1000, 76, 7, 79, -1000, -19, -1000,
	-1000, -1000, -1000, 287, 119, 352, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -19, -1000, -1000, 287, 287, 287, -1,
	-1000, -1000, -27, 255, 353, -1000, 91, 117, 140, -1000,
	313, -1, -28, -1000, 287, 74, 4, 73, -1000, 71,
	15, 0, 287, 68, -1000, 66, 223, -35, 376, -1000,
	353, 287, 287, 287, 287, 287, 287, 287, 287, 287,
	287, -1000, -1000, -1000, -1000, -1000, -1000, 18, 376, 18,
	-1000, 353, 57, -1000, -36, -1000, 14, 287, -1000, -39,
	116, 287, -1000, 353, 287, 353, -1000, 51, 33, -35,
	-1000, 287, 31, -1, -5, 287, 287, 190, -1000, -10,
	-1000, -13, 287, 112, 393, 409, 376, 95, 53, 53,
	18, 18, 18, 18, 107, -15, -30, -1000, -1000, -23,
	-1000, 287, 161, -1000, -1000, 376, -1000, -1000, -1000, -1000,
	-1000, -35, -1000, -1000, -1000, -35, -35, -1000, -1000, -1000,
	376, -1000, -1000, -1000, 37, -1000, 376, 11, -1000, 287,
	376,
}
var yyPgo = []int{

	0, 172, 58, 171, 169, 168, 0, 145, 167, 166,
	163, 161, 3, 160, 2, 16, 7, 159, 5, 82,
	158, 157, 40, 156, 155, 154, 4, 152, 151, 143,
	1, 140,
}
var yyR1 = []int{

	0, 1, 2, 2, 5, 3, 3, 3, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 9, 9, 10,
	10, 17, 17, 18, 4, 4, 4, 8, 8, 8,
	8, 13, 13, 11, 11, 20, 20, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	25, 25, 25, 25, 25, 25, 19, 19, 19, 22,
	22, 15, 15, 26, 26, 21, 12, 27, 28, 28,
	28, 28, 16, 16, 24, 24, 14, 14, 23, 23,
	29, 29, 31, 30, 30,
}
var yyR2 = []int{

	0, 1, 1, 2, 3, 1, 2, 0, 1, 3,
	5, 4, 2, 3, 3, 3, 1, 4, 4, 1,
	3, 1, 3, 3, 1, 1, 2, 2, 4, 4,
	4, 1, 3, 1, 3, 1, 3, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 2, 2, 2,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	1, 1, 1, 1, 1, 1, 1, 4, 3, 1,
	1, 2, 4, 2, 3, 2, 3, 3, 1, 1,
	3, 0, 1, 3, 2, 3, 1, 3, 2, 3,
	1, 3, 1, 3, 5,
}
var yyChk = []int{

	-1000, -1, -2, -3, -7, -8, 6, 25, -9, 21,
	14, 12, -13, -15, 16, 26, 11, -19, -22, 33,
	-4, -7, 5, 20, -2, -6, 17, 23, 10, 27,
	28, 34, -21, -22, -23, -24, 18, 51, 39, 12,
	-19, -15, 54, 52, 6, -5, -2, -10, -17, -18,
	-6, -11, -20, 33, 45, 48, -16, 12, 33, 46,
	33, -16, 52, 50, -26, 49, 46, -14, -6, 9,
	6, 19, 4, -25, 35, 39, 38, 40, 41, 44,
	42, 36, 37, 29, 30, 31, 32, -6, -6, -6,
	-12, -27, 46, 55, -29, -30, 33, 52, 53, -14,
	-2, 24, 9, 8, 7, 22, -12, 49, 50, -14,
	-19, 45, 48, 33, -16, 45, 15, -6, 33, 33,
	47, -14, 48, -2, -6, -6, -6, -6, -6, -6,
	-6, -6, -6, -6, -2, -28, -16, 34, 55, -31,
	48, 45, -6, 53, 9, -6, -2, -18, -2, 33,
	33, -14, 33, -12, 47, -14, -14, 53, -26, 47,
	-6, 9, 9, 47, 48, -30, -6, 53, 34, 45,
	-6,
}
var yyDef = []int{

	7, -2, 1, 2, 5, 8, 7, 0, 0, 7,
	0, 0, 0, -2, 0, 0, 0, -2, 0, 66,
	3, 6, 24, 25, 0, 0, 37, 38, 39, 40,
	41, 42, 43, 44, 45, 46, 0, 0, 0, 0,
	69, 70, 0, 0, 7, 12, 0, 0, 19, 21,
	0, 0, 33, 35, 0, 0, 27, 0, 82, 0,
	82, 0, 0, 0, 71, 0, 0, 26, 86, 9,
	7, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 60, 61, 62, 63, 64, 65, 47, 48, 49,
	75, 7, 81, 88, 0, 90, 0, 0, 84, 0,
	0, 0, 13, 7, 0, 7, 14, 0, 0, 15,
	-2, 0, 0, 0, 0, 0, 0, 0, 68, 0,
	73, 0, 0, 0, 50, 51, 52, 53, 54, 55,
	56, 57, 58, 59, 0, 0, 78, 79, 89, 0,
	92, 0, 0, 85, 11, 4, 20, 22, 23, 34,
	36, 28, 83, 29, 30, 17, 18, 67, 72, 74,
	87, 10, 76, 77, 0, 91, 93, 0, 80, 0,
	94,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 51, 3, 42, 3, 3,
	46, 47, 40, 38, 48, 39, 50, 41, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 49, 3,
	36, 45, 37, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 52, 3, 53, 44, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 54, 3, 55, 43,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(yyToknames) {
		if yyToknames[c-4] != "" {
			return yyToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(c), uint(char))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line lua.y:51
		{
			endChunk(yylex, &yyVAL)
		}
	case 3:
		//line lua.y:55
		{
			opAppend(yylex, &yyVAL, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 4:
		//line lua.y:58
		{
			op2(yylex, &yyVAL, OP_UNTIL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 6:
		//line lua.y:62
		{
			opAppend(yylex, &yyVAL, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 9:
		//line lua.y:67
		{
			yyVAL = yyS[yypt-1]
		}
	case 10:
		//line lua.y:68
		{
			op2(yylex, &yyVAL, OP_WHILE, &yyS[yypt-3], &yyS[yypt-1])
		}
	case 11:
		//line lua.y:69
		{
			opForBind(yylex, &yyVAL, &yyS[yypt-3], &yyS[yypt-1])
		}
	case 12:
		//line lua.y:70
		{
			yyVAL = yyS[yypt-0]
		}
	case 13:
		//line lua.y:71
		{
			yyVAL = yyS[yypt-1]
		}
	case 14:
		//line lua.y:72
		{
			bindFuncName(yylex, &yyS[yypt-0], &yyS[yypt-1], "")
			op2(yylex, &yyVAL, OP_ASSIGN, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 15:
		//line lua.y:76
		{
			op2(yylex, &yyVAL, OP_ASSIGN, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 17:
		//line lua.y:80
		{
			opFor(yylex, &yyVAL, OP_FOR, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 18:
		//line lua.y:81
		{
			opFor(yylex, &yyVAL, OP_FORIN, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 20:
		//line lua.y:85
		{
			opIf(yylex, &yyVAL, nil, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 22:
		//line lua.y:89
		{
			opIf(yylex, &yyVAL, nil, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 23:
		//line lua.y:92
		{
			opIf(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0], nil)
		}
	case 25:
		//line lua.y:96
		{
			op1(yylex, &yyVAL, OP_RETURN, nil)
		}
	case 26:
		//line lua.y:97
		{
			op1(yylex, &yyVAL, OP_RETURN, &yyS[yypt-0])
		}
	case 27:
		//line lua.y:100
		{
			opLocal(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 28:
		//line lua.y:101
		{
			opLocal(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 29:
		//line lua.y:102
		{
			bindFuncName(yylex, &yyS[yypt-0], nil, yyS[yypt-1].token.image)
			var tmp yySymType
			nameAppend(yylex, &tmp, &yyS[yypt-1], nil)
			opLocal(yylex, &yyVAL, &tmp, &yyS[yypt-0])
		}
	case 30:
		//line lua.y:108
		{
			opClosure(yylex, &yyVAL, &yyS[yypt-1])
		}
	case 32:
		//line lua.y:112
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 34:
		//line lua.y:116
		{
			opValueExt(&yyS[yypt-0], yyS[yypt-0].token.image)
			op2(yylex, &yyVAL, OP_SELFM, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 35:
		//line lua.y:122
		{
			opVar(&yyVAL, &yyS[yypt-0])
		}
	case 36:
		//line lua.y:123
		{
			opValueExt(&yyS[yypt-0], yyS[yypt-0].token.image)
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 37:
		//line lua.y:129
		{
			opValue(yylex, &yyVAL)
		}
	case 38:
		//line lua.y:130
		{
			opValue(yylex, &yyVAL)
		}
	case 39:
		//line lua.y:131
		{
			opValue(yylex, &yyVAL)
		}
	case 40:
		//line lua.y:132
		{
			opValue(yylex, &yyVAL)
		}
	case 41:
		//line lua.y:133
		{
			opValue(yylex, &yyVAL)
		}
	case 47:
		//line lua.y:139
		{
			op1(yylex, &yyVAL, OP_NOT, &yyS[yypt-0])
		}
	case 48:
		//line lua.y:140
		{
			op1(yylex, &yyVAL, OP_LEN, &yyS[yypt-0])
		}
	case 49:
		//line lua.y:141
		{
			op1(yylex, &yyVAL, OP_NSIGN, &yyS[yypt-0])
		}
	case 50:
		//line lua.y:142
		{
			op2(yylex, &yyVAL, OP_OR, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 51:
		//line lua.y:143
		{
			op2(yylex, &yyVAL, OP_AND, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 52:
		//line lua.y:144
		{
			op2(yylex, &yyVAL, yyS[yypt-1].op, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 53:
		//line lua.y:145
		{
			op2(yylex, &yyVAL, OP_STRADD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 54:
		//line lua.y:146
		{
			op2(yylex, &yyVAL, OP_SUB, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 55:
		//line lua.y:147
		{
			op2(yylex, &yyVAL, OP_ADD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 56:
		//line lua.y:148
		{
			op2(yylex, &yyVAL, OP_MUL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 57:
		//line lua.y:149
		{
			op2(yylex, &yyVAL, OP_DIV, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 58:
		//line lua.y:150
		{
			op2(yylex, &yyVAL, OP_PMUL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 59:
		//line lua.y:151
		{
			op2(yylex, &yyVAL, OP_MOD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 60:
		//line lua.y:154
		{
			opFlag(&yyVAL, OP_LT)
		}
	case 61:
		//line lua.y:155
		{
			opFlag(&yyVAL, OP_GT)
		}
	case 62:
		//line lua.y:156
		{
			opFlag(&yyVAL, OP_LTEQ)
		}
	case 63:
		//line lua.y:157
		{
			opFlag(&yyVAL, OP_GTEQ)
		}
	case 64:
		//line lua.y:158
		{
			opFlag(&yyVAL, OP_EQ)
		}
	case 65:
		//line lua.y:159
		{
			opFlag(&yyVAL, OP_NOTEQ)
		}
	case 66:
		//line lua.y:162
		{
			opVar(&yyVAL, &yyS[yypt-0])
		}
	case 67:
		//line lua.y:163
		{
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-3], &yyS[yypt-1])
		}
	case 68:
		//line lua.y:164
		{
			opValueExt(&yyS[yypt-0], yyS[yypt-0].token.image)
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 71:
		//line lua.y:174
		{
			op2(yylex, &yyVAL, OP_CALL, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 72:
		//line lua.y:175
		{
			var tmp yySymType
			op2(yylex, &tmp, OP_SELFM, &yyS[yypt-3], &yyS[yypt-1])
			op2(yylex, &yyVAL, OP_CALL, &tmp, &yyS[yypt-0])
		}
	case 74:
		//line lua.y:183
		{
			yyVAL = yyS[yypt-1]
		}
	case 75:
		//line lua.y:186
		{
			yyVAL = yyS[yypt-0]
		}
	case 76:
		//line lua.y:189
		{
			opFunc(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-1])
		}
	case 77:
		//line lua.y:192
		{
			yyVAL = yyS[yypt-1]
		}
	case 79:
		//line lua.y:196
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 80:
		//line lua.y:197
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 82:
		//line lua.y:202
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 83:
		//line lua.y:203
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 84:
		//line lua.y:206
		{
			op1(yylex, &yyVAL, OP_ARRAY, nil)
		}
	case 85:
		//line lua.y:207
		{
			op1(yylex, &yyVAL, OP_ARRAY, &yyS[yypt-1])
		}
	case 86:
		//line lua.y:210
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 87:
		//line lua.y:211
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 88:
		//line lua.y:214
		{
			op1(yylex, &yyVAL, OP_TABLE, nil)
		}
	case 89:
		//line lua.y:215
		{
			op1(yylex, &yyVAL, OP_TABLE, &yyS[yypt-1])
		}
	case 91:
		//line lua.y:219
		{
			opAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 93:
		//line lua.y:225
		{
			opValueExt(&yyS[yypt-2], yyS[yypt-2].token.image)
			op2(yylex, &yyVAL, OP_FIELD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 94:
		//line lua.y:229
		{
			op2(yylex, &yyVAL, OP_FIELD, &yyS[yypt-3], &yyS[yypt-0])
		}
	}
	goto yystack /* stack new state and value */
}
