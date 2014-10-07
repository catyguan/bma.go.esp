//line lua.y:1
package goyacc

import __yyfmt__ "fmt"

//line lua.y:3
const AND = 57346
const BREAK = 57347
const DO = 57348
const ELSE = 57349
const ELSEIF = 57350
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
const NUMBER = 57368
const STRING = 57369
const SLTEQ = 57370
const SGTEQ = 57371
const SEQ = 57372
const SNOTEQ = 57373
const NAME = 57374
const MORE = 57375
const STRADD = 57376

var yyToknames = []string{
	"AND",
	"BREAK",
	"DO",
	"ELSE",
	"ELSEIF",
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

//line lua.y:216

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 13,
	46, 69,
	48, 69,
	50, 69,
	51, 69,
	-2, 16,
	-1, 16,
	44, 30,
	45, 30,
	-2, 68,
	-1, 108,
	44, 31,
	45, 31,
	-2, 68,
}

const yyNprod = 97
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 488

var yyAct = []int{

	66, 93, 65, 48, 62, 137, 88, 55, 24, 119,
	161, 49, 94, 64, 135, 138, 157, 90, 119, 119,
	40, 13, 141, 59, 13, 32, 17, 13, 95, 17,
	13, 162, 17, 113, 91, 17, 85, 86, 87, 94,
	39, 16, 106, 105, 16, 97, 168, 16, 139, 63,
	16, 60, 112, 61, 64, 95, 107, 104, 75, 76,
	78, 114, 77, 110, 13, 109, 110, 118, 77, 17,
	121, 122, 123, 124, 125, 126, 127, 128, 129, 130,
	17, 2, 53, 54, 16, 150, 167, 56, 23, 13,
	150, 45, 57, 134, 17, 108, 140, 148, 133, 147,
	143, 116, 115, 49, 111, 18, 145, 57, 58, 16,
	13, 52, 149, 153, 99, 17, 154, 160, 151, 159,
	158, 156, 13, 142, 13, 98, 43, 17, 72, 17,
	16, 74, 73, 75, 76, 78, 100, 77, 163, 67,
	164, 136, 16, 92, 16, 81, 82, 83, 84, 132,
	120, 72, 79, 80, 74, 73, 75, 76, 78, 4,
	77, 101, 102, 20, 89, 71, 34, 169, 33, 170,
	31, 131, 172, 27, 51, 38, 47, 152, 12, 50,
	25, 35, 46, 144, 8, 146, 26, 5, 44, 28,
	29, 27, 19, 38, 3, 18, 30, 1, 25, 35,
	0, 37, 0, 0, 26, 0, 0, 28, 29, 0,
	36, 42, 0, 18, 30, 117, 41, 0, 0, 37,
	0, 0, 0, 0, 27, 0, 38, 0, 36, 42,
	96, 25, 35, 0, 41, 0, 0, 26, 0, 70,
	28, 29, 0, 0, 0, 0, 18, 30, 0, 0,
	0, 0, 37, 0, 69, 0, 0, 0, 0, 0,
	0, 36, 42, 81, 82, 83, 84, 41, 70, 72,
	79, 80, 74, 73, 75, 76, 78, 0, 77, 0,
	0, 0, 0, 69, 165, 0, 0, 0, 0, 0,
	0, 70, 81, 82, 83, 84, 0, 0, 72, 79,
	80, 74, 73, 75, 76, 78, 69, 77, 0, 0,
	0, 0, 0, 155, 70, 81, 82, 83, 84, 0,
	0, 72, 79, 80, 74, 73, 75, 76, 78, 69,
	77, 0, 171, 0, 0, 0, 0, 0, 81, 82,
	83, 84, 70, 0, 72, 79, 80, 74, 73, 75,
	76, 78, 0, 77, 0, 166, 0, 69, 0, 0,
	103, 0, 0, 0, 0, 0, 81, 82, 83, 84,
	0, 0, 72, 79, 80, 74, 73, 75, 76, 78,
	70, 77, 68, 6, 0, 0, 0, 0, 15, 11,
	0, 10, 0, 14, 0, 69, 0, 0, 9, 0,
	0, 0, 7, 70, 81, 82, 83, 84, 0, 18,
	72, 79, 80, 74, 73, 75, 76, 78, 69, 77,
	70, 0, 0, 0, 0, 0, 0, 81, 82, 83,
	84, 0, 0, 72, 79, 80, 74, 73, 75, 76,
	78, 0, 77, 0, 81, 82, 83, 84, 0, 0,
	72, 79, 80, 74, 73, 75, 76, 78, 0, 77,
	21, 6, 0, 0, 0, 0, 15, 11, 0, 10,
	0, 14, 0, 0, 0, 22, 9, 0, 0, 0,
	7, 0, 0, 0, 0, 0, 0, 18,
}
var yyPact = []int{

	377, -1000, -1000, 455, -1000, -1000, 377, 214, 120, 377,
	214, 79, 38, -1000, 75, 76, -1000, 3, -1000, -1000,
	-1000, -1000, 214, 130, 376, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 3, -1000, -1000, 214, 214, 214, -34, -1000,
	-1000, -20, 181, 377, -1000, 90, 127, 154, -1000, 338,
	-34, -3, -1000, 214, 73, 21, 72, -1000, 8, 18,
	214, 70, -1000, 69, 163, -26, 399, -1000, 377, 214,
	214, 214, 214, 214, 214, 214, 214, 214, 214, -1000,
	-1000, -1000, -1000, -1000, -1000, 25, 399, 25, -1000, 377,
	60, -1000, -40, -1000, 4, 214, -1000, -27, 114, 214,
	-1000, 377, 214, 377, -1000, 67, 65, -26, -1000, 214,
	58, -34, 214, 214, 264, -1000, -38, -1000, -36, 214,
	110, 416, 117, 399, 94, 19, 19, 25, 25, 25,
	25, 108, -42, -14, -1000, -1000, 7, -1000, -1000, 214,
	235, -1000, -1000, 399, -1000, -1000, -1000, -1000, -1000, -26,
	-1000, -1000, -1000, 310, -26, -1000, -1000, -1000, 399, -1000,
	-1000, -1000, 53, -1000, 399, 2, 214, -1000, 214, 287,
	399, 214, 399,
}
var yyPgo = []int{

	0, 197, 81, 194, 192, 188, 0, 159, 187, 184,
	182, 179, 6, 178, 2, 20, 177, 7, 176, 3,
	40, 174, 170, 25, 168, 166, 165, 4, 164, 149,
	143, 1, 141,
}
var yyR1 = []int{

	0, 1, 2, 2, 5, 3, 3, 3, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 9, 9, 10,
	10, 18, 18, 19, 4, 4, 4, 8, 8, 8,
	13, 13, 11, 11, 21, 21, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 26,
	26, 26, 26, 26, 26, 20, 20, 20, 23, 23,
	15, 15, 27, 27, 22, 12, 28, 29, 29, 29,
	29, 17, 17, 25, 25, 14, 14, 16, 16, 24,
	24, 30, 30, 32, 32, 31, 31,
}
var yyR2 = []int{

	0, 1, 1, 2, 3, 1, 2, 0, 1, 3,
	5, 4, 2, 3, 3, 3, 1, 4, 4, 1,
	3, 1, 3, 3, 1, 1, 2, 2, 4, 4,
	1, 3, 1, 3, 1, 3, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 2, 2, 2, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 1,
	1, 1, 1, 1, 1, 1, 4, 3, 1, 1,
	2, 4, 2, 3, 2, 3, 3, 1, 1, 3,
	0, 1, 3, 2, 3, 1, 3, 3, 5, 2,
	3, 1, 3, 1, 1, 3, 5,
}
var yyChk = []int{

	-1000, -1, -2, -3, -7, -8, 6, 25, -9, 21,
	14, 12, -13, -15, 16, 11, -20, -23, 32, -4,
	-7, 5, 20, -2, -6, 17, 23, 10, 26, 27,
	33, -22, -23, -24, -25, 18, 47, 38, 12, -20,
	-15, 53, 48, 6, -5, -2, -10, -18, -19, -6,
	-11, -21, 32, 44, 45, -17, 12, 32, 32, -17,
	48, 50, -27, 46, 51, -14, -6, 9, 6, 19,
	4, -26, 34, 38, 37, 39, 40, 43, 41, 35,
	36, 28, 29, 30, 31, -6, -6, -6, -12, -28,
	51, 54, -30, -31, 32, 48, 49, -14, -2, 24,
	9, 7, 8, 22, -12, 46, 45, -14, -20, 44,
	45, 32, 44, 15, -6, 32, 32, 52, -14, 45,
	-2, -6, -6, -6, -6, -6, -6, -6, -6, -6,
	-6, -2, -29, -17, 33, 54, -32, 45, 55, 44,
	-6, 49, 9, -6, -2, -19, -2, 32, 32, -14,
	32, -12, -16, -6, -14, 49, -27, 52, -6, 9,
	9, 52, 45, -31, -6, 49, 45, 33, 44, -6,
	-6, 45, -6,
}
var yyDef = []int{

	7, -2, 1, 2, 5, 8, 7, 0, 0, 7,
	0, 0, 0, -2, 0, 0, -2, 0, 65, 3,
	6, 24, 25, 0, 0, 36, 37, 38, 39, 40,
	41, 42, 43, 44, 45, 0, 0, 0, 0, 68,
	69, 0, 0, 7, 12, 0, 0, 19, 21, 0,
	0, 32, 34, 0, 0, 27, 0, 81, 81, 0,
	0, 0, 70, 0, 0, 26, 85, 9, 7, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 59,
	60, 61, 62, 63, 64, 46, 47, 48, 74, 7,
	80, 89, 0, 91, 0, 0, 83, 0, 0, 0,
	13, 7, 0, 7, 14, 0, 0, 15, -2, 0,
	0, 0, 0, 0, 0, 67, 0, 72, 0, 0,
	0, 49, 50, 51, 52, 53, 54, 55, 56, 57,
	58, 0, 0, 77, 78, 90, 0, 93, 94, 0,
	0, 84, 11, 4, 20, 22, 23, 33, 35, 28,
	82, 29, 17, 0, 18, 66, 71, 73, 86, 10,
	75, 76, 0, 92, 95, 0, 0, 79, 0, 87,
	96, 0, 88,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 47, 3, 41, 3, 3,
	51, 52, 39, 37, 45, 38, 50, 40, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 46, 55,
	35, 44, 36, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 48, 3, 49, 43, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 53, 3, 54, 42,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34,
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
		//line lua.y:50
		{
			endChunk(yylex, &yyVAL)
		}
	case 3:
		//line lua.y:54
		{
			opAppend(yylex, &yyVAL, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 4:
		//line lua.y:57
		{
			op2(yylex, &yyVAL, OP_UNTIL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 6:
		//line lua.y:61
		{
			opAppend(yylex, &yyVAL, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 9:
		//line lua.y:66
		{
			yyVAL = yyS[yypt-1]
		}
	case 10:
		//line lua.y:67
		{
			op2(yylex, &yyVAL, OP_WHILE, &yyS[yypt-3], &yyS[yypt-1])
		}
	case 12:
		//line lua.y:69
		{
			yyVAL = yyS[yypt-0]
		}
	case 15:
		//line lua.y:72
		{
			op2(yylex, &yyVAL, OP_ASSIGN, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 25:
		//line lua.y:92
		{
			op1(yylex, &yyVAL, OP_RETURN, nil)
		}
	case 26:
		//line lua.y:93
		{
			op1(yylex, &yyVAL, OP_RETURN, &yyS[yypt-0])
		}
	case 27:
		//line lua.y:96
		{
			opLocal(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 28:
		//line lua.y:97
		{
			opLocal(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 36:
		//line lua.y:113
		{
			opValue(yylex, &yyVAL)
		}
	case 37:
		//line lua.y:114
		{
			opValue(yylex, &yyVAL)
		}
	case 38:
		//line lua.y:115
		{
			opValue(yylex, &yyVAL)
		}
	case 39:
		//line lua.y:116
		{
			opValue(yylex, &yyVAL)
		}
	case 40:
		//line lua.y:117
		{
			opValue(yylex, &yyVAL)
		}
	case 46:
		//line lua.y:123
		{
			op1(yylex, &yyVAL, OP_NOT, &yyS[yypt-0])
		}
	case 47:
		//line lua.y:124
		{
			op1(yylex, &yyVAL, OP_LEN, &yyS[yypt-0])
		}
	case 48:
		//line lua.y:125
		{
			op1(yylex, &yyVAL, OP_NSIGN, &yyS[yypt-0])
		}
	case 49:
		//line lua.y:126
		{
			op2(yylex, &yyVAL, OP_OR, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 50:
		//line lua.y:127
		{
			op2(yylex, &yyVAL, OP_AND, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 51:
		//line lua.y:128
		{
			op2(yylex, &yyVAL, yyS[yypt-1].op, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 52:
		//line lua.y:129
		{
			op2(yylex, &yyVAL, OP_STRADD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 53:
		//line lua.y:130
		{
			op2(yylex, &yyVAL, OP_SUB, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 54:
		//line lua.y:131
		{
			op2(yylex, &yyVAL, OP_ADD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 55:
		//line lua.y:132
		{
			op2(yylex, &yyVAL, OP_MUL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 56:
		//line lua.y:133
		{
			op2(yylex, &yyVAL, OP_DIV, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 57:
		//line lua.y:134
		{
			op2(yylex, &yyVAL, OP_PMUL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 58:
		//line lua.y:135
		{
			op2(yylex, &yyVAL, OP_MOD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 59:
		//line lua.y:138
		{
			opFlag(&yyVAL, OP_LT)
		}
	case 60:
		//line lua.y:139
		{
			opFlag(&yyVAL, OP_GT)
		}
	case 61:
		//line lua.y:140
		{
			opFlag(&yyVAL, OP_LTEQ)
		}
	case 62:
		//line lua.y:141
		{
			opFlag(&yyVAL, OP_GTEQ)
		}
	case 63:
		//line lua.y:142
		{
			opFlag(&yyVAL, OP_EQ)
		}
	case 64:
		//line lua.y:143
		{
			opFlag(&yyVAL, OP_NOTEQ)
		}
	case 65:
		//line lua.y:146
		{
			opVar(&yyVAL, &yyS[yypt-0])
		}
	case 66:
		//line lua.y:147
		{
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-3], &yyS[yypt-1])
		}
	case 67:
		//line lua.y:148
		{
			opValueExt(&yyS[yypt-0], yyS[yypt-0].token.image)
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 81:
		//line lua.y:182
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 82:
		//line lua.y:183
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 83:
		//line lua.y:186
		{
			op1(yylex, &yyVAL, OP_ARRAY, nil)
		}
	case 84:
		//line lua.y:187
		{
			op1(yylex, &yyVAL, OP_ARRAY, &yyS[yypt-1])
		}
	case 85:
		//line lua.y:190
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 86:
		//line lua.y:191
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 89:
		//line lua.y:198
		{
			op1(yylex, &yyVAL, OP_TABLE, nil)
		}
	case 90:
		//line lua.y:199
		{
			op1(yylex, &yyVAL, OP_TABLE, &yyS[yypt-1])
		}
	case 92:
		//line lua.y:203
		{
			opAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 95:
		//line lua.y:210
		{
			opValueExt(&yyS[yypt-2], yyS[yypt-2].token.image)
			op2(yylex, &yyVAL, OP_FIELD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 96:
		//line lua.y:214
		{
			op2(yylex, &yyVAL, OP_FIELD, &yyS[yypt-3], &yyS[yypt-0])
		}
	}
	goto yystack /* stack new state and value */
}
