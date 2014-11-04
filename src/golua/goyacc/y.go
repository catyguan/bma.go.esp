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
const CONTINUE = 57369
const NUMBER = 57370
const STRING = 57371
const SLTEQ = 57372
const SGTEQ = 57373
const SEQ = 57374
const SNOTEQ = 57375
const NAME = 57376
const MORE = 57377
const STRADD = 57378

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
	"CONTINUE",
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

//line lua.y:225

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 13,
	47, 70,
	50, 70,
	52, 70,
	-2, 16,
	-1, 17,
	46, 33,
	49, 33,
	-2, 69,
	-1, 109,
	46, 34,
	49, 34,
	-2, 69,
}

const yyNprod = 97
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 491

var yyAct = []int{

	69, 68, 95, 57, 50, 96, 90, 143, 26, 96,
	147, 51, 67, 142, 146, 65, 107, 64, 121, 61,
	166, 42, 13, 97, 59, 13, 168, 97, 13, 167,
	93, 13, 162, 121, 157, 111, 111, 60, 87, 88,
	89, 92, 114, 34, 18, 115, 99, 18, 174, 144,
	18, 41, 17, 18, 85, 17, 101, 108, 17, 106,
	110, 17, 159, 111, 113, 117, 58, 13, 55, 120,
	155, 56, 153, 123, 124, 125, 126, 127, 128, 129,
	130, 131, 132, 133, 134, 135, 136, 137, 59, 18,
	155, 173, 118, 13, 2, 59, 140, 17, 145, 112,
	18, 25, 149, 19, 47, 51, 59, 141, 109, 151,
	62, 54, 154, 13, 116, 18, 158, 165, 160, 156,
	164, 148, 163, 17, 102, 13, 70, 13, 83, 84,
	86, 45, 85, 104, 103, 18, 94, 4, 139, 91,
	100, 21, 66, 17, 73, 170, 169, 18, 163, 18,
	36, 35, 33, 53, 49, 17, 63, 17, 80, 72,
	12, 82, 81, 83, 84, 86, 122, 85, 52, 48,
	76, 77, 78, 79, 8, 175, 80, 74, 75, 82,
	81, 83, 84, 86, 73, 85, 138, 5, 46, 20,
	3, 1, 0, 171, 0, 0, 0, 0, 150, 72,
	152, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	76, 77, 78, 79, 0, 0, 80, 74, 75, 82,
	81, 83, 84, 86, 29, 85, 40, 0, 0, 0,
	0, 27, 37, 161, 0, 0, 0, 28, 0, 29,
	0, 40, 30, 31, 0, 0, 27, 37, 19, 32,
	0, 0, 28, 0, 39, 0, 0, 30, 31, 0,
	0, 0, 0, 19, 32, 38, 44, 172, 43, 39,
	0, 29, 0, 40, 0, 0, 0, 119, 27, 37,
	38, 44, 0, 43, 28, 0, 29, 0, 40, 30,
	31, 0, 0, 27, 37, 19, 32, 0, 0, 28,
	0, 39, 0, 0, 30, 31, 73, 0, 0, 0,
	19, 32, 38, 44, 98, 43, 39, 0, 0, 0,
	0, 72, 0, 0, 105, 0, 0, 38, 44, 0,
	43, 0, 76, 77, 78, 79, 0, 0, 80, 74,
	75, 82, 81, 83, 84, 86, 73, 85, 71, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 72, 73, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 76, 77, 78, 79, 0, 72, 80, 74,
	75, 82, 81, 83, 84, 86, 73, 85, 76, 77,
	78, 79, 0, 0, 80, 74, 75, 82, 81, 83,
	84, 86, 0, 85, 0, 0, 0, 0, 0, 0,
	0, 0, 76, 77, 78, 79, 0, 0, 80, 74,
	75, 82, 81, 83, 84, 86, 0, 85, 76, 77,
	78, 79, 0, 0, 80, 74, 75, 82, 81, 83,
	84, 86, 0, 85, 22, 6, 0, 0, 0, 0,
	16, 11, 0, 10, 0, 14, 0, 0, 0, 24,
	9, 0, 6, 0, 7, 15, 23, 16, 11, 0,
	10, 0, 14, 19, 0, 0, 0, 9, 0, 0,
	0, 7, 15, 0, 0, 0, 0, 0, 0, 0,
	19,
}
var yyPact = []int{

	456, -1000, -1000, 439, -1000, -1000, 456, 276, 125, 456,
	276, 77, 22, -1000, 54, -10, 76, -1000, -35, -1000,
	-1000, -1000, -1000, -1000, 276, 117, 342, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -35, -1000, -1000, 276, 276, 276,
	-6, -1000, -1000, -25, 261, 456, -1000, 32, 115, 126,
	-1000, 302, -6, -34, -1000, 276, 69, 14, 65, -1000,
	61, -13, -4, 99, 276, 58, -1000, 229, -31, 358,
	-1000, 456, 276, 276, 276, 276, 276, 276, 276, 276,
	276, 276, 276, 276, 276, 276, 276, 9, 358, 9,
	-1000, 456, 72, -1000, -42, -1000, 3, 276, -1000, -39,
	112, 276, -1000, 456, 276, 456, -1000, 38, -31, -1000,
	276, 36, -6, -14, 276, 28, 276, 180, -1000, -1000,
	-16, 276, 111, 382, 398, 122, 122, 122, 122, 122,
	122, 122, 87, 87, 9, 9, 9, 9, 108, -28,
	-20, -1000, -1000, -29, 276, 140, -1000, 214, -1000, 358,
	-1000, -1000, -1000, -1000, -31, -1000, -1000, -1000, -31, -1000,
	-31, -1000, -1000, 358, -1000, -1000, -1000, 56, -1000, -1000,
	358, 2, -1000, -1000, 276, 358,
}
var yyPgo = []int{

	0, 191, 94, 190, 189, 188, 0, 137, 187, 174,
	169, 168, 6, 160, 1, 21, 156, 154, 4, 3,
	51, 153, 152, 43, 151, 150, 142, 139, 138, 136,
	2,
}
var yyR1 = []int{

	0, 1, 2, 2, 5, 3, 3, 3, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 9, 9, 10,
	10, 17, 17, 18, 4, 4, 4, 4, 8, 8,
	8, 8, 8, 13, 13, 11, 21, 21, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 20, 20, 20, 23,
	23, 15, 26, 26, 22, 12, 27, 28, 28, 28,
	28, 16, 16, 19, 19, 25, 25, 25, 14, 14,
	24, 24, 24, 29, 29, 30, 30,
}
var yyR2 = []int{

	0, 1, 1, 2, 3, 1, 2, 0, 1, 3,
	5, 4, 2, 3, 3, 3, 1, 4, 4, 1,
	3, 1, 3, 3, 1, 1, 1, 2, 2, 4,
	4, 4, 2, 1, 3, 1, 1, 3, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 2, 2,
	2, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 1, 4, 3, 1,
	1, 2, 2, 3, 2, 3, 3, 1, 1, 3,
	0, 1, 3, 1, 3, 2, 3, 4, 1, 3,
	2, 3, 4, 1, 3, 3, 5,
}
var yyChk = []int{

	-1000, -1, -2, -3, -7, -8, 6, 25, -9, 21,
	14, 12, -13, -15, 16, 26, 11, -20, -23, 34,
	-4, -7, 5, 27, 20, -2, -6, 17, 23, 10,
	28, 29, 35, -22, -23, -24, -25, 18, 51, 40,
	12, -20, -15, 54, 52, 6, -5, -2, -10, -17,
	-18, -6, -11, -21, 34, 46, 49, -19, 12, 34,
	47, -19, 34, -16, 52, 50, -26, 47, -14, -6,
	9, 6, 19, 4, 37, 38, 30, 31, 32, 33,
	36, 40, 39, 41, 42, 45, 43, -6, -6, -6,
	-12, -27, 47, 55, -29, -30, 34, 52, 53, -14,
	-2, 24, 9, 8, 7, 22, -12, 50, -14, -20,
	46, 49, 34, -19, 46, 49, 15, -6, 34, 48,
	-14, 49, -2, -6, -6, -6, -6, -6, -6, -6,
	-6, -6, -6, -6, -6, -6, -6, -6, -2, -28,
	-19, 35, 55, 49, 46, -6, 53, 49, 9, -6,
	-2, -18, -2, 34, -14, 34, -12, 48, -14, 34,
	-14, 53, 48, -6, 9, 9, 48, 49, 55, -30,
	-6, 53, 53, 35, 46, -6,
}
var yyDef = []int{

	7, -2, 1, 2, 5, 8, 7, 0, 0, 7,
	0, 0, 0, -2, 0, 0, 0, -2, 0, 66,
	3, 6, 24, 25, 26, 0, 0, 38, 39, 40,
	41, 42, 43, 44, 45, 46, 47, 0, 0, 0,
	0, 69, 70, 0, 0, 7, 12, 0, 0, 19,
	21, 0, 0, 35, 36, 0, 0, 28, 0, 83,
	0, 32, 81, 0, 0, 0, 71, 0, 27, 88,
	9, 7, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 48, 49, 50,
	74, 7, 80, 90, 0, 93, 0, 0, 85, 0,
	0, 0, 13, 7, 0, 7, 14, 0, 15, -2,
	0, 0, 0, 0, 0, 0, 0, 0, 68, 72,
	0, 0, 0, 51, 52, 53, 54, 55, 56, 57,
	58, 59, 60, 61, 62, 63, 64, 65, 0, 0,
	77, 78, 91, 0, 0, 0, 86, 0, 11, 4,
	20, 22, 23, 37, 29, 84, 30, 31, 17, 82,
	18, 67, 73, 89, 10, 75, 76, 0, 92, 94,
	95, 0, 87, 79, 0, 96,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 51, 3, 43, 3, 3,
	47, 48, 41, 39, 49, 40, 50, 42, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	37, 46, 38, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 52, 3, 53, 45, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 54, 3, 55, 44,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36,
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
		//line lua.y:52
		{
			endChunk(yylex, &yyVAL)
		}
	case 3:
		//line lua.y:56
		{
			opAppend(yylex, &yyVAL, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 4:
		//line lua.y:59
		{
			op2(yylex, &yyVAL, OP_UNTIL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 6:
		//line lua.y:63
		{
			opAppend(yylex, &yyVAL, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 7:
		//line lua.y:64
		{
			yyVAL.value = nil
		}
	case 9:
		//line lua.y:68
		{
			yyVAL = yyS[yypt-1]
		}
	case 10:
		//line lua.y:69
		{
			op2(yylex, &yyVAL, OP_WHILE, &yyS[yypt-3], &yyS[yypt-1])
		}
	case 11:
		//line lua.y:70
		{
			opForBind(yylex, &yyVAL, &yyS[yypt-3], &yyS[yypt-1])
		}
	case 12:
		//line lua.y:71
		{
			yyVAL = yyS[yypt-0]
		}
	case 13:
		//line lua.y:72
		{
			yyVAL = yyS[yypt-1]
		}
	case 14:
		//line lua.y:73
		{
			bindFuncName(yylex, &yyS[yypt-0], &yyS[yypt-1], "")
			op2(yylex, &yyVAL, OP_ASSIGN, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 15:
		//line lua.y:77
		{
			op2(yylex, &yyVAL, OP_ASSIGN, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 17:
		//line lua.y:81
		{
			opFor(yylex, &yyVAL, OP_FOR, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 18:
		//line lua.y:82
		{
			opFor(yylex, &yyVAL, OP_FORIN, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 20:
		//line lua.y:86
		{
			opIf(yylex, &yyVAL, nil, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 22:
		//line lua.y:90
		{
			opIf(yylex, &yyVAL, nil, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 23:
		//line lua.y:93
		{
			opIf(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0], nil)
		}
	case 24:
		//line lua.y:96
		{
			op0(yylex, &yyVAL, OP_BREAK, &yyS[yypt-0])
		}
	case 25:
		//line lua.y:97
		{
			op0(yylex, &yyVAL, OP_CONTINUE, &yyS[yypt-0])
		}
	case 26:
		//line lua.y:98
		{
			op1(yylex, &yyVAL, OP_RETURN, nil)
		}
	case 27:
		//line lua.y:99
		{
			op1(yylex, &yyVAL, OP_RETURN, &yyS[yypt-0])
		}
	case 28:
		//line lua.y:102
		{
			opLocal(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 29:
		//line lua.y:103
		{
			opLocal(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 30:
		//line lua.y:104
		{
			bindFuncName(yylex, &yyS[yypt-0], nil, yyS[yypt-1].token.image)
			var tmp yySymType
			nameAppend(yylex, &tmp, &yyS[yypt-1], nil)
			opLocal(yylex, &yyVAL, &tmp, &yyS[yypt-0])
		}
	case 31:
		//line lua.y:110
		{
			opClosure(yylex, &yyVAL, &yyS[yypt-1])
		}
	case 32:
		//line lua.y:111
		{
			opClosure(yylex, &yyVAL, &yyS[yypt-0])
		}
	case 34:
		//line lua.y:115
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 36:
		//line lua.y:121
		{
			opVar(&yyVAL, &yyS[yypt-0])
		}
	case 37:
		//line lua.y:122
		{
			opValueExt(&yyS[yypt-0], yyS[yypt-0].token.image)
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 38:
		//line lua.y:128
		{
			opValue(yylex, &yyVAL)
		}
	case 39:
		//line lua.y:129
		{
			opValue(yylex, &yyVAL)
		}
	case 40:
		//line lua.y:130
		{
			opValue(yylex, &yyVAL)
		}
	case 41:
		//line lua.y:131
		{
			opValue(yylex, &yyVAL)
		}
	case 42:
		//line lua.y:132
		{
			opValue(yylex, &yyVAL)
		}
	case 43:
		//line lua.y:133
		{
			opVar(&yyVAL, &yyS[yypt-0])
		}
	case 48:
		//line lua.y:138
		{
			op1(yylex, &yyVAL, OP_NOT, &yyS[yypt-0])
		}
	case 49:
		//line lua.y:139
		{
			op1(yylex, &yyVAL, OP_LEN, &yyS[yypt-0])
		}
	case 50:
		//line lua.y:140
		{
			op1(yylex, &yyVAL, OP_NSIGN, &yyS[yypt-0])
		}
	case 51:
		//line lua.y:141
		{
			op2(yylex, &yyVAL, OP_OR, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 52:
		//line lua.y:142
		{
			op2(yylex, &yyVAL, OP_AND, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 53:
		//line lua.y:143
		{
			op2(yylex, &yyVAL, OP_LT, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 54:
		//line lua.y:144
		{
			op2(yylex, &yyVAL, OP_GT, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 55:
		//line lua.y:145
		{
			op2(yylex, &yyVAL, OP_LTEQ, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 56:
		//line lua.y:146
		{
			op2(yylex, &yyVAL, OP_GTEQ, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 57:
		//line lua.y:147
		{
			op2(yylex, &yyVAL, OP_EQ, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 58:
		//line lua.y:148
		{
			op2(yylex, &yyVAL, OP_NOTEQ, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 59:
		//line lua.y:149
		{
			op2(yylex, &yyVAL, OP_STRADD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 60:
		//line lua.y:150
		{
			op2(yylex, &yyVAL, OP_SUB, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 61:
		//line lua.y:151
		{
			op2(yylex, &yyVAL, OP_ADD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 62:
		//line lua.y:152
		{
			op2(yylex, &yyVAL, OP_MUL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 63:
		//line lua.y:153
		{
			op2(yylex, &yyVAL, OP_DIV, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 64:
		//line lua.y:154
		{
			op2(yylex, &yyVAL, OP_PMUL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 65:
		//line lua.y:155
		{
			op2(yylex, &yyVAL, OP_MOD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 66:
		//line lua.y:158
		{
			opVar(&yyVAL, &yyS[yypt-0])
		}
	case 67:
		//line lua.y:159
		{
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-3], &yyS[yypt-1])
		}
	case 68:
		//line lua.y:160
		{
			opValueExt(&yyS[yypt-0], yyS[yypt-0].token.image)
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 71:
		//line lua.y:170
		{
			op2(yylex, &yyVAL, OP_CALL, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 73:
		//line lua.y:174
		{
			yyVAL = yyS[yypt-1]
		}
	case 74:
		//line lua.y:177
		{
			yyVAL = yyS[yypt-0]
		}
	case 75:
		//line lua.y:180
		{
			opFunc(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-1])
		}
	case 76:
		//line lua.y:183
		{
			yyVAL = yyS[yypt-1]
		}
	case 78:
		//line lua.y:187
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 79:
		//line lua.y:188
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 80:
		//line lua.y:189
		{
			nameAppend(yylex, &yyVAL, nil, nil)
		}
	case 81:
		//line lua.y:193
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 82:
		//line lua.y:194
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 83:
		//line lua.y:197
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 84:
		//line lua.y:198
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 85:
		//line lua.y:201
		{
			op1(yylex, &yyVAL, OP_ARRAY, nil)
		}
	case 86:
		//line lua.y:202
		{
			op1(yylex, &yyVAL, OP_ARRAY, &yyS[yypt-1])
		}
	case 87:
		//line lua.y:203
		{
			op1(yylex, &yyVAL, OP_ARRAY, &yyS[yypt-2])
		}
	case 88:
		//line lua.y:206
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 89:
		//line lua.y:207
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 90:
		//line lua.y:210
		{
			op1(yylex, &yyVAL, OP_TABLE, nil)
		}
	case 91:
		//line lua.y:211
		{
			op1(yylex, &yyVAL, OP_TABLE, &yyS[yypt-1])
		}
	case 92:
		//line lua.y:212
		{
			op1(yylex, &yyVAL, OP_TABLE, &yyS[yypt-2])
		}
	case 94:
		//line lua.y:216
		{
			opAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 95:
		//line lua.y:219
		{
			opValueExt(&yyS[yypt-2], yyS[yypt-2].token.image)
			op2(yylex, &yyVAL, OP_FIELD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 96:
		//line lua.y:223
		{
			op2(yylex, &yyVAL, OP_FIELD, &yyS[yypt-3], &yyS[yypt-0])
		}
	}
	goto yystack /* stack new state and value */
}
