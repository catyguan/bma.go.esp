//line gom.y:1
package goyacc

import __yyfmt__ "fmt"

//line gom.y:3
const OBJECT = 57346
const SERVICE = 57347
const STRUCT = 57348
const TRUE = 57349
const FALSE = 57350
const NIL = 57351
const NUMBER = 57352
const STRING = 57353
const NAME = 57354

var yyToknames = []string{
	"OBJECT",
	"SERVICE",
	"STRUCT",
	"TRUE",
	"FALSE",
	"NIL",
	"NUMBER",
	"STRING",
	"NAME",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line gom.y:172

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 78
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 175

var yyAct = []int{

	101, 10, 15, 107, 72, 108, 74, 73, 53, 61,
	91, 67, 125, 22, 113, 82, 34, 36, 37, 35,
	38, 39, 81, 118, 105, 80, 42, 36, 37, 35,
	38, 39, 92, 93, 43, 97, 42, 36, 37, 35,
	38, 39, 52, 116, 43, 55, 42, 62, 115, 68,
	123, 75, 60, 122, 43, 66, 56, 27, 63, 14,
	57, 27, 14, 86, 49, 22, 102, 77, 76, 14,
	20, 22, 85, 111, 33, 29, 30, 31, 22, 32,
	89, 87, 88, 83, 84, 85, 62, 89, 95, 27,
	68, 100, 47, 109, 19, 75, 104, 96, 112, 98,
	69, 54, 20, 78, 79, 103, 114, 94, 19, 71,
	20, 54, 22, 99, 119, 117, 48, 51, 109, 121,
	120, 76, 13, 20, 124, 69, 70, 20, 12, 19,
	64, 20, 106, 26, 58, 18, 17, 16, 27, 25,
	27, 45, 50, 19, 11, 20, 27, 110, 28, 20,
	76, 69, 20, 20, 19, 24, 20, 9, 4, 90,
	65, 46, 21, 59, 44, 41, 40, 8, 23, 7,
	6, 5, 3, 2, 1,
}
var yyPact = []int{

	131, -1000, -1000, 131, -1000, -1000, -1000, -1000, -1000, -1000,
	131, -1000, -1000, -1000, 133, -1000, 82, 82, 82, -1000,
	82, -1000, -1000, -1000, -1000, -1000, -1000, 62, 30, 125,
	76, 48, 127, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 100, 20, -1000, 117, -1000, 113, -1000, 109,
	30, -1000, 86, -1000, 6, -1000, -3, -1000, -1000, 66,
	-1000, -1000, 142, 44, -1000, 64, -1000, -1000, 139, -12,
	-1000, 15, -1000, -1000, -1000, 138, -12, -1000, -1000, 90,
	30, -1000, 10, -1000, 96, -1000, 54, -1000, 88, -1000,
	5, 135, -1000, 56, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -6, -1000, -1000, 54, 25, -1000, -1000, 135,
	4, -1000, -1000, 54, -1000, -1000, 135, -1000, 54, 32,
	-1000, -1000, -1000, 54, -9, -1000,
}
var yyPgo = []int{

	0, 174, 173, 172, 158, 171, 170, 169, 167, 58,
	1, 2, 16, 157, 166, 165, 144, 164, 163, 6,
	9, 0, 128, 161, 160, 7, 11, 159, 132, 3,
	5, 122, 116, 109, 4, 56, 42, 8,
}
var yyR1 = []int{

	0, 1, 2, 3, 3, 4, 4, 4, 4, 9,
	9, 10, 10, 11, 5, 5, 13, 13, 12, 12,
	12, 12, 12, 12, 12, 6, 6, 16, 17, 17,
	17, 18, 18, 19, 19, 20, 21, 21, 21, 7,
	7, 22, 23, 23, 23, 24, 24, 25, 25, 26,
	27, 28, 28, 28, 29, 29, 30, 8, 8, 31,
	32, 32, 32, 33, 33, 34, 34, 15, 15, 15,
	35, 35, 14, 14, 14, 36, 36, 37,
}
var yyR2 = []int{

	0, 1, 1, 1, 2, 1, 1, 1, 1, 1,
	3, 1, 2, 4, 1, 2, 0, 3, 1, 1,
	1, 1, 1, 1, 1, 1, 2, 3, 2, 3,
	4, 1, 3, 1, 2, 3, 1, 4, 6, 1,
	2, 3, 2, 3, 4, 1, 3, 1, 2, 4,
	3, 1, 3, 0, 1, 2, 3, 1, 2, 3,
	2, 3, 4, 1, 3, 1, 1, 2, 3, 4,
	1, 3, 2, 3, 4, 1, 3, 3,
}
var yyChk = []int{

	-1000, -1, -2, -3, -4, -5, -6, -7, -8, -13,
	-10, -16, -22, -31, -9, -11, 6, 5, 4, 12,
	14, -4, -11, -13, -16, -22, -31, 13, 15, -9,
	-9, -9, -9, 12, -12, 9, 7, 8, 10, 11,
	-14, -15, 16, 24, -17, 16, -23, 16, -32, 16,
	15, 17, -36, -37, 11, 25, -35, -12, 17, -18,
	-19, -20, -10, -9, 17, -24, -25, -26, -10, 12,
	17, -33, -34, -25, -19, -10, 12, -12, 17, 18,
	19, 25, 18, 17, 18, -20, 19, 17, 18, -26,
	-27, 22, 17, 18, 17, -37, -12, 25, -12, 17,
	-19, -21, 12, 17, -25, 19, -28, -29, -30, -10,
	12, 17, -34, 20, -21, 23, 18, -30, 19, -21,
	-29, -21, 21, 18, -21, 21,
}
var yyDef = []int{

	16, -2, 1, 2, 3, 5, 6, 7, 8, 14,
	16, 25, 39, 57, 0, 11, 0, 0, 0, 9,
	0, 4, 12, 15, 26, 40, 58, 0, 0, 0,
	0, 0, 0, 10, 17, 18, 19, 20, 21, 22,
	23, 24, 0, 0, 27, 0, 41, 0, 59, 0,
	0, 72, 0, 75, 0, 67, 0, 70, 28, 0,
	31, 33, 0, 0, 42, 0, 45, 47, 0, 0,
	60, 0, 63, 65, 66, 0, 9, 13, 73, 0,
	0, 68, 0, 29, 0, 34, 0, 43, 0, 48,
	0, 53, 61, 0, 74, 76, 77, 69, 71, 30,
	32, 35, 36, 44, 46, 0, 0, 51, 54, 0,
	0, 62, 64, 0, 49, 50, 0, 55, 0, 0,
	52, 56, 37, 0, 0, 38,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	22, 23, 3, 3, 18, 3, 13, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 19, 3,
	20, 15, 21, 3, 14, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 24, 3, 25, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 16, 3, 17,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12,
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
		//line gom.y:19
		{
			endGOM(yylex, &yyVAL)
		}
	case 10:
		//line gom.y:36
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 11:
		//line gom.y:39
		{
			annoAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 12:
		//line gom.y:40
		{
			annoAppend(yylex, &yyVAL, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 13:
		//line gom.y:43
		{
			defineAnnotation(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 14:
		//line gom.y:46
		{
			commitValue(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 15:
		//line gom.y:47
		{
			commitValue(yylex, &yyVAL, &yyS[yypt-0], &yyS[yypt-1])
		}
	case 17:
		//line gom.y:50
		{
			op2(yylex, &yyVAL, OP_VALUE, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 18:
		//line gom.y:53
		{
			opValue(yylex, &yyVAL)
		}
	case 19:
		//line gom.y:54
		{
			opValue(yylex, &yyVAL)
		}
	case 20:
		//line gom.y:55
		{
			opValue(yylex, &yyVAL)
		}
	case 21:
		//line gom.y:56
		{
			opValue(yylex, &yyVAL)
		}
	case 22:
		//line gom.y:57
		{
			opValue(yylex, &yyVAL)
		}
	case 25:
		//line gom.y:62
		{
			commitStruct(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 26:
		//line gom.y:63
		{
			commitStruct(yylex, &yyVAL, &yyS[yypt-0], &yyS[yypt-1])
		}
	case 27:
		//line gom.y:66
		{
			op2(yylex, &yyVAL, OP_STRUCT, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 28:
		//line gom.y:69
		{
			opN(yylex, &yyVAL, OP_STRUCT_BODY, nil)
		}
	case 29:
		//line gom.y:70
		{
			opN(yylex, &yyVAL, OP_STRUCT_BODY, &yyS[yypt-1])
		}
	case 30:
		//line gom.y:71
		{
			opN(yylex, &yyVAL, OP_STRUCT_BODY, &yyS[yypt-2])
		}
	case 31:
		//line gom.y:74
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 32:
		//line gom.y:75
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 33:
		//line gom.y:78
		{
			commitStructField(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 34:
		//line gom.y:79
		{
			commitStructField(yylex, &yyVAL, &yyS[yypt-0], &yyS[yypt-1])
		}
	case 35:
		//line gom.y:82
		{
			op2(yylex, &yyVAL, OP_SFIELD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 36:
		//line gom.y:85
		{
			op2(yylex, &yyVAL, OP_TYPE, &yyS[yypt-0], nil)
		}
	case 37:
		//line gom.y:86
		{
			op2(yylex, &yyVAL, OP_TYPE, &yyS[yypt-3], &yyS[yypt-1])
		}
	case 38:
		//line gom.y:87
		{
			op2(yylex, &yyS[yypt-3], OP_TYPE, &yyS[yypt-3], &yyS[yypt-1])
			op2(yylex, &yyVAL, OP_TYPE, &yyS[yypt-5], &yyS[yypt-3])
		}
	case 39:
		//line gom.y:93
		{
			commitService(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 40:
		//line gom.y:94
		{
			commitService(yylex, &yyVAL, &yyS[yypt-0], &yyS[yypt-1])
		}
	case 41:
		//line gom.y:97
		{
			op2(yylex, &yyVAL, OP_SERVICE, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 42:
		//line gom.y:100
		{
			opN(yylex, &yyVAL, OP_SERVICE_BODY, nil)
		}
	case 43:
		//line gom.y:101
		{
			opN(yylex, &yyVAL, OP_SERVICE_BODY, &yyS[yypt-1])
		}
	case 44:
		//line gom.y:102
		{
			opN(yylex, &yyVAL, OP_SERVICE_BODY, &yyS[yypt-2])
		}
	case 45:
		//line gom.y:105
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 46:
		//line gom.y:106
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 47:
		//line gom.y:109
		{
			commitServiceMethod(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 48:
		//line gom.y:110
		{
			commitServiceMethod(yylex, &yyVAL, &yyS[yypt-0], &yyS[yypt-1])
		}
	case 49:
		//line gom.y:113
		{
			op3(yylex, &yyVAL, OP_SMETHOD, &yyS[yypt-3], &yyS[yypt-2], &yyS[yypt-0])
		}
	case 50:
		//line gom.y:116
		{
			opN(yylex, &yyVAL, OP_SM_PARAMS, &yyS[yypt-1])
		}
	case 51:
		//line gom.y:119
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 52:
		//line gom.y:120
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 53:
		//line gom.y:121
		{
			nodeAppend(yylex, &yyVAL, nil, nil)
		}
	case 54:
		//line gom.y:125
		{
			commitMethodParam(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 55:
		//line gom.y:126
		{
			commitMethodParam(yylex, &yyVAL, &yyS[yypt-0], &yyS[yypt-1])
		}
	case 56:
		//line gom.y:129
		{
			op2(yylex, &yyVAL, OP_SM_PARAM, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 57:
		//line gom.y:132
		{
			commitObject(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 58:
		//line gom.y:133
		{
			commitObject(yylex, &yyVAL, &yyS[yypt-0], &yyS[yypt-1])
		}
	case 59:
		//line gom.y:136
		{
			op2(yylex, &yyVAL, OP_OBJECT, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 60:
		//line gom.y:139
		{
			opN(yylex, &yyVAL, OP_OBJECT_BODY, nil)
		}
	case 61:
		//line gom.y:140
		{
			opN(yylex, &yyVAL, OP_SERVICE_BODY, &yyS[yypt-1])
		}
	case 62:
		//line gom.y:141
		{
			opN(yylex, &yyVAL, OP_OBJECT_BODY, &yyS[yypt-2])
		}
	case 63:
		//line gom.y:144
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 64:
		//line gom.y:145
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 67:
		//line gom.y:152
		{
			defineArray(yylex, &yyVAL, nil)
		}
	case 68:
		//line gom.y:153
		{
			defineArray(yylex, &yyVAL, &yyS[yypt-1])
		}
	case 69:
		//line gom.y:154
		{
			defineArray(yylex, &yyVAL, &yyS[yypt-2])
		}
	case 70:
		//line gom.y:157
		{
			beNode(yylex, &yyVAL, &yyS[yypt-0])
		}
	case 71:
		//line gom.y:158
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 72:
		//line gom.y:161
		{
			defineTable(yylex, &yyVAL, nil)
		}
	case 73:
		//line gom.y:162
		{
			defineTable(yylex, &yyVAL, &yyS[yypt-1])
		}
	case 74:
		//line gom.y:163
		{
			defineTable(yylex, &yyVAL, &yyS[yypt-2])
		}
	case 75:
		//line gom.y:166
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 76:
		//line gom.y:167
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 77:
		//line gom.y:170
		{
			defineField(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	}
	goto yystack /* stack new state and value */
}
