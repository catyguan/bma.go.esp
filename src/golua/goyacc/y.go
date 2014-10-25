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

//line lua.y:226

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 13,
	46, 70,
	49, 70,
	51, 70,
	-2, 16,
	-1, 17,
	45, 32,
	48, 32,
	-2, 69,
	-1, 109,
	45, 33,
	48, 33,
	-2, 69,
}

const yyNprod = 97
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 482

var yyAct = []int{

	68, 67, 95, 49, 107, 142, 96, 90, 25, 141,
	56, 50, 96, 121, 138, 33, 18, 157, 121, 18,
	137, 2, 18, 162, 97, 18, 60, 163, 24, 66,
	97, 46, 64, 93, 63, 152, 111, 87, 88, 89,
	40, 17, 111, 114, 17, 99, 115, 17, 74, 161,
	17, 76, 75, 77, 78, 80, 108, 79, 92, 106,
	18, 169, 110, 139, 117, 111, 100, 79, 120, 58,
	113, 18, 123, 124, 125, 126, 127, 128, 129, 130,
	131, 132, 59, 41, 13, 17, 18, 13, 54, 57,
	13, 55, 122, 13, 150, 168, 109, 154, 140, 58,
	136, 150, 144, 135, 148, 50, 118, 18, 146, 58,
	58, 17, 149, 133, 112, 19, 153, 61, 155, 18,
	151, 18, 158, 53, 72, 145, 101, 147, 13, 77,
	78, 80, 17, 79, 116, 160, 159, 143, 102, 71,
	165, 164, 69, 158, 17, 44, 17, 104, 103, 83,
	84, 85, 86, 94, 13, 74, 81, 82, 76, 75,
	77, 78, 80, 4, 79, 72, 134, 21, 91, 65,
	170, 73, 166, 35, 34, 13, 32, 52, 48, 62,
	71, 12, 51, 47, 8, 5, 45, 13, 20, 13,
	83, 84, 85, 86, 3, 1, 74, 81, 82, 76,
	75, 77, 78, 80, 28, 79, 39, 0, 0, 0,
	0, 26, 36, 156, 0, 0, 0, 27, 0, 0,
	0, 29, 30, 0, 0, 0, 0, 19, 31, 0,
	0, 0, 0, 38, 0, 28, 0, 39, 0, 0,
	0, 0, 26, 36, 37, 43, 167, 42, 27, 0,
	0, 0, 29, 30, 0, 0, 0, 0, 19, 31,
	0, 0, 0, 0, 38, 0, 28, 0, 39, 0,
	0, 0, 119, 26, 36, 37, 43, 0, 42, 27,
	0, 0, 0, 29, 30, 0, 0, 0, 0, 19,
	31, 0, 0, 0, 0, 38, 0, 28, 0, 39,
	0, 0, 0, 0, 26, 36, 37, 43, 98, 42,
	27, 0, 0, 0, 29, 30, 72, 0, 0, 0,
	19, 31, 0, 0, 0, 0, 38, 0, 0, 0,
	0, 71, 0, 0, 105, 0, 0, 37, 43, 0,
	42, 83, 84, 85, 86, 0, 0, 74, 81, 82,
	76, 75, 77, 78, 80, 72, 79, 70, 6, 0,
	0, 0, 0, 16, 11, 0, 10, 0, 14, 0,
	71, 0, 0, 9, 0, 0, 0, 7, 15, 72,
	83, 84, 85, 86, 0, 19, 74, 81, 82, 76,
	75, 77, 78, 80, 71, 79, 72, 0, 0, 0,
	0, 0, 0, 0, 83, 84, 85, 86, 0, 0,
	74, 81, 82, 76, 75, 77, 78, 80, 0, 79,
	0, 83, 84, 85, 86, 0, 0, 74, 81, 82,
	76, 75, 77, 78, 80, 0, 79, 83, 84, 85,
	86, 0, 0, 74, 81, 82, 76, 75, 77, 78,
	80, 0, 79, 22, 6, 0, 0, 0, 0, 16,
	11, 0, 10, 0, 14, 0, 0, 0, 23, 9,
	0, 0, 0, 7, 15, 0, 0, 0, 0, 0,
	0, 19,
}
var yyPact = []int{

	352, -1000, -1000, 448, -1000, -1000, 352, 287, 139, 352,
	287, 90, 43, -1000, 77, 36, 84, -1000, -17, -1000,
	-1000, -1000, -1000, 287, 133, 351, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -17, -1000, -1000, 287, 287, 287, 12,
	-1000, -1000, -21, 256, 352, -1000, 102, 129, 140, -1000,
	312, 12, -45, -1000, 287, 82, 17, 81, -1000, 76,
	-6, -2, 119, 287, 73, -1000, 225, -35, 375, -1000,
	352, 287, 287, 287, 287, 287, 287, 287, 287, 287,
	287, -1000, -1000, -1000, -1000, -1000, -1000, 23, 375, 23,
	-1000, 352, 66, -1000, -34, -1000, 18, 287, -1000, -43,
	128, 287, -1000, 352, 287, 352, -1000, 71, -35, -1000,
	287, 68, 12, -12, 287, 64, 287, 161, -1000, -1000,
	-30, 287, 127, 392, 408, 375, 13, 89, 89, 23,
	23, 23, 23, 126, 2, -25, -1000, -1000, -27, 287,
	120, -1000, 194, -1000, 375, -1000, -1000, -1000, -1000, -35,
	-1000, -1000, -1000, -35, -1000, -35, -1000, -1000, 375, -1000,
	-1000, -1000, 61, -1000, -1000, 375, 16, -1000, -1000, 287,
	375,
}
var yyPgo = []int{

	0, 195, 21, 194, 188, 186, 0, 163, 185, 184,
	183, 182, 7, 181, 1, 83, 179, 178, 3, 10,
	40, 177, 176, 15, 174, 173, 171, 169, 168, 166,
	153, 2,
}
var yyR1 = []int{

	0, 1, 2, 2, 5, 3, 3, 3, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 9, 9, 10,
	10, 17, 17, 18, 4, 4, 4, 8, 8, 8,
	8, 8, 13, 13, 11, 21, 21, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	26, 26, 26, 26, 26, 26, 20, 20, 20, 23,
	23, 15, 27, 27, 22, 12, 28, 29, 29, 29,
	29, 16, 16, 19, 19, 25, 25, 25, 14, 14,
	24, 24, 24, 30, 30, 31, 31,
}
var yyR2 = []int{

	0, 1, 1, 2, 3, 1, 2, 0, 1, 3,
	5, 4, 2, 3, 3, 3, 1, 4, 4, 1,
	3, 1, 3, 3, 1, 1, 2, 2, 4, 4,
	4, 2, 1, 3, 1, 1, 3, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 2, 2, 2,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	1, 1, 1, 1, 1, 1, 1, 4, 3, 1,
	1, 2, 2, 3, 2, 3, 3, 1, 1, 3,
	0, 1, 3, 1, 3, 2, 3, 4, 1, 3,
	2, 3, 4, 1, 3, 3, 5,
}
var yyChk = []int{

	-1000, -1, -2, -3, -7, -8, 6, 25, -9, 21,
	14, 12, -13, -15, 16, 26, 11, -20, -23, 33,
	-4, -7, 5, 20, -2, -6, 17, 23, 10, 27,
	28, 34, -22, -23, -24, -25, 18, 50, 39, 12,
	-20, -15, 53, 51, 6, -5, -2, -10, -17, -18,
	-6, -11, -21, 33, 45, 48, -19, 12, 33, 46,
	-19, 33, -16, 51, 49, -27, 46, -14, -6, 9,
	6, 19, 4, -26, 35, 39, 38, 40, 41, 44,
	42, 36, 37, 29, 30, 31, 32, -6, -6, -6,
	-12, -28, 46, 54, -30, -31, 33, 51, 52, -14,
	-2, 24, 9, 8, 7, 22, -12, 49, -14, -20,
	45, 48, 33, -19, 45, 48, 15, -6, 33, 47,
	-14, 48, -2, -6, -6, -6, -6, -6, -6, -6,
	-6, -6, -6, -2, -29, -19, 34, 54, 48, 45,
	-6, 52, 48, 9, -6, -2, -18, -2, 33, -14,
	33, -12, 47, -14, 33, -14, 52, 47, -6, 9,
	9, 47, 48, 54, -31, -6, 52, 52, 34, 45,
	-6,
}
var yyDef = []int{

	7, -2, 1, 2, 5, 8, 7, 0, 0, 7,
	0, 0, 0, -2, 0, 0, 0, -2, 0, 66,
	3, 6, 24, 25, 0, 0, 37, 38, 39, 40,
	41, 42, 43, 44, 45, 46, 0, 0, 0, 0,
	69, 70, 0, 0, 7, 12, 0, 0, 19, 21,
	0, 0, 34, 35, 0, 0, 27, 0, 83, 0,
	31, 81, 0, 0, 0, 71, 0, 26, 88, 9,
	7, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 60, 61, 62, 63, 64, 65, 47, 48, 49,
	74, 7, 80, 90, 0, 93, 0, 0, 85, 0,
	0, 0, 13, 7, 0, 7, 14, 0, 15, -2,
	0, 0, 0, 0, 0, 0, 0, 0, 68, 72,
	0, 0, 0, 50, 51, 52, 53, 54, 55, 56,
	57, 58, 59, 0, 0, 77, 78, 91, 0, 0,
	0, 86, 0, 11, 4, 20, 22, 23, 36, 28,
	84, 29, 30, 17, 82, 18, 67, 73, 89, 10,
	75, 76, 0, 92, 94, 95, 0, 87, 79, 0,
	96,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 50, 3, 42, 3, 3,
	46, 47, 40, 38, 48, 39, 49, 41, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	36, 45, 37, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 51, 3, 52, 44, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 53, 3, 54, 43,
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
	case 7:
		//line lua.y:63
		{
			yyVAL.value = nil
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
	case 24:
		//line lua.y:95
		{
			op0(yylex, &yyVAL, OP_BREAK, &yyS[yypt-0])
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
	case 31:
		//line lua.y:109
		{
			opClosure(yylex, &yyVAL, &yyS[yypt-0])
		}
	case 33:
		//line lua.y:113
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 35:
		//line lua.y:119
		{
			opVar(&yyVAL, &yyS[yypt-0])
		}
	case 36:
		//line lua.y:120
		{
			opValueExt(&yyS[yypt-0], yyS[yypt-0].token.image)
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 37:
		//line lua.y:126
		{
			opValue(yylex, &yyVAL)
		}
	case 38:
		//line lua.y:127
		{
			opValue(yylex, &yyVAL)
		}
	case 39:
		//line lua.y:128
		{
			opValue(yylex, &yyVAL)
		}
	case 40:
		//line lua.y:129
		{
			opValue(yylex, &yyVAL)
		}
	case 41:
		//line lua.y:130
		{
			opValue(yylex, &yyVAL)
		}
	case 42:
		//line lua.y:131
		{
			opVar(&yyVAL, &yyS[yypt-0])
		}
	case 47:
		//line lua.y:136
		{
			op1(yylex, &yyVAL, OP_NOT, &yyS[yypt-0])
		}
	case 48:
		//line lua.y:137
		{
			op1(yylex, &yyVAL, OP_LEN, &yyS[yypt-0])
		}
	case 49:
		//line lua.y:138
		{
			op1(yylex, &yyVAL, OP_NSIGN, &yyS[yypt-0])
		}
	case 50:
		//line lua.y:139
		{
			op2(yylex, &yyVAL, OP_OR, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 51:
		//line lua.y:140
		{
			op2(yylex, &yyVAL, OP_AND, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 52:
		//line lua.y:141
		{
			op2(yylex, &yyVAL, yyS[yypt-1].op, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 53:
		//line lua.y:142
		{
			op2(yylex, &yyVAL, OP_STRADD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 54:
		//line lua.y:143
		{
			op2(yylex, &yyVAL, OP_SUB, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 55:
		//line lua.y:144
		{
			op2(yylex, &yyVAL, OP_ADD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 56:
		//line lua.y:145
		{
			op2(yylex, &yyVAL, OP_MUL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 57:
		//line lua.y:146
		{
			op2(yylex, &yyVAL, OP_DIV, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 58:
		//line lua.y:147
		{
			op2(yylex, &yyVAL, OP_PMUL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 59:
		//line lua.y:148
		{
			op2(yylex, &yyVAL, OP_MOD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 60:
		//line lua.y:151
		{
			opFlag(&yyVAL, OP_LT)
		}
	case 61:
		//line lua.y:152
		{
			opFlag(&yyVAL, OP_GT)
		}
	case 62:
		//line lua.y:153
		{
			opFlag(&yyVAL, OP_LTEQ)
		}
	case 63:
		//line lua.y:154
		{
			opFlag(&yyVAL, OP_GTEQ)
		}
	case 64:
		//line lua.y:155
		{
			opFlag(&yyVAL, OP_EQ)
		}
	case 65:
		//line lua.y:156
		{
			opFlag(&yyVAL, OP_NOTEQ)
		}
	case 66:
		//line lua.y:159
		{
			opVar(&yyVAL, &yyS[yypt-0])
		}
	case 67:
		//line lua.y:160
		{
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-3], &yyS[yypt-1])
		}
	case 68:
		//line lua.y:161
		{
			opValueExt(&yyS[yypt-0], yyS[yypt-0].token.image)
			op2(yylex, &yyVAL, OP_MEMBER, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 71:
		//line lua.y:171
		{
			op2(yylex, &yyVAL, OP_CALL, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 73:
		//line lua.y:175
		{
			yyVAL = yyS[yypt-1]
		}
	case 74:
		//line lua.y:178
		{
			yyVAL = yyS[yypt-0]
		}
	case 75:
		//line lua.y:181
		{
			opFunc(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-1])
		}
	case 76:
		//line lua.y:184
		{
			yyVAL = yyS[yypt-1]
		}
	case 78:
		//line lua.y:188
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 79:
		//line lua.y:189
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 80:
		//line lua.y:190
		{
			nameAppend(yylex, &yyVAL, nil, nil)
		}
	case 81:
		//line lua.y:194
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 82:
		//line lua.y:195
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 83:
		//line lua.y:198
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 84:
		//line lua.y:199
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 85:
		//line lua.y:202
		{
			op1(yylex, &yyVAL, OP_ARRAY, nil)
		}
	case 86:
		//line lua.y:203
		{
			op1(yylex, &yyVAL, OP_ARRAY, &yyS[yypt-1])
		}
	case 87:
		//line lua.y:204
		{
			op1(yylex, &yyVAL, OP_ARRAY, &yyS[yypt-2])
		}
	case 88:
		//line lua.y:207
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 89:
		//line lua.y:208
		{
			opExpList(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 90:
		//line lua.y:211
		{
			op1(yylex, &yyVAL, OP_TABLE, nil)
		}
	case 91:
		//line lua.y:212
		{
			op1(yylex, &yyVAL, OP_TABLE, &yyS[yypt-1])
		}
	case 92:
		//line lua.y:213
		{
			op1(yylex, &yyVAL, OP_TABLE, &yyS[yypt-2])
		}
	case 94:
		//line lua.y:217
		{
			opAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 95:
		//line lua.y:220
		{
			opValueExt(&yyS[yypt-2], yyS[yypt-2].token.image)
			op2(yylex, &yyVAL, OP_FIELD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 96:
		//line lua.y:224
		{
			op2(yylex, &yyVAL, OP_FIELD, &yyS[yypt-3], &yyS[yypt-0])
		}
	}
	goto yystack /* stack new state and value */
}
