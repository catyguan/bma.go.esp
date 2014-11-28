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
const ANNOTATION = 57379

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
	"ANNOTATION",
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

//line lua.y:228

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 14,
	48, 71,
	51, 71,
	53, 71,
	-2, 17,
	-1, 18,
	47, 34,
	50, 34,
	-2, 70,
	-1, 110,
	47, 35,
	50, 35,
	-2, 70,
}

const yyNprod = 98
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 497

var yyAct = []int{

	70, 69, 96, 51, 148, 144, 97, 91, 147, 27,
	58, 143, 52, 97, 68, 115, 108, 66, 116, 65,
	163, 122, 43, 14, 122, 98, 14, 62, 169, 167,
	14, 111, 98, 14, 112, 94, 158, 112, 81, 88,
	89, 90, 83, 82, 84, 85, 87, 100, 86, 35,
	19, 56, 160, 19, 57, 168, 112, 19, 109, 60,
	19, 107, 93, 84, 85, 87, 118, 86, 86, 14,
	121, 175, 114, 61, 124, 125, 126, 127, 128, 129,
	130, 131, 132, 133, 134, 135, 136, 137, 138, 2,
	42, 18, 145, 156, 18, 14, 19, 26, 18, 146,
	48, 18, 154, 150, 141, 119, 52, 19, 59, 152,
	156, 174, 60, 155, 113, 14, 20, 159, 63, 161,
	55, 157, 19, 164, 60, 142, 30, 14, 41, 14,
	60, 102, 117, 28, 38, 166, 101, 18, 165, 29,
	149, 103, 19, 71, 31, 32, 171, 170, 110, 164,
	20, 33, 105, 104, 19, 46, 19, 40, 95, 140,
	4, 92, 123, 18, 22, 67, 37, 36, 39, 45,
	173, 44, 74, 34, 54, 50, 176, 64, 13, 53,
	49, 9, 139, 18, 5, 47, 21, 73, 3, 1,
	0, 0, 0, 0, 151, 18, 153, 18, 77, 78,
	79, 80, 74, 0, 81, 0, 75, 76, 83, 82,
	84, 85, 87, 0, 86, 0, 0, 73, 0, 0,
	0, 0, 172, 0, 0, 0, 0, 0, 77, 78,
	79, 80, 0, 0, 81, 0, 75, 76, 83, 82,
	84, 85, 87, 30, 86, 41, 0, 30, 0, 41,
	28, 38, 162, 0, 28, 38, 29, 0, 0, 0,
	29, 31, 32, 0, 0, 31, 32, 20, 33, 0,
	0, 20, 33, 0, 40, 0, 0, 0, 40, 0,
	0, 0, 120, 0, 0, 39, 45, 0, 44, 39,
	45, 99, 44, 30, 0, 41, 0, 0, 0, 0,
	28, 38, 0, 0, 0, 74, 29, 0, 0, 0,
	0, 31, 32, 0, 0, 0, 0, 20, 33, 0,
	73, 0, 0, 106, 40, 0, 0, 0, 0, 0,
	0, 77, 78, 79, 80, 39, 45, 81, 44, 75,
	76, 83, 82, 84, 85, 87, 74, 86, 72, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 73, 0, 0, 74, 0, 0, 0, 0, 0,
	0, 0, 77, 78, 79, 80, 0, 0, 81, 73,
	75, 76, 83, 82, 84, 85, 87, 74, 86, 0,
	77, 78, 79, 80, 0, 0, 81, 0, 75, 76,
	83, 82, 84, 85, 87, 0, 86, 0, 0, 0,
	0, 0, 0, 77, 78, 79, 80, 0, 0, 81,
	0, 75, 76, 83, 82, 84, 85, 87, 0, 86,
	77, 78, 79, 80, 0, 0, 81, 0, 75, 76,
	83, 82, 84, 85, 87, 0, 86, 23, 7, 0,
	0, 0, 0, 17, 12, 0, 11, 0, 15, 0,
	0, 0, 25, 10, 0, 7, 0, 8, 16, 24,
	17, 12, 0, 11, 0, 15, 20, 0, 0, 6,
	10, 0, 0, 0, 8, 16, 0, 0, 0, 0,
	0, 0, 0, 20, 0, 0, 6,
}
var yyPact = []int{

	459, -1000, -1000, 442, -1000, -1000, -1000, 459, 283, 149,
	459, 283, 86, 4, -1000, 96, 25, 84, -1000, -34,
	-1000, -1000, -1000, -1000, -1000, 283, 134, 342, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -34, -1000, -1000, 283, 283,
	283, 14, -1000, -1000, -21, 237, 459, -1000, 107, 132,
	145, -1000, 301, 14, -35, -1000, 283, 82, -16, 80,
	-1000, 78, 6, -32, 117, 283, 71, -1000, 233, -26,
	360, -1000, 459, 283, 283, 283, 283, 283, 283, 283,
	283, 283, 283, 283, 283, 283, 283, 283, 22, 360,
	22, -1000, 459, 90, -1000, -45, -1000, 45, 283, -1000,
	-46, 131, 283, -1000, 459, 283, 459, -1000, 68, -26,
	-1000, 283, 59, 14, -13, 283, 18, 283, 198, -1000,
	-1000, -29, 283, 129, 383, 400, 2, 2, 2, 2,
	2, 2, 2, 21, 21, 22, 22, 22, 22, 126,
	-20, 5, -1000, -1000, -28, 283, 168, -1000, 116, -1000,
	360, -1000, -1000, -1000, -1000, -26, -1000, -1000, -1000, -26,
	-1000, -26, -1000, -1000, 360, -1000, -1000, -1000, 76, -1000,
	-1000, 360, 24, -1000, -1000, 283, 360,
}
var yyPgo = []int{

	0, 189, 89, 188, 186, 185, 0, 160, 184, 181,
	180, 179, 7, 178, 1, 22, 177, 175, 3, 10,
	90, 174, 173, 49, 167, 166, 165, 161, 159, 158,
	2,
}
var yyR1 = []int{

	0, 1, 2, 2, 5, 3, 3, 3, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 9, 9,
	10, 10, 17, 17, 18, 4, 4, 4, 4, 8,
	8, 8, 8, 8, 13, 13, 11, 21, 21, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 20, 20, 20,
	23, 23, 15, 26, 26, 22, 12, 27, 28, 28,
	28, 28, 16, 16, 19, 19, 25, 25, 25, 14,
	14, 24, 24, 24, 29, 29, 30, 30,
}
var yyR2 = []int{

	0, 1, 1, 2, 3, 1, 2, 0, 1, 1,
	3, 5, 4, 2, 3, 3, 3, 1, 4, 4,
	1, 3, 1, 3, 3, 1, 1, 1, 2, 2,
	4, 4, 4, 2, 1, 3, 1, 1, 3, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 2,
	2, 2, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 1, 4, 3,
	1, 1, 2, 2, 3, 2, 3, 3, 1, 1,
	3, 0, 1, 3, 1, 3, 2, 3, 4, 1,
	3, 2, 3, 4, 1, 3, 3, 5,
}
var yyChk = []int{

	-1000, -1, -2, -3, -7, -8, 37, 6, 25, -9,
	21, 14, 12, -13, -15, 16, 26, 11, -20, -23,
	34, -4, -7, 5, 27, 20, -2, -6, 17, 23,
	10, 28, 29, 35, -22, -23, -24, -25, 18, 52,
	41, 12, -20, -15, 55, 53, 6, -5, -2, -10,
	-17, -18, -6, -11, -21, 34, 47, 50, -19, 12,
	34, 48, -19, 34, -16, 53, 51, -26, 48, -14,
	-6, 9, 6, 19, 4, 38, 39, 30, 31, 32,
	33, 36, 41, 40, 42, 43, 46, 44, -6, -6,
	-6, -12, -27, 48, 56, -29, -30, 34, 53, 54,
	-14, -2, 24, 9, 8, 7, 22, -12, 51, -14,
	-20, 47, 50, 34, -19, 47, 50, 15, -6, 34,
	49, -14, 50, -2, -6, -6, -6, -6, -6, -6,
	-6, -6, -6, -6, -6, -6, -6, -6, -6, -2,
	-28, -19, 35, 56, 50, 47, -6, 54, 50, 9,
	-6, -2, -18, -2, 34, -14, 34, -12, 49, -14,
	34, -14, 54, 49, -6, 9, 9, 49, 50, 56,
	-30, -6, 54, 54, 35, 47, -6,
}
var yyDef = []int{

	7, -2, 1, 2, 5, 8, 9, 7, 0, 0,
	7, 0, 0, 0, -2, 0, 0, 0, -2, 0,
	67, 3, 6, 25, 26, 27, 0, 0, 39, 40,
	41, 42, 43, 44, 45, 46, 47, 48, 0, 0,
	0, 0, 70, 71, 0, 0, 7, 13, 0, 0,
	20, 22, 0, 0, 36, 37, 0, 0, 29, 0,
	84, 0, 33, 82, 0, 0, 0, 72, 0, 28,
	89, 10, 7, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 49, 50,
	51, 75, 7, 81, 91, 0, 94, 0, 0, 86,
	0, 0, 0, 14, 7, 0, 7, 15, 0, 16,
	-2, 0, 0, 0, 0, 0, 0, 0, 0, 69,
	73, 0, 0, 0, 52, 53, 54, 55, 56, 57,
	58, 59, 60, 61, 62, 63, 64, 65, 66, 0,
	0, 78, 79, 92, 0, 0, 0, 87, 0, 12,
	4, 21, 23, 24, 38, 30, 85, 31, 32, 18,
	83, 19, 68, 74, 90, 11, 76, 77, 0, 93,
	95, 96, 0, 88, 80, 0, 97,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 52, 3, 44, 3, 3,
	48, 49, 42, 40, 50, 41, 51, 43, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	38, 47, 39, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 53, 3, 54, 46, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 55, 3, 56, 45,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37,
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
		//line lua.y:54
		{
			endChunk(yylex, &yyVAL)
		}
	case 3:
		//line lua.y:58
		{
			opAppend(yylex, &yyVAL, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 4:
		//line lua.y:61
		{
			op2(yylex, &yyVAL, OP_UNTIL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 6:
		//line lua.y:65
		{
			opAppend(yylex, &yyVAL, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 7:
		//line lua.y:66
		{
			yyVAL.value = nil
		}
	case 9:
		//line lua.y:70
		{
			defineAnnotation(yylex, &yyVAL)
		}
	case 10:
		//line lua.y:71
		{
			yyVAL = yyS[yypt-1]
		}
	case 11:
		//line lua.y:72
		{
			op2(yylex, &yyVAL, OP_WHILE, &yyS[yypt-3], &yyS[yypt-1])
		}
	case 12:
		//line lua.y:73
		{
			opForBind(yylex, &yyVAL, &yyS[yypt-3], &yyS[yypt-1])
		}
	case 13:
		//line lua.y:74
		{
			yyVAL = yyS[yypt-0]
		}
	case 14:
		//line lua.y:75
		{
			yyVAL = yyS[yypt-1]
		}
	case 15:
		//line lua.y:76
		{
			bindFuncName(yylex, &yyS[yypt-0], &yyS[yypt-1], "")
			op2(yylex, &yyVAL, OP_ASSIGN, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 16:
		//line lua.y:80
		{
			op2(yylex, &yyVAL, OP_ASSIGN, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 18:
		//line lua.y:84
		{
			opFor(yylex, &yyVAL, OP_FOR, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 19:
		//line lua.y:85
		{
			opFor(yylex, &yyVAL, OP_FORIN, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 21:
		//line lua.y:89
		{
			opIf(yylex, &yyVAL, nil, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 23:
		//line lua.y:93
		{
			opIf(yylex, &yyVAL, nil, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 24:
		//line lua.y:96
		{
			opIf(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0], nil)
		}
	case 25:
		//line lua.y:99
		{
			op0(yylex, &yyVAL, OP_BREAK, &yyS[yypt-0])
		}
	case 26:
		//line lua.y:100
		{
			op0(yylex, &yyVAL, OP_CONTINUE, &yyS[yypt-0])
		}
	case 27:
		//line lua.y:101
		{
			op1(yylex, &yyVAL, OP_RETURN, nil)
		}
	case 28:
		//line lua.y:102
		{
			op1(yylex, &yyVAL, OP_RETURN, &yyS[yypt-0])
		}
	case 29:
		//line lua.y:105
		{
			opLocal(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 30:
		//line lua.y:106
		{
			opLocal(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 31:
		//line lua.y:107
		{
			bindFuncName(yylex, &yyS[yypt-0], nil, yyS[yypt-1].token.image)
			var tmp yySymType
			nameAppend(yylex, &tmp, &yyS[yypt-1], nil)
			opLocal(yylex, &yyVAL, &tmp, &yyS[yypt-0])
		}
	case 32:
		//line lua.y:113
		{
			opClosure(yylex, &yyVAL, &yyS[yypt-1])
		}
	case 33:
		//line lua.y:114
		{
			opClosure(yylex, &yyVAL, &yyS[yypt-0])
		}
	case 35:
		//line lua.y:118
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 37:
		//line lua.y:124
		{
			opVar(&yyVAL, &yyS[yypt-0])
		}
	case 38:
		//line lua.y:125
		{
			opValueExt(&yyS[yypt-0], yyS[yypt-0].token.image)
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-2], &yyS[yypt-0])
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
	case 42:
		//line lua.y:134
		{
			opValue(yylex, &yyVAL)
		}
	case 43:
		//line lua.y:135
		{
			opValue(yylex, &yyVAL)
		}
	case 44:
		//line lua.y:136
		{
			opVar(&yyVAL, &yyS[yypt-0])
		}
	case 49:
		//line lua.y:141
		{
			op1(yylex, &yyVAL, OP_NOT, &yyS[yypt-0])
		}
	case 50:
		//line lua.y:142
		{
			op1(yylex, &yyVAL, OP_LEN, &yyS[yypt-0])
		}
	case 51:
		//line lua.y:143
		{
			op1(yylex, &yyVAL, OP_NSIGN, &yyS[yypt-0])
		}
	case 52:
		//line lua.y:144
		{
			op2(yylex, &yyVAL, OP_OR, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 53:
		//line lua.y:145
		{
			op2(yylex, &yyVAL, OP_AND, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 54:
		//line lua.y:146
		{
			op2(yylex, &yyVAL, OP_LT, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 55:
		//line lua.y:147
		{
			op2(yylex, &yyVAL, OP_GT, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 56:
		//line lua.y:148
		{
			op2(yylex, &yyVAL, OP_LTEQ, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 57:
		//line lua.y:149
		{
			op2(yylex, &yyVAL, OP_GTEQ, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 58:
		//line lua.y:150
		{
			op2(yylex, &yyVAL, OP_EQ, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 59:
		//line lua.y:151
		{
			op2(yylex, &yyVAL, OP_NOTEQ, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 60:
		//line lua.y:152
		{
			op2(yylex, &yyVAL, OP_STRADD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 61:
		//line lua.y:153
		{
			op2(yylex, &yyVAL, OP_SUB, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 62:
		//line lua.y:154
		{
			op2(yylex, &yyVAL, OP_ADD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 63:
		//line lua.y:155
		{
			op2(yylex, &yyVAL, OP_MUL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 64:
		//line lua.y:156
		{
			op2(yylex, &yyVAL, OP_DIV, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 65:
		//line lua.y:157
		{
			op2(yylex, &yyVAL, OP_PMUL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 66:
		//line lua.y:158
		{
			op2(yylex, &yyVAL, OP_MOD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 67:
		//line lua.y:161
		{
			opVar(&yyVAL, &yyS[yypt-0])
		}
	case 68:
		//line lua.y:162
		{
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-3], &yyS[yypt-1])
		}
	case 69:
		//line lua.y:163
		{
			opValueExt(&yyS[yypt-0], yyS[yypt-0].token.image)
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 72:
		//line lua.y:173
		{
			op2(yylex, &yyVAL, OP_CALL, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 74:
		//line lua.y:177
		{
			yyVAL = yyS[yypt-1]
		}
	case 75:
		//line lua.y:180
		{
			yyVAL = yyS[yypt-0]
		}
	case 76:
		//line lua.y:183
		{
			opFunc(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-1])
		}
	case 77:
		//line lua.y:186
		{
			yyVAL = yyS[yypt-1]
		}
	case 79:
		//line lua.y:190
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 80:
		//line lua.y:191
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 81:
		//line lua.y:192
		{
			nameAppend(yylex, &yyVAL, nil, nil)
		}
	case 82:
		//line lua.y:196
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 83:
		//line lua.y:197
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 84:
		//line lua.y:200
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 85:
		//line lua.y:201
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 86:
		//line lua.y:204
		{
			op1(yylex, &yyVAL, OP_ARRAY, nil)
		}
	case 87:
		//line lua.y:205
		{
			op1(yylex, &yyVAL, OP_ARRAY, &yyS[yypt-1])
		}
	case 88:
		//line lua.y:206
		{
			op1(yylex, &yyVAL, OP_ARRAY, &yyS[yypt-2])
		}
	case 89:
		//line lua.y:209
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 90:
		//line lua.y:210
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 91:
		//line lua.y:213
		{
			op1(yylex, &yyVAL, OP_TABLE, nil)
		}
	case 92:
		//line lua.y:214
		{
			op1(yylex, &yyVAL, OP_TABLE, &yyS[yypt-1])
		}
	case 93:
		//line lua.y:215
		{
			op1(yylex, &yyVAL, OP_TABLE, &yyS[yypt-2])
		}
	case 95:
		//line lua.y:219
		{
			opAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 96:
		//line lua.y:222
		{
			opValueExt(&yyS[yypt-2], yyS[yypt-2].token.image)
			op2(yylex, &yyVAL, OP_FIELD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 97:
		//line lua.y:226
		{
			op2(yylex, &yyVAL, OP_FIELD, &yyS[yypt-3], &yyS[yypt-0])
		}
	}
	goto yystack /* stack new state and value */
}
