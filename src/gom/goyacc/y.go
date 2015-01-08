//line gom.y:1
package goyacc

import __yyfmt__ "fmt"

//line gom.y:3
const SERVICE = 57346
const STRUCT = 57347
const TRUE = 57348
const FALSE = 57349
const NIL = 57350
const NUMBER = 57351
const STRING = 57352
const NAME = 57353

var yyToknames = []string{
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

//line gom.y:150

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 67
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 147

var yyAct = []int{

	85, 9, 12, 91, 13, 92, 59, 53, 29, 60,
	98, 54, 68, 77, 19, 97, 46, 25, 26, 67,
	27, 31, 32, 30, 33, 34, 105, 107, 100, 104,
	37, 31, 32, 30, 33, 34, 95, 89, 38, 81,
	37, 66, 55, 56, 61, 23, 86, 50, 38, 48,
	28, 72, 63, 31, 32, 30, 33, 34, 56, 62,
	19, 17, 37, 16, 87, 17, 19, 71, 83, 16,
	38, 75, 55, 56, 47, 80, 61, 82, 84, 93,
	78, 88, 79, 73, 74, 69, 70, 64, 65, 62,
	96, 17, 23, 45, 57, 42, 101, 11, 19, 99,
	93, 103, 102, 16, 47, 17, 106, 22, 51, 23,
	44, 49, 40, 15, 14, 23, 23, 43, 24, 94,
	16, 17, 17, 62, 16, 17, 17, 10, 8, 4,
	90, 76, 58, 18, 41, 52, 39, 21, 20, 36,
	35, 7, 6, 5, 3, 2, 1,
}
var yyPact = []int{

	109, -1000, -1000, 109, -1000, -1000, -1000, -1000, -1000, 109,
	-1000, -1000, 104, -1000, 58, 58, -1000, 58, -1000, -1000,
	-1000, -1000, -1000, 39, 47, 97, 80, 103, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 94, 25, -1000,
	92, -1000, 78, 47, -1000, 71, -1000, 23, -1000, -5,
	-1000, -1000, 69, -1000, -1000, 113, 33, -1000, 67, -1000,
	-1000, 112, -8, -1000, -1000, 64, 47, -1000, 15, -1000,
	52, -1000, 35, -1000, 48, -1000, 19, 108, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 17, -1000, -1000, 35,
	-7, -1000, -1000, 108, 10, 35, -1000, -1000, 108, -1000,
	35, 9, -1000, -1000, -1000, 35, 7, -1000,
}
var yyPgo = []int{

	0, 146, 145, 144, 129, 143, 142, 141, 2, 1,
	4, 8, 128, 140, 139, 127, 136, 135, 7, 11,
	0, 97, 134, 132, 6, 9, 131, 130, 3, 5,
	111, 93, 16,
}
var yyR1 = []int{

	0, 1, 2, 3, 3, 4, 4, 4, 8, 8,
	9, 9, 10, 5, 5, 12, 12, 11, 11, 11,
	11, 11, 11, 11, 6, 6, 15, 16, 16, 16,
	17, 17, 18, 18, 19, 20, 20, 20, 7, 7,
	21, 22, 22, 22, 23, 23, 24, 24, 25, 26,
	27, 27, 27, 28, 28, 29, 14, 14, 14, 30,
	30, 13, 13, 13, 31, 31, 32,
}
var yyR2 = []int{

	0, 1, 1, 1, 2, 1, 1, 1, 1, 3,
	1, 2, 4, 1, 2, 0, 3, 1, 1, 1,
	1, 1, 1, 1, 1, 2, 3, 2, 3, 4,
	1, 3, 1, 2, 3, 1, 4, 6, 1, 2,
	3, 2, 3, 4, 1, 3, 1, 2, 4, 3,
	1, 3, 0, 1, 2, 3, 2, 3, 4, 1,
	3, 2, 3, 4, 1, 3, 3,
}
var yyChk = []int{

	-1000, -1, -2, -3, -4, -5, -6, -7, -12, -9,
	-15, -21, -8, -10, 5, 4, 11, 13, -4, -10,
	-12, -15, -21, 12, 14, -8, -8, -8, 11, -11,
	8, 6, 7, 9, 10, -13, -14, 15, 23, -16,
	15, -22, 15, 14, 16, -31, -32, 10, 24, -30,
	-11, 16, -17, -18, -19, -9, -8, 16, -23, -24,
	-25, -9, 11, -11, 16, 17, 18, 24, 17, 16,
	17, -19, 18, 16, 17, -25, -26, 21, 16, -32,
	-11, 24, -11, 16, -18, -20, 11, 16, -24, 18,
	-27, -28, -29, -9, 11, 19, -20, 22, 17, -29,
	18, -20, -28, -20, 20, 17, -20, 20,
}
var yyDef = []int{

	15, -2, 1, 2, 3, 5, 6, 7, 13, 15,
	24, 38, 0, 10, 0, 0, 8, 0, 4, 11,
	14, 25, 39, 0, 0, 0, 0, 0, 9, 16,
	17, 18, 19, 20, 21, 22, 23, 0, 0, 26,
	0, 40, 0, 0, 61, 0, 64, 0, 56, 0,
	59, 27, 0, 30, 32, 0, 0, 41, 0, 44,
	46, 0, 0, 12, 62, 0, 0, 57, 0, 28,
	0, 33, 0, 42, 0, 47, 0, 52, 63, 65,
	66, 58, 60, 29, 31, 34, 35, 43, 45, 0,
	0, 50, 53, 0, 0, 0, 48, 49, 0, 54,
	0, 0, 51, 55, 36, 0, 0, 37,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	21, 22, 3, 3, 17, 3, 12, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 18, 3,
	19, 14, 20, 3, 13, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 23, 3, 24, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 15, 3, 16,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
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
		//line gom.y:18
		{
			endGOM(yylex, &yyVAL)
		}
	case 9:
		//line gom.y:34
		{
			nameAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 10:
		//line gom.y:37
		{
			annoAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 11:
		//line gom.y:38
		{
			annoAppend(yylex, &yyVAL, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 12:
		//line gom.y:41
		{
			defineAnnotation(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 13:
		//line gom.y:44
		{
			commitValue(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 14:
		//line gom.y:45
		{
			commitValue(yylex, &yyVAL, &yyS[yypt-0], &yyS[yypt-1])
		}
	case 16:
		//line gom.y:48
		{
			op2(yylex, &yyVAL, OP_VALUE, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 17:
		//line gom.y:51
		{
			opValue(yylex, &yyVAL)
		}
	case 18:
		//line gom.y:52
		{
			opValue(yylex, &yyVAL)
		}
	case 19:
		//line gom.y:53
		{
			opValue(yylex, &yyVAL)
		}
	case 20:
		//line gom.y:54
		{
			opValue(yylex, &yyVAL)
		}
	case 21:
		//line gom.y:55
		{
			opValue(yylex, &yyVAL)
		}
	case 24:
		//line gom.y:60
		{
			commitStruct(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 25:
		//line gom.y:61
		{
			commitStruct(yylex, &yyVAL, &yyS[yypt-0], &yyS[yypt-1])
		}
	case 26:
		//line gom.y:64
		{
			op2(yylex, &yyVAL, OP_STRUCT, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 27:
		//line gom.y:67
		{
			opN(yylex, &yyVAL, OP_STRUCT_BODY, nil)
		}
	case 28:
		//line gom.y:68
		{
			opN(yylex, &yyVAL, OP_STRUCT_BODY, &yyS[yypt-1])
		}
	case 29:
		//line gom.y:69
		{
			opN(yylex, &yyVAL, OP_STRUCT_BODY, &yyS[yypt-2])
		}
	case 30:
		//line gom.y:72
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 31:
		//line gom.y:73
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 32:
		//line gom.y:76
		{
			commitStructField(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 33:
		//line gom.y:77
		{
			commitStructField(yylex, &yyVAL, &yyS[yypt-0], &yyS[yypt-1])
		}
	case 34:
		//line gom.y:80
		{
			op2(yylex, &yyVAL, OP_SFIELD, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 35:
		//line gom.y:83
		{
			op2(yylex, &yyVAL, OP_TYPE, &yyS[yypt-0], nil)
		}
	case 36:
		//line gom.y:84
		{
			op2(yylex, &yyVAL, OP_TYPE, &yyS[yypt-3], &yyS[yypt-1])
		}
	case 37:
		//line gom.y:85
		{
			op2(yylex, &yyS[yypt-3], OP_TYPE, &yyS[yypt-3], &yyS[yypt-1])
			op2(yylex, &yyVAL, OP_TYPE, &yyS[yypt-5], &yyS[yypt-3])
		}
	case 38:
		//line gom.y:91
		{
			commitService(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 39:
		//line gom.y:92
		{
			commitService(yylex, &yyVAL, &yyS[yypt-0], &yyS[yypt-1])
		}
	case 40:
		//line gom.y:95
		{
			op2(yylex, &yyVAL, OP_SERVICE, &yyS[yypt-1], &yyS[yypt-0])
		}
	case 41:
		//line gom.y:98
		{
			opN(yylex, &yyVAL, OP_SERVICE_BODY, nil)
		}
	case 42:
		//line gom.y:99
		{
			opN(yylex, &yyVAL, OP_SERVICE_BODY, &yyS[yypt-1])
		}
	case 43:
		//line gom.y:100
		{
			opN(yylex, &yyVAL, OP_SERVICE_BODY, &yyS[yypt-2])
		}
	case 44:
		//line gom.y:103
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 45:
		//line gom.y:104
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 46:
		//line gom.y:107
		{
			commitServiceMethod(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 47:
		//line gom.y:108
		{
			commitServiceMethod(yylex, &yyVAL, &yyS[yypt-0], &yyS[yypt-1])
		}
	case 48:
		//line gom.y:111
		{
			op3(yylex, &yyVAL, OP_SMETHOD, &yyS[yypt-3], &yyS[yypt-2], &yyS[yypt-0])
		}
	case 49:
		//line gom.y:114
		{
			opN(yylex, &yyVAL, OP_SM_PARAMS, &yyS[yypt-1])
		}
	case 50:
		//line gom.y:117
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 51:
		//line gom.y:118
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 52:
		//line gom.y:119
		{
			nodeAppend(yylex, &yyVAL, nil, nil)
		}
	case 53:
		//line gom.y:123
		{
			commitMethodParam(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 54:
		//line gom.y:124
		{
			commitMethodParam(yylex, &yyVAL, &yyS[yypt-0], &yyS[yypt-1])
		}
	case 55:
		//line gom.y:127
		{
			op2(yylex, &yyVAL, OP_SM_PARAM, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 56:
		//line gom.y:130
		{
			defineArray(yylex, &yyVAL, nil)
		}
	case 57:
		//line gom.y:131
		{
			defineArray(yylex, &yyVAL, &yyS[yypt-1])
		}
	case 58:
		//line gom.y:132
		{
			defineArray(yylex, &yyVAL, &yyS[yypt-2])
		}
	case 59:
		//line gom.y:135
		{
			beNode(yylex, &yyVAL, &yyS[yypt-0])
		}
	case 60:
		//line gom.y:136
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 61:
		//line gom.y:139
		{
			defineTable(yylex, &yyVAL, nil)
		}
	case 62:
		//line gom.y:140
		{
			defineTable(yylex, &yyVAL, &yyS[yypt-1])
		}
	case 63:
		//line gom.y:141
		{
			defineTable(yylex, &yyVAL, &yyS[yypt-2])
		}
	case 64:
		//line gom.y:144
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-0], nil)
		}
	case 65:
		//line gom.y:145
		{
			nodeAppend(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	case 66:
		//line gom.y:148
		{
			defineField(yylex, &yyVAL, &yyS[yypt-2], &yyS[yypt-0])
		}
	}
	goto yystack /* stack new state and value */
}
