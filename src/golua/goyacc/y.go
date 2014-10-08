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

//line lua.y:229

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 13,
	46, 69,
	47, 69,
	49, 69,
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

const yyNprod = 94
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 442

var yyAct = []int{

	66, 93, 65, 47, 62, 137, 88, 55, 24, 94,
	159, 49, 63, 61, 135, 60, 64, 64, 119, 90,
	40, 13, 113, 59, 13, 155, 95, 13, 119, 119,
	13, 91, 160, 140, 94, 165, 85, 86, 87, 32,
	17, 105, 106, 17, 138, 97, 17, 109, 110, 17,
	112, 95, 110, 53, 54, 77, 107, 104, 39, 16,
	56, 114, 16, 149, 13, 16, 147, 118, 16, 146,
	121, 122, 123, 124, 125, 126, 127, 128, 129, 130,
	57, 2, 116, 17, 149, 164, 57, 134, 23, 13,
	115, 45, 111, 18, 17, 58, 139, 52, 133, 99,
	142, 158, 16, 49, 157, 141, 144, 100, 17, 67,
	13, 101, 148, 108, 102, 151, 152, 43, 150, 4,
	156, 154, 13, 20, 13, 98, 136, 16, 72, 17,
	92, 74, 73, 75, 76, 78, 132, 77, 161, 162,
	89, 17, 71, 17, 34, 75, 76, 78, 16, 77,
	120, 27, 33, 38, 31, 51, 48, 12, 25, 35,
	16, 50, 16, 46, 26, 8, 166, 28, 29, 5,
	44, 131, 19, 18, 30, 3, 1, 0, 0, 37,
	0, 0, 0, 143, 27, 145, 38, 0, 0, 36,
	42, 25, 35, 117, 41, 0, 0, 26, 0, 0,
	28, 29, 0, 27, 0, 38, 18, 30, 0, 0,
	25, 35, 37, 0, 0, 0, 26, 0, 0, 28,
	29, 0, 36, 42, 96, 18, 30, 41, 70, 0,
	0, 37, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 36, 42, 69, 0, 0, 41, 0, 0, 0,
	0, 70, 81, 82, 83, 84, 0, 0, 72, 79,
	80, 74, 73, 75, 76, 78, 69, 77, 0, 0,
	0, 0, 0, 0, 163, 81, 82, 83, 84, 0,
	70, 72, 79, 80, 74, 73, 75, 76, 78, 0,
	77, 0, 0, 0, 0, 69, 0, 153, 103, 0,
	0, 0, 0, 0, 81, 82, 83, 84, 0, 0,
	72, 79, 80, 74, 73, 75, 76, 78, 70, 77,
	68, 6, 0, 0, 0, 0, 15, 11, 0, 10,
	0, 14, 0, 69, 0, 0, 9, 0, 0, 0,
	7, 70, 81, 82, 83, 84, 0, 18, 72, 79,
	80, 74, 73, 75, 76, 78, 69, 77, 70, 0,
	0, 0, 0, 0, 0, 81, 82, 83, 84, 0,
	0, 72, 79, 80, 74, 73, 75, 76, 78, 0,
	77, 0, 81, 82, 83, 84, 0, 0, 72, 79,
	80, 74, 73, 75, 76, 78, 0, 77, 81, 82,
	83, 84, 0, 0, 72, 79, 80, 74, 73, 75,
	76, 78, 0, 77, 21, 6, 0, 0, 0, 0,
	15, 11, 0, 10, 0, 14, 0, 0, 0, 22,
	9, 0, 0, 0, 7, 0, 0, 0, 0, 0,
	0, 18,
}
var yyPact = []int{

	315, -1000, -1000, 409, -1000, -1000, 315, 193, 111, 315,
	193, 65, 9, -1000, 48, 63, -1000, -34, -1000, -1000,
	-1000, -1000, 193, 100, 314, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -34, -1000, -1000, 193, 193, 193, -32, -1000,
	-1000, -23, 174, 315, -1000, 75, 98, 103, 107, 276,
	-32, -5, -1000, 193, 61, 3, 60, -1000, 6, 7,
	193, 58, -1000, 50, 141, -16, 337, -1000, 315, 193,
	193, 193, 193, 193, 193, 193, 193, 193, 193, -1000,
	-1000, -1000, -1000, -1000, -1000, 12, 337, 12, -1000, 315,
	54, -1000, -40, -1000, 0, 193, -1000, -17, 96, 193,
	-1000, 315, 193, 315, -1000, 37, 34, -16, -1000, 193,
	31, -32, 193, 193, 247, -1000, -35, -1000, -27, 193,
	95, 354, 370, 337, 94, 106, 106, 12, 12, 12,
	12, 92, -42, -13, -1000, -1000, 2, -1000, 193, 224,
	-1000, -1000, 337, -1000, -1000, -1000, -1000, -1000, -16, -1000,
	-1000, -16, -16, -1000, -1000, -1000, 337, -1000, -1000, -1000,
	52, -1000, 337, -9, -1000, 193, 337,
}
var yyPgo = []int{

	0, 176, 81, 175, 172, 170, 0, 119, 169, 165,
	163, 161, 6, 157, 2, 20, 7, 3, 156, 58,
	155, 154, 39, 152, 144, 142, 4, 140, 136, 130,
	1, 126,
}
var yyR1 = []int{

	0, 1, 2, 2, 5, 3, 3, 3, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 9, 9, 10,
	10, 17, 17, 18, 4, 4, 4, 8, 8, 8,
	13, 13, 11, 11, 20, 20, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 25,
	25, 25, 25, 25, 25, 19, 19, 19, 22, 22,
	15, 15, 26, 26, 21, 12, 27, 28, 28, 28,
	28, 16, 16, 24, 24, 14, 14, 23, 23, 29,
	29, 31, 30, 30,
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
	0, 1, 3, 2, 3, 1, 3, 2, 3, 1,
	3, 1, 3, 5,
}
var yyChk = []int{

	-1000, -1, -2, -3, -7, -8, 6, 25, -9, 21,
	14, 12, -13, -15, 16, 11, -19, -22, 32, -4,
	-7, 5, 20, -2, -6, 17, 23, 10, 26, 27,
	33, -21, -22, -23, -24, 18, 48, 38, 12, -19,
	-15, 53, 49, 6, -5, -2, -10, -17, -18, -6,
	-11, -20, 32, 44, 45, -16, 12, 32, 32, -16,
	49, 47, -26, 46, 51, -14, -6, 9, 6, 19,
	4, -25, 34, 38, 37, 39, 40, 43, 41, 35,
	36, 28, 29, 30, 31, -6, -6, -6, -12, -27,
	51, 54, -29, -30, 32, 49, 50, -14, -2, 24,
	9, 8, 7, 22, -12, 46, 47, -14, -19, 44,
	45, 32, 44, 15, -6, 32, 32, 52, -14, 45,
	-2, -6, -6, -6, -6, -6, -6, -6, -6, -6,
	-6, -2, -28, -16, 33, 54, -31, 45, 44, -6,
	50, 9, -6, -2, -17, -2, 32, 32, -14, 32,
	-12, -14, -14, 50, -26, 52, -6, 9, 9, 52,
	45, -30, -6, 50, 33, 44, -6,
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
	80, 87, 0, 89, 0, 0, 83, 0, 0, 0,
	13, 7, 0, 7, 14, 0, 0, 15, -2, 0,
	0, 0, 0, 0, 0, 67, 0, 72, 0, 0,
	0, 49, 50, 51, 52, 53, 54, 55, 56, 57,
	58, 0, 0, 77, 78, 88, 0, 91, 0, 0,
	84, 11, 4, 20, 22, 23, 33, 35, 28, 82,
	29, 17, 18, 66, 71, 73, 86, 10, 75, 76,
	0, 90, 92, 0, 79, 0, 93,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 48, 3, 41, 3, 3,
	51, 52, 39, 37, 45, 38, 47, 40, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 46, 3,
	35, 44, 36, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 49, 3, 50, 43, 3, 3, 3, 3, 3,
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
	case 13:
		//line lua.y:70
		{
			yyVAL = yyS[yypt-1]
		}
	case 14:
		//line lua.y:71
		{
			bindFuncName(yylex, &yyS[yypt-0], &yyS[yypt-1], "")
			op2(yylex, &yyVAL, OP_ASSIGN, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 15:
		//line lua.y:75
		{
			op2(yylex, &yyVAL, OP_ASSIGN, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 17:
		//line lua.y:79
		{
			opFor(yylex, &yyVAL, OP_FOR, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 18:
		//line lua.y:80
		{
			opFor(yylex, &yyVAL, OP_FORIN, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 20:
		//line lua.y:84
		{
			opIf(yylex, &yyVAL, nil, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 22:
		//line lua.y:88
		{
			opIf(yylex, &yyVAL, nil, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 23:
		//line lua.y:91
		{
			opIf(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0], nil)
		}
	case 25:
		//line lua.y:95
		{
			op1(yylex, &yyVAL, OP_RETURN, nil)
		}
	case 26:
		//line lua.y:96
		{
			op1(yylex, &yyVAL, OP_RETURN, &yyS[yypt-0])
		}
	case 27:
		//line lua.y:99
		{
			opLocal(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 28:
		//line lua.y:100
		{
			opLocal(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 29:
		//line lua.y:101
		{
			bindFuncName(yylex, &yyS[yypt-0], nil, yyS[yypt-1].token.image)
			var tmp yySymType
			nameAppend(yylex, &tmp, &yyS[yypt-1], nil)
			opLocal(yylex, &yyVAL, &tmp, &yyS[yypt-0])
		}
	case 31:
		//line lua.y:110
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 33:
		//line lua.y:114
		{
			opValueExt(&yyS[yypt-0], yyS[yypt-0].token.image)
			op2(yylex, &yyVAL, OP_SELFM, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 34:
		//line lua.y:120
		{
			opVar(&yyVAL, &yyS[yypt-0])
		}
	case 35:
		//line lua.y:121
		{
			opValueExt(&yyS[yypt-0], yyS[yypt-0].token.image)
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 36:
		//line lua.y:127
		{
			opValue(yylex, &yyVAL)
		}
	case 37:
		//line lua.y:128
		{
			opValue(yylex, &yyVAL)
		}
	case 38:
		//line lua.y:129
		{
			opValue(yylex, &yyVAL)
		}
	case 39:
		//line lua.y:130
		{
			opValue(yylex, &yyVAL)
		}
	case 40:
		//line lua.y:131
		{
			opValue(yylex, &yyVAL)
		}
	case 46:
		//line lua.y:137
		{
			op1(yylex, &yyVAL, OP_NOT, &yyS[yypt-0])
		}
	case 47:
		//line lua.y:138
		{
			op1(yylex, &yyVAL, OP_LEN, &yyS[yypt-0])
		}
	case 48:
		//line lua.y:139
		{
			op1(yylex, &yyVAL, OP_NSIGN, &yyS[yypt-0])
		}
	case 49:
		//line lua.y:140
		{
			op2(yylex, &yyVAL, OP_OR, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 50:
		//line lua.y:141
		{
			op2(yylex, &yyVAL, OP_AND, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 51:
		//line lua.y:142
		{
			op2(yylex, &yyVAL, yyS[yypt-1].op, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 52:
		//line lua.y:143
		{
			op2(yylex, &yyVAL, OP_STRADD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 53:
		//line lua.y:144
		{
			op2(yylex, &yyVAL, OP_SUB, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 54:
		//line lua.y:145
		{
			op2(yylex, &yyVAL, OP_ADD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 55:
		//line lua.y:146
		{
			op2(yylex, &yyVAL, OP_MUL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 56:
		//line lua.y:147
		{
			op2(yylex, &yyVAL, OP_DIV, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 57:
		//line lua.y:148
		{
			op2(yylex, &yyVAL, OP_PMUL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 58:
		//line lua.y:149
		{
			op2(yylex, &yyVAL, OP_MOD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 59:
		//line lua.y:152
		{
			opFlag(&yyVAL, OP_LT)
		}
	case 60:
		//line lua.y:153
		{
			opFlag(&yyVAL, OP_GT)
		}
	case 61:
		//line lua.y:154
		{
			opFlag(&yyVAL, OP_LTEQ)
		}
	case 62:
		//line lua.y:155
		{
			opFlag(&yyVAL, OP_GTEQ)
		}
	case 63:
		//line lua.y:156
		{
			opFlag(&yyVAL, OP_EQ)
		}
	case 64:
		//line lua.y:157
		{
			opFlag(&yyVAL, OP_NOTEQ)
		}
	case 65:
		//line lua.y:160
		{
			opVar(&yyVAL, &yyS[yypt-0])
		}
	case 66:
		//line lua.y:161
		{
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-3], &yyS[yypt-1])
		}
	case 67:
		//line lua.y:162
		{
			opValueExt(&yyS[yypt-0], yyS[yypt-0].token.image)
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 70:
		//line lua.y:172
		{
			op2(yylex, &yyVAL, OP_CALL, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 71:
		//line lua.y:173
		{
			var tmp yySymType
			op2(yylex, &tmp, OP_SELFM, &yyS[yypt-3], &yyS[yypt-1])
			op2(yylex, &yyVAL, OP_CALL, &tmp, &yyS[yypt-0])
		}
	case 73:
		//line lua.y:181
		{
			yyVAL = yyS[yypt-1]
		}
	case 74:
		//line lua.y:184
		{
			yyVAL = yyS[yypt-0]
		}
	case 75:
		//line lua.y:187
		{
			opFunc(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-1])
		}
	case 76:
		//line lua.y:190
		{
			yyVAL = yyS[yypt-1]
		}
	case 78:
		//line lua.y:194
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 79:
		//line lua.y:195
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 81:
		//line lua.y:200
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 82:
		//line lua.y:201
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 83:
		//line lua.y:204
		{
			op1(yylex, &yyVAL, OP_ARRAY, nil)
		}
	case 84:
		//line lua.y:205
		{
			op1(yylex, &yyVAL, OP_ARRAY, &yyS[yypt-1])
		}
	case 85:
		//line lua.y:208
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 86:
		//line lua.y:209
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 87:
		//line lua.y:212
		{
			op1(yylex, &yyVAL, OP_TABLE, nil)
		}
	case 88:
		//line lua.y:213
		{
			op1(yylex, &yyVAL, OP_TABLE, &yyS[yypt-1])
		}
	case 90:
		//line lua.y:217
		{
			opAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 92:
		//line lua.y:223
		{
			opValueExt(&yyS[yypt-2], yyS[yypt-2].token.image)
			op2(yylex, &yyVAL, OP_FIELD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 93:
		//line lua.y:227
		{
			op2(yylex, &yyVAL, OP_FIELD, &yyS[yypt-3], &yyS[yypt-0])
		}
	}
	goto yystack /* stack new state and value */
}
