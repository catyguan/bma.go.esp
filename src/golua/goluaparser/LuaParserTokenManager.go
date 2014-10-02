package goluaparser

import "bytes"

var lexStateNames []string
var jjnewLexState []int
var jjtoToken []int64
var jjtoSkip []int64
var jjtoSpecial []int64
var jjtoMore []int64

func init() {
	lexStateNames = []string{
		"DEFAULT",
		"IN_COMMENT",
		"IN_LC0",
		"IN_LC1",
		"IN_LC2",
		"IN_LC3",
		"IN_LCN",
		"IN_LS0",
		"IN_LS1",
		"IN_LS2",
		"IN_LS3",
		"IN_LSN",
	}
	/** Lex State array. */
	jjnewLexState = []int{
		-1, -1, -1, -1, -1, -1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 1, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	}
	jjtoToken = []int64{
		0x601fffffef800001, 0x7fffffe2,
	}

	jjtoSkip = []int64{
		0x7e003e, 0x0,
	}

	jjtoSpecial = []int64{
		0x7e0000, 0x0,
	}
	jjtoMore = []int64{
		0x1001ffc0, 0x0,
	}
}

type LuaParserTokenManager struct {
	input_stream *SimpleCharStream
	jjrounds     []uint
	jjstateSet   []int

	jjimage       *bytes.Buffer
	image         *bytes.Buffer
	jjimageLen    int
	lengthOfMatch int

	curChar rune

	curLexState     int
	defaultLexState int
	jjnewStateCnt   int
	jjround         uint
	jjmatchedPos    int
	jjmatchedKind   int
}

func newLuaParserTokenManager() *LuaParserTokenManager {
	r := new(LuaParserTokenManager)
	r.jjrounds = make([]uint, 66)
	r.jjstateSet = make([]int, 2*66)
	r.jjimage = bytes.NewBuffer([]byte{})
	r.image = r.jjimage
	return r
}

func newLuaParserTokenManager1(stream *SimpleCharStream) *LuaParserTokenManager {
	r := newLuaParserTokenManager()
	r.input_stream = stream
	return r
}

func newLuaParserTokenManager2(stream *SimpleCharStream, lexState int) *LuaParserTokenManager {
	r := newLuaParserTokenManager()
	r.ReInit2(stream, lexState)
	return r
}

func (this *LuaParserTokenManager) jjStopAtPos(pos int, kind int) int {
	this.jjmatchedKind = kind
	this.jjmatchedPos = pos
	return pos + 1
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa0_2() int {
	switch this.curChar {
	case 93:
		return this.jjMoveStringLiteralDfa1_2(0x40000)
	default:
		return 1
	}
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa1_2(active0 int64) int {
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 1
	}
	switch this.curChar {
	case 93:
		if (active0 & 0x40000) != 0 {
			return this.jjStopAtPos(1, 18)
		}
		break
	default:
		return 2
	}
	return 2
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa0_11() int {
	return this.jjMoveNfa_11(6, 0)
}
func (this *LuaParserTokenManager) jjMoveNfa_11(startState int, curPos int) int {
	startsAt := 0
	this.jjnewStateCnt = 7
	i := 1
	this.jjstateSet[0] = startState
	kind := 0x7fffffff
	for {
		this.jjround++
		if this.jjround == 0x7fffffff {
			this.ReInitRounds()
		}
		if this.curChar < 64 {
			// long l = 1L << uint(uint(this.curChar))
			for {
				i--
				switch this.jjstateSet[i] {
				case 0, 1:
					if this.curChar == 61 {
						this.jjCheckNAddTwoStates(1, 2)
					}
					break
				case 3:
					if this.curChar == 61 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 0
					}
					break
				case 4:
					if this.curChar == 61 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 3
					}
					break
				case 5:
					if this.curChar == 61 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 4
					}
					break
				default:
					break
				}
				if i == startsAt {
					break
				}
			}
		} else if this.curChar < 128 {
			// long l = 1L << (uint(this.curChar) & 077);
			for {
				i--
				switch this.jjstateSet[i] {
				case 2:
					if this.curChar == 93 && kind > 27 {
						kind = 27
					}
					break
				case 6:
					if this.curChar == 93 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 5
					}
					break
				default:
					break
				}
				if i == startsAt {
					break
				}
			}
		} else {
			// int i2 = (uint(this.curChar) & 0xff) >> 6;
			// long l2 = 1L << (uint(this.curChar) & 077);
			for {
				i--
				switch this.jjstateSet[i] {
				default:
					break
				}
				if i == startsAt {
					break
				}
			}
		}
		if kind != 0x7fffffff {
			this.jjmatchedKind = kind
			this.jjmatchedPos = curPos
			kind = 0x7fffffff
		}
		curPos++
		i = this.jjnewStateCnt
		this.jjnewStateCnt = startsAt
		startsAt = 7 - startsAt
		if i == startsAt {
			return curPos
		}
		this.curChar = this.input_stream.readChar()
		if this.curChar == 0 {
			return curPos
		}
	}
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa0_10() int {
	switch this.curChar {
	case 93:
		return this.jjMoveStringLiteralDfa1_10(0x4000000)
	default:
		return 1
	}
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa1_10(active0 int64) int {
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 1
	}
	switch this.curChar {
	case 61:
		return this.jjMoveStringLiteralDfa2_10(active0, 0x4000000)
	default:
		return 2
	}
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa2_10(old0 int64, active0 int64) int {
	active0 = active0 & old0
	if active0 == 0 {
		return 2
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 2
	}
	switch this.curChar {
	case 61:
		return this.jjMoveStringLiteralDfa3_10(active0, 0x4000000)
	default:
		return 3
	}
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa3_10(old0 int64, active0 int64) int {
	active0 = active0 & old0
	if active0 == 0 {
		return 3
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 3
	}
	switch this.curChar {
	case 61:
		return this.jjMoveStringLiteralDfa4_10(active0, 0x4000000)
	default:
		return 4
	}
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa4_10(old0 int64, active0 int64) int {
	active0 = active0 & old0
	if active0 == 0 {
		return 4
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 4
	}
	switch this.curChar {
	case 93:
		if active0&0x4000000 != 0 {
			return this.jjStopAtPos(4, 26)
		}
		break
	default:
		return 5
	}
	return 5
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa0_9() int {
	switch this.curChar {
	case 93:
		return this.jjMoveStringLiteralDfa1_9(0x2000000)
	default:
		return 1
	}
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa1_9(active0 int64) int {
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 1
	}
	switch this.curChar {
	case 61:
		return this.jjMoveStringLiteralDfa2_9(active0, 0x2000000)
	default:
		return 2
	}
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa2_9(old0 int64, active0 int64) int {
	active0 = active0 & old0
	if active0 == 0 {
		return 2
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 2
	}
	switch this.curChar {
	case 61:
		return this.jjMoveStringLiteralDfa3_9(active0, 0x2000000)
	default:
		return 3
	}
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa3_9(old0 int64, active0 int64) int {
	active0 = active0 & old0
	if active0 == 0 {
		return 3
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 3
	}
	switch this.curChar {
	case 93:
		if active0&0x2000000 != 0 {
			return this.jjStopAtPos(3, 25)
		}
		break
	default:
		return 4
	}
	return 4
}
func (this *LuaParserTokenManager) jjMoveStringLiteralDfa0_8() int {
	switch this.curChar {
	case 93:
		return this.jjMoveStringLiteralDfa1_8(0x1000000)
	default:
		return 1
	}
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa1_8(active0 int64) int {
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 1
	}
	switch this.curChar {
	case 61:
		return this.jjMoveStringLiteralDfa2_8(active0, 0x1000000)
	default:
		return 2
	}
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa2_8(old0 int64, active0 int64) int {
	active0 = active0 & old0
	if active0 == 0 {
		return 2
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 2
	}
	switch this.curChar {
	case 93:
		if (active0 & 0x1000000) != 0 {
			return this.jjStopAtPos(2, 24)
		}
		break
	default:
		return 3
	}
	return 3
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa0_7() int {
	switch this.curChar {
	case 93:
		return this.jjMoveStringLiteralDfa1_7(0x800000)
	default:
		return 1
	}
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa1_7(active0 int64) int {
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 1
	}
	switch this.curChar {
	case 93:
		if (active0 & 0x800000) != 0 {
			return this.jjStopAtPos(1, 23)
		}
		break
	default:
		return 2
	}
	return 2
}

func (this *LuaParserTokenManager) jjStopStringLiteralDfa_0(pos int, active0 int64, active1 int64) int {
	switch pos {
	case 0:
		if (active0&0x7800) != 0 || (active1&0x2000) != 0 {
			return 14
		}
		if (active1 & 0x1008200) != 0 {
			return 31
		}
		if (active0 & 0x7ffffe0000000) != 0 {
			this.jjmatchedKind = 51
			return 17
		}
		if (active0&0x103c0) != 0 || (active1&0x80000) != 0 {
			return 7
		}
		return -1
	case 1:
		if (active0 & 0x103c0) != 0 {
			return 6
		}
		if (active0 & 0x7000) != 0 {
			return 13
		}
		if (active0 & 0x118080000000) != 0 {
			return 17
		}
		if (active0 & 0x7ee7f60000000) != 0 {
			if this.jjmatchedPos != 1 {
				this.jjmatchedKind = 51
				this.jjmatchedPos = 1
			}
			return 17
		}
		return -1
	case 2:
		if (active0 & 0x7e26b40000000) != 0 {
			this.jjmatchedKind = 51
			this.jjmatchedPos = 2
			return 17
		}
		if (active0 & 0x6000) != 0 {
			return 12
		}
		if (active0 & 0x3c0) != 0 {
			return 5
		}
		if (active0 & 0xc1420000000) != 0 {
			return 17
		}
		return -1
	case 3:
		if (active0 & 0x380) != 0 {
			return 4
		}
		if (active0 & 0x6622840000000) != 0 {
			if this.jjmatchedPos != 3 {
				this.jjmatchedKind = 51
				this.jjmatchedPos = 3
			}
			return 17
		}
		if (active0 & 0x1804300000000) != 0 {
			return 17
		}
		if (active0 & 0x4000) != 0 {
			return 9
		}
		return -1
	case 4:
		if (active0 & 0x602200000000) != 0 {
			this.jjmatchedKind = 51
			this.jjmatchedPos = 4
			return 17
		}
		if (active0 & 0x300) != 0 {
			return 3
		}
		if (active0 & 0x6020840000000) != 0 {
			return 17
		}
		return -1
	case 5:
		if (active0 & 0x200) != 0 {
			return 0
		}
		if (active0 & 0x600200000000) != 0 {
			return 17
		}
		if (active0 & 0x2000000000) != 0 {
			this.jjmatchedKind = 51
			this.jjmatchedPos = 5
			return 17
		}
		return -1
	case 6:
		if (active0 & 0x2000000000) != 0 {
			this.jjmatchedKind = 51
			this.jjmatchedPos = 6
			return 17
		}
		return -1
	default:
		return -1
	}
}

func (this *LuaParserTokenManager) jjStartNfa_0(pos int, active0 int64, active1 int64) int {
	return this.jjMoveNfa_0(this.jjStopStringLiteralDfa_0(pos, active0, active1), pos+1)
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa0_0() int {
	switch this.curChar {
	case 35:
		return this.jjStopAtPos(0, 69)
	case 37:
		return this.jjStopAtPos(0, 87)
	case 40:
		return this.jjStopAtPos(0, 75)
	case 41:
		return this.jjStopAtPos(0, 76)
	case 42:
		return this.jjStopAtPos(0, 84)
	case 43:
		return this.jjStopAtPos(0, 82)
	case 44:
		return this.jjStopAtPos(0, 72)
	case 45:
		this.jjmatchedKind = 83
		return this.jjMoveStringLiteralDfa1_0(0x103c0, 0x0)
	case 46:
		this.jjmatchedKind = 73
		return this.jjMoveStringLiteralDfa1_0(0x0, 0x1008000)
	case 47:
		return this.jjStopAtPos(0, 85)
	case 58:
		this.jjmatchedKind = 74
		return this.jjMoveStringLiteralDfa1_0(0x0, 0x2)
	case 59:
		return this.jjStopAtPos(0, 70)
	case 60:
		this.jjmatchedKind = 89
		return this.jjMoveStringLiteralDfa1_0(0x0, 0x4000000)
	case 61:
		this.jjmatchedKind = 71
		return this.jjMoveStringLiteralDfa1_0(0x0, 0x20000000)
	case 62:
		this.jjmatchedKind = 91
		return this.jjMoveStringLiteralDfa1_0(0x0, 0x10000000)
	case 91:
		this.jjmatchedKind = 77
		return this.jjMoveStringLiteralDfa1_0(0x7800, 0x0)
	case 93:
		return this.jjStopAtPos(0, 78)
	case 94:
		return this.jjStopAtPos(0, 86)
	case 97:
		return this.jjMoveStringLiteralDfa1_0(0x20000000, 0x0)
	case 98:
		return this.jjMoveStringLiteralDfa1_0(0x40000000, 0x0)
	case 100:
		return this.jjMoveStringLiteralDfa1_0(0x80000000, 0x0)
	case 101:
		return this.jjMoveStringLiteralDfa1_0(0x700000000, 0x0)
	case 102:
		return this.jjMoveStringLiteralDfa1_0(0x3800000000, 0x0)
	case 103:
		return this.jjMoveStringLiteralDfa1_0(0x4000000000, 0x0)
	case 105:
		return this.jjMoveStringLiteralDfa1_0(0x18000000000, 0x0)
	case 108:
		return this.jjMoveStringLiteralDfa1_0(0x20000000000, 0x0)
	case 110:
		return this.jjMoveStringLiteralDfa1_0(0xc0000000000, 0x0)
	case 111:
		return this.jjMoveStringLiteralDfa1_0(0x100000000000, 0x0)
	case 114:
		return this.jjMoveStringLiteralDfa1_0(0x600000000000, 0x0)
	case 116:
		return this.jjMoveStringLiteralDfa1_0(0x1800000000000, 0x0)
	case 117:
		return this.jjMoveStringLiteralDfa1_0(0x2000000000000, 0x0)
	case 119:
		return this.jjMoveStringLiteralDfa1_0(0x4000000000000, 0x0)
	case 123:
		return this.jjStopAtPos(0, 80)
	case 125:
		return this.jjStopAtPos(0, 81)
	case 126:
		return this.jjMoveStringLiteralDfa1_0(0x0, 0x40000000)
	default:
		return this.jjMoveNfa_0(8, 0)
	}
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa1_0(active0 int64, active1 int64) int {
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		this.jjStopStringLiteralDfa_0(0, active0, active1)
		return 1
	}
	switch this.curChar {
	case 45:
		if (active0 & 0x10000) != 0 {
			this.jjmatchedKind = 16
			this.jjmatchedPos = 1
		}
		return this.jjMoveStringLiteralDfa2_0(active0, 0x3c0, active1, 0)
	case 46:
		if (active1 & 0x1000000) != 0 {
			this.jjmatchedKind = 88
			this.jjmatchedPos = 1
		}
		return this.jjMoveStringLiteralDfa2_0(active0, 0, active1, 0x8000)
	case 58:
		if (active1 & 0x2) != 0 {
			return this.jjStopAtPos(1, 65)
		}
		break
	case 61:
		if (active1 & 0x4000000) != 0 {
			return this.jjStopAtPos(1, 90)
		} else if (active1 & 0x10000000) != 0 {
			return this.jjStopAtPos(1, 92)
		} else if (active1 & 0x20000000) != 0 {
			return this.jjStopAtPos(1, 93)
		} else if (active1 & 0x40000000) != 0 {
			return this.jjStopAtPos(1, 94)
		}
		return this.jjMoveStringLiteralDfa2_0(active0, 0x7000, active1, 0)
	case 91:
		if (active0 & 0x800) != 0 {
			return this.jjStopAtPos(1, 11)
		}
		break
	case 97:
		return this.jjMoveStringLiteralDfa2_0(active0, 0x800000000, active1, 0)
	case 101:
		return this.jjMoveStringLiteralDfa2_0(active0, 0x600000000000, active1, 0)
	case 102:
		if (active0 & 0x8000000000) != 0 {
			return this.jjStartNfaWithStates_0(1, 39, 17)
		}
		break
	case 104:
		return this.jjMoveStringLiteralDfa2_0(active0, 0x4800000000000, active1, 0)
	case 105:
		return this.jjMoveStringLiteralDfa2_0(active0, 0x40000000000, active1, 0)
	case 108:
		return this.jjMoveStringLiteralDfa2_0(active0, 0x300000000, active1, 0)
	case 110:
		if (active0 & 0x10000000000) != 0 {
			return this.jjStartNfaWithStates_0(1, 40, 17)
		}
		return this.jjMoveStringLiteralDfa2_0(active0, 0x2000420000000, active1, 0)
	case 111:
		if (active0 & 0x80000000) != 0 {
			return this.jjStartNfaWithStates_0(1, 31, 17)
		}
		return this.jjMoveStringLiteralDfa2_0(active0, 0xa5000000000, active1, 0)
	case 114:
		if (active0 & 0x100000000000) != 0 {
			return this.jjStartNfaWithStates_0(1, 44, 17)
		}
		return this.jjMoveStringLiteralDfa2_0(active0, 0x1000040000000, active1, 0)
	case 117:
		return this.jjMoveStringLiteralDfa2_0(active0, 0x2000000000, active1, 0)
	default:
		break
	}
	return this.jjStartNfa_0(0, active0, active1)
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa2_0(old0 int64, active0 int64, old1 int64, active1 int64) int {
	active0 &= old0
	active1 &= old1
	if (active0 | active1) == 0 {
		return this.jjStartNfa_0(0, old0, old1)
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		this.jjStopStringLiteralDfa_0(1, active0, active1)
		return 2
	}
	switch this.curChar {
	case 46:
		if (active1 & 0x8000) != 0 {
			return this.jjStopAtPos(2, 79)
		}
		break
	case 61:
		return this.jjMoveStringLiteralDfa3_0(active0, 0x6000, active1, 0)
	case 91:
		if (active0 & 0x1000) != 0 {
			return this.jjStopAtPos(2, 12)
		}
		return this.jjMoveStringLiteralDfa3_0(active0, 0x3c0, active1, 0)
	case 99:
		return this.jjMoveStringLiteralDfa3_0(active0, 0x20000000000, active1, 0)
	case 100:
		if (active0 & 0x20000000) != 0 {
			return this.jjStartNfaWithStates_0(2, 29, 17)
		} else if (active0 & 0x400000000) != 0 {
			return this.jjStartNfaWithStates_0(2, 34, 17)
		}
		break
	case 101:
		return this.jjMoveStringLiteralDfa3_0(active0, 0x800040000000, active1, 0)
	case 105:
		return this.jjMoveStringLiteralDfa3_0(active0, 0x4000000000000, active1, 0)
	case 108:
		if (active0 & 0x40000000000) != 0 {
			return this.jjStartNfaWithStates_0(2, 42, 17)
		}
		return this.jjMoveStringLiteralDfa3_0(active0, 0x800000000, active1, 0)
	case 110:
		return this.jjMoveStringLiteralDfa3_0(active0, 0x2000000000, active1, 0)
	case 112:
		return this.jjMoveStringLiteralDfa3_0(active0, 0x400000000000, active1, 0)
	case 114:
		if (active0 & 0x1000000000) != 0 {
			return this.jjStartNfaWithStates_0(2, 36, 17)
		}
		break
	case 115:
		return this.jjMoveStringLiteralDfa3_0(active0, 0x300000000, active1, 0)
	case 116:
		if (active0 & 0x80000000000) != 0 {
			return this.jjStartNfaWithStates_0(2, 43, 17)
		}
		return this.jjMoveStringLiteralDfa3_0(active0, 0x2204000000000, active1, 0)
	case 117:
		return this.jjMoveStringLiteralDfa3_0(active0, 0x1000000000000, active1, 0)
	default:
		break
	}
	return this.jjStartNfa_0(1, active0, active1)
}
func (this *LuaParserTokenManager) jjMoveStringLiteralDfa3_0(old0 int64, active0 int64, old1 int64, active1 int64) int {
	active0 &= old0
	active1 &= old1
	if (active0 | active1) == 0 {
		return this.jjStartNfa_0(1, old0, old1)
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		this.jjStopStringLiteralDfa_0(2, active0, 0)
		return 3
	}
	switch this.curChar {
	case 61:
		return this.jjMoveStringLiteralDfa4_0(active0, 0x4380)
	case 91:
		if (active0 & 0x40) != 0 {
			return this.jjStopAtPos(3, 6)
		} else if (active0 & 0x2000) != 0 {
			return this.jjStopAtPos(3, 13)
		}
		break
	case 97:
		return this.jjMoveStringLiteralDfa4_0(active0, 0x20040000000)
	case 99:
		return this.jjMoveStringLiteralDfa4_0(active0, 0x2000000000)
	case 101:
		if (active0 & 0x100000000) != 0 {
			this.jjmatchedKind = 32
			this.jjmatchedPos = 3
		} else if (active0 & 0x1000000000000) != 0 {
			return this.jjStartNfaWithStates_0(3, 48, 17)
		}
		return this.jjMoveStringLiteralDfa4_0(active0, 0x400200000000)
	case 105:
		return this.jjMoveStringLiteralDfa4_0(active0, 0x2000000000000)
	case 108:
		return this.jjMoveStringLiteralDfa4_0(active0, 0x4000000000000)
	case 110:
		if (active0 & 0x800000000000) != 0 {
			return this.jjStartNfaWithStates_0(3, 47, 17)
		}
		break
	case 111:
		if (active0 & 0x4000000000) != 0 {
			return this.jjStartNfaWithStates_0(3, 38, 17)
		}
		break
	case 115:
		return this.jjMoveStringLiteralDfa4_0(active0, 0x800000000)
	case 117:
		return this.jjMoveStringLiteralDfa4_0(active0, 0x200000000000)
	default:
		break
	}
	return this.jjStartNfa_0(2, active0, 0)
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa4_0(old0 int64, active0 int64) int {
	active0 &= old0
	if active0 == 0 {
		return this.jjStartNfa_0(2, old0, 0)
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		this.jjStopStringLiteralDfa_0(3, active0, 0)
		return 4
	}
	switch this.curChar {
	case 61:
		return this.jjMoveStringLiteralDfa5_0(active0, 0x300)
	case 91:
		if (active0 & 0x80) != 0 {
			return this.jjStopAtPos(4, 7)
		} else if (active0 & 0x4000) != 0 {
			return this.jjStopAtPos(4, 14)
		}
		break
	case 97:
		return this.jjMoveStringLiteralDfa5_0(active0, 0x400000000000)
	case 101:
		if (active0 & 0x800000000) != 0 {
			return this.jjStartNfaWithStates_0(4, 35, 17)
		} else if (active0 & 0x4000000000000) != 0 {
			return this.jjStartNfaWithStates_0(4, 50, 17)
		}
		break
	case 105:
		return this.jjMoveStringLiteralDfa5_0(active0, 0x200000000)
	case 107:
		if (active0 & 0x40000000) != 0 {
			return this.jjStartNfaWithStates_0(4, 30, 17)
		}
		break
	case 108:
		if (active0 & 0x20000000000) != 0 {
			return this.jjStartNfaWithStates_0(4, 41, 17)
		} else if (active0 & 0x2000000000000) != 0 {
			return this.jjStartNfaWithStates_0(4, 49, 17)
		}
		break
	case 114:
		return this.jjMoveStringLiteralDfa5_0(active0, 0x200000000000)
	case 116:
		return this.jjMoveStringLiteralDfa5_0(active0, 0x2000000000)
	default:
		break
	}
	return this.jjStartNfa_0(3, active0, 0)
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa5_0(old0 int64, active0 int64) int {
	active0 &= old0
	if active0 == 0 {
		return this.jjStartNfa_0(3, old0, 0)
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		this.jjStopStringLiteralDfa_0(4, active0, 0)
		return 5
	}
	switch this.curChar {
	case 61:
		return this.jjMoveStringLiteralDfa6_0(active0, 0x200)
	case 91:
		if (active0 & 0x100) != 0 {
			return this.jjStopAtPos(5, 8)
		}
		break
	case 102:
		if (active0 & 0x200000000) != 0 {
			return this.jjStartNfaWithStates_0(5, 33, 17)
		}
		break
	case 105:
		return this.jjMoveStringLiteralDfa6_0(active0, 0x2000000000)
	case 110:
		if (active0 & 0x200000000000) != 0 {
			return this.jjStartNfaWithStates_0(5, 45, 17)
		}
		break
	case 116:
		if (active0 & 0x400000000000) != 0 {
			return this.jjStartNfaWithStates_0(5, 46, 17)
		}
		break
	default:
		break
	}
	return this.jjStartNfa_0(4, active0, 0)
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa6_0(old0 int64, active0 int64) int {
	active0 &= old0
	if active0 == 0 {
		return this.jjStartNfa_0(4, old0, 0)
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		this.jjStopStringLiteralDfa_0(5, active0, 0)
		return 6
	}
	switch this.curChar {
	case 91:
		if (active0 & 0x200) != 0 {
			return this.jjStopAtPos(6, 9)
		}
		break
	case 111:
		return this.jjMoveStringLiteralDfa7_0(active0, 0x2000000000)
	default:
		break
	}
	return this.jjStartNfa_0(5, active0, 0)
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa7_0(old0 int64, active0 int64) int {
	active0 &= old0
	if active0 == 0 {
		return this.jjStartNfa_0(5, old0, 0)
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		this.jjStopStringLiteralDfa_0(6, active0, 0)
		return 7
	}
	switch this.curChar {
	case 110:
		if (active0 & 0x2000000000) != 0 {
			return this.jjStartNfaWithStates_0(7, 37, 17)
		}
		break
	default:
		break
	}
	return this.jjStartNfa_0(6, active0, 0)
}

func (this *LuaParserTokenManager) jjStartNfaWithStates_0(pos int, kind int, state int) int {
	this.jjmatchedKind = kind
	this.jjmatchedPos = pos
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return pos + 1
	}
	return this.jjMoveNfa_0(state, pos+1)
}

var jjbitVec0 []uint64

func init() {
	jjbitVec0 = []uint64{
		0x0, 0x0, 0xffffffffffffffff, 0xffffffffffffffff,
	}
}

func (this *LuaParserTokenManager) jjMoveNfa_0(startState int, curPos int) int {
	startsAt := 0
	this.jjnewStateCnt = 66
	i := 1
	this.jjstateSet[0] = startState
	kind := 0x7fffffff
	for {
		this.jjround++
		if this.jjround == 0x7fffffff {
			this.ReInitRounds()
		}
		if this.curChar < 64 {
			// long l = 1L << uint(uint(this.curChar)) ;
			l := uint64(1) << uint(uint(this.curChar))
			for {
				i--
				switch this.jjstateSet[i] {
				case 8:
					if (0x3ff000000000000 & l) != 0 {
						if kind > 52 {
							kind = 52
						}
						this.jjCheckNAddStates(0, 3)
					} else if this.curChar == 39 {
						this.jjCheckNAddStates(4, 6)
					} else if this.curChar == 34 {
						this.jjCheckNAddStates(7, 9)
					} else if this.curChar == 46 {
						this.jjCheckNAdd(31)
					} else if this.curChar == 45 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 7
					}
					if this.curChar == 48 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 19
					}
					break
				case 0, 1:
					if this.curChar == 61 {
						this.jjCheckNAddTwoStates(1, 2)
					}
					break
				case 3:
					if this.curChar == 61 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 0
					}
					break
				case 4:
					if this.curChar == 61 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 3
					}
					break
				case 5:
					if this.curChar == 61 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 4
					}
					break
				case 7:
					if this.curChar == 45 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 6
					}
					break
				case 9, 10:
					if this.curChar == 61 {
						this.jjCheckNAddTwoStates(10, 11)
					}
					break
				case 12:
					if this.curChar == 61 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 9
					}
					break
				case 13:
					if this.curChar == 61 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 12
					}
					break
				case 14:
					if this.curChar == 61 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 13
					}
					break
				case 17:
					if (0x3ff000000000000 & l) == 0 {
						break
					}
					if kind > 51 {
						kind = 51
					}
					v := this.jjnewStateCnt
					this.jjnewStateCnt++
					this.jjstateSet[v] = 17
					break
				case 18:
					if this.curChar == 48 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 19
					}
					break
				case 20:
					if this.curChar == 46 {
						this.jjCheckNAdd(21)
					}
					break
				case 21:
					if (0x3ff000000000000 & l) == 0 {
						break
					}
					if kind > 52 {
						kind = 52
					}
					this.jjCheckNAddTwoStates(21, 22)
					break
				case 23:
					if (0x280000000000 & l) != 0 {
						this.jjCheckNAdd(24)
					}
					break
				case 24:
					if (0x3ff000000000000 & l) == 0 {
						break
					}
					if kind > 52 {
						kind = 52
					}
					this.jjCheckNAdd(24)
					break
				case 25:
					if (0x3ff000000000000 & l) == 0 {
						break
					}
					if kind > 52 {
						kind = 52
					}
					this.jjCheckNAddStates(10, 13)
					break
				case 26:
					if (0x3ff000000000000 & l) != 0 {
						this.jjCheckNAddTwoStates(26, 27)
					}
					break
				case 27:
					if this.curChar != 46 {
						break
					}
					if kind > 52 {
						kind = 52
					}
					this.jjCheckNAddTwoStates(28, 22)
					break
				case 28:
					if (0x3ff000000000000 & l) == 0 {
						break
					}
					if kind > 52 {
						kind = 52
					}
					this.jjCheckNAddTwoStates(28, 22)
					break
				case 29:
					if (0x3ff000000000000 & l) == 0 {
						break
					}
					if kind > 52 {
						kind = 52
					}
					this.jjCheckNAddTwoStates(29, 22)
					break
				case 30:
					if this.curChar == 46 {
						this.jjCheckNAdd(31)
					}
					break
				case 31:
					if (0x3ff000000000000 & l) == 0 {
						break
					}
					if kind > 52 {
						kind = 52
					}
					this.jjCheckNAddTwoStates(31, 32)
					break
				case 33:
					if (0x280000000000 & l) != 0 {
						this.jjCheckNAdd(34)
					}
					break
				case 34:
					if (0x3ff000000000000 & l) == 0 {
						break
					}
					if kind > 52 {
						kind = 52
					}
					this.jjCheckNAdd(34)
					break
				case 35:
					if this.curChar == 34 {
						this.jjCheckNAddStates(7, 9)
					}
					break
				case 36:
					if (0xfffffffbffffffff & l) != 0 {
						this.jjCheckNAddStates(7, 9)
					}
					break
				case 37:
					if this.curChar == 34 && kind > 61 {
						kind = 61
					}
					break
				case 39:
					this.jjCheckNAddStates(7, 9)
					break
				case 41:
					if (0x3ff000000000000 & l) != 0 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 42
					}
					break
				case 42:
					if (0x3ff000000000000 & l) != 0 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 43
					}
					break
				case 43:
					if (0x3ff000000000000 & l) != 0 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 44
					}
					break
				case 44, 47:
					if (0x3ff000000000000 & l) != 0 {
						this.jjCheckNAddStates(7, 9)
					}
					break
				case 45:
					if (0x3ff000000000000 & l) != 0 {
						this.jjCheckNAddStates(14, 17)
					}
					break
				case 46:
					if (0x3ff000000000000 & l) != 0 {
						this.jjCheckNAddStates(18, 21)
					}
					break
				case 48:
					if this.curChar == 39 {
						this.jjCheckNAddStates(4, 6)
					}
					break
				case 49:
					if (0xffffff7fffffffff & l) != 0 {
						this.jjCheckNAddStates(4, 6)
					}
					break
				case 50:
					if this.curChar == 39 && kind > 62 {
						kind = 62
					}
					break
				case 52:
					this.jjCheckNAddStates(4, 6)
					break
				case 54:
					if (0x3ff000000000000 & l) != 0 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 55
					}
					break
				case 55:
					if (0x3ff000000000000 & l) != 0 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 56
					}
					break
				case 56:
					if (0x3ff000000000000 & l) != 0 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 57
					}
					break
				case 57, 60:
					if (0x3ff000000000000 & l) != 0 {
						this.jjCheckNAddStates(4, 6)
					}
					break
				case 58:
					if (0x3ff000000000000 & l) != 0 {
						this.jjCheckNAddStates(22, 25)
					}
					break
				case 59:
					if (0x3ff000000000000 & l) != 0 {
						this.jjCheckNAddStates(26, 29)
					}
					break
				case 61:
					if (0x3ff000000000000 & l) == 0 {
						break
					}
					if kind > 52 {
						kind = 52
					}
					this.jjCheckNAddStates(0, 3)
					break
				case 62:
					if (0x3ff000000000000 & l) != 0 {
						this.jjCheckNAddTwoStates(62, 63)
					}
					break
				case 63:
					if this.curChar != 46 {
						break
					}
					if kind > 52 {
						kind = 52
					}
					this.jjCheckNAddTwoStates(64, 32)
					break
				case 64:
					if (0x3ff000000000000 & l) == 0 {
						break
					}
					if kind > 52 {
						kind = 52
					}
					this.jjCheckNAddTwoStates(64, 32)
					break
				case 65:
					if (0x3ff000000000000 & l) == 0 {
						break
					}
					if kind > 52 {
						kind = 52
					}
					this.jjCheckNAddTwoStates(65, 32)
					break
				default:
					break
				}
				if i == startsAt {
					break
				}
			}
		} else if this.curChar < 128 {
			l := uint64(1) << (uint(this.curChar) & 077)
			for {
				i--
				switch this.jjstateSet[i] {
				case 8:
					if (0x7fffffe87fffffe & l) != 0 {
						if kind > 51 {
							kind = 51
						}
						this.jjCheckNAdd(17)
					} else if this.curChar == 91 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 14
					}
					break
				case 2:
					if this.curChar == 91 && kind > 10 {
						kind = 10
					}
					break
				case 6:
					if this.curChar == 91 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 5
					}
					break
				case 11:
					if this.curChar == 91 && kind > 15 {
						kind = 15
					}
					break
				case 15:
					if this.curChar == 91 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 14
					}
					break
				case 16, 17:
					if (0x7fffffe87fffffe & l) == 0 {
						break
					}
					if kind > 51 {
						kind = 51
					}
					this.jjCheckNAdd(17)
					break
				case 19:
					if (0x100000001000000 & l) != 0 {
						this.jjAddStates(30, 31)
					}
					break
				case 21:
					if (0x7e0000007e & l) == 0 {
						break
					}
					if kind > 52 {
						kind = 52
					}
					this.jjCheckNAddTwoStates(21, 22)
					break
				case 22:
					if (0x1002000010020 & l) != 0 {
						this.jjAddStates(32, 33)
					}
					break
				case 25:
					if (0x7e0000007e & l) == 0 {
						break
					}
					if kind > 52 {
						kind = 52
					}
					this.jjCheckNAddStates(10, 13)
					break
				case 26:
					if (0x7e0000007e & l) != 0 {
						this.jjCheckNAddTwoStates(26, 27)
					}
					break
				case 28:
					if (0x7e0000007e & l) == 0 {
						break
					}
					if kind > 52 {
						kind = 52
					}
					this.jjCheckNAddTwoStates(28, 22)
					break
				case 29:
					if (0x7e0000007e & l) == 0 {
						break
					}
					if kind > 52 {
						kind = 52
					}
					this.jjCheckNAddTwoStates(29, 22)
					break
				case 32:
					if (0x2000000020 & l) != 0 {
						this.jjAddStates(34, 35)
					}
					break
				case 36:
					if (0xffffffffefffffff & l) != 0 {
						this.jjCheckNAddStates(7, 9)
					}
					break
				case 38:
					if this.curChar == 92 {
						this.jjAddStates(36, 38)
					}
					break
				case 39:
					this.jjCheckNAddStates(7, 9)
					break
				case 40:
					if this.curChar == 117 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 41
					}
					break
				case 41:
					if (0x7e0000007e & l) != 0 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 42
					}
					break
				case 42:
					if (0x7e0000007e & l) != 0 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 43
					}
					break
				case 43:
					if (0x7e0000007e & l) != 0 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 44
					}
					break
				case 44:
					if (0x7e0000007e & l) != 0 {
						this.jjCheckNAddStates(7, 9)
					}
					break
				case 49:
					if (0xffffffffefffffff & l) != 0 {
						this.jjCheckNAddStates(4, 6)
					}
					break
				case 51:
					if this.curChar == 92 {
						this.jjAddStates(39, 41)
					}
					break
				case 52:
					this.jjCheckNAddStates(4, 6)
					break
				case 53:
					if this.curChar == 117 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 54
					}
					break
				case 54:
					if (0x7e0000007e & l) != 0 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 55
					}
					break
				case 55:
					if (0x7e0000007e & l) != 0 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 56
					}
					break
				case 56:
					if (0x7e0000007e & l) != 0 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 57
					}
					break
				case 57:
					if (0x7e0000007e & l) != 0 {
						this.jjCheckNAddStates(4, 6)
					}
					break
				default:
					break
				}
				if i == startsAt {
					break
				}
			}
		} else {
			i2 := (uint(this.curChar) & 0xff) >> 6
			l2 := uint64(1) << (uint(this.curChar) & 077)
			for {
				i--
				switch this.jjstateSet[i] {
				case 36, 39:
					if (jjbitVec0[i2] & l2) != 0 {
						this.jjCheckNAddStates(7, 9)
					}
					break
				case 49, 52:
					if (jjbitVec0[i2] & l2) != 0 {
						this.jjCheckNAddStates(4, 6)
					}
					break
				default:
					break
				}
				if i == startsAt {
					break
				}
			}
		}
		if kind != 0x7fffffff {
			this.jjmatchedKind = kind
			this.jjmatchedPos = curPos
			kind = 0x7fffffff
		}
		curPos++
		i = this.jjnewStateCnt
		this.jjnewStateCnt = startsAt
		startsAt = 66 - startsAt
		if i == startsAt {
			return curPos
		}
		this.curChar = this.input_stream.readChar()
		if this.curChar == 0 {
			return curPos
		}
	}
}
func (this *LuaParserTokenManager) jjMoveStringLiteralDfa0_1() int {
	return this.jjMoveNfa_1(4, 0)
}
func (this *LuaParserTokenManager) jjMoveNfa_1(startState int, curPos int) int {
	startsAt := 0
	this.jjnewStateCnt = 4
	i := 1
	this.jjstateSet[0] = startState
	kind := 0x7fffffff
	for {
		this.jjround++
		if this.jjround == 0x7fffffff {
			this.ReInitRounds()
		}
		if this.curChar < 64 {
			l := uint64(1) << uint(uint(this.curChar))
			for {
				i--
				switch this.jjstateSet[i] {
				case 4:
					if (0xffffffffffffdbff & l) != 0 {
						if kind > 17 {
							kind = 17
						}
						this.jjCheckNAddStates(42, 44)
					} else if (0x2400 & l) != 0 {
						if kind > 17 {
							kind = 17
						}
					}
					if this.curChar == 13 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 2
					}
					break
				case 0:
					if (0xffffffffffffdbff & l) == 0 {
						break
					}
					kind = 17
					this.jjCheckNAddStates(42, 44)
					break
				case 1:
					if (0x2400&l) != 0 && kind > 17 {
						kind = 17
					}
					break
				case 2:
					if this.curChar == 10 && kind > 17 {
						kind = 17
					}
					break
				case 3:
					if this.curChar == 13 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 2
					}
					break
				default:
					break
				}
				if i == startsAt {
					break
				}
			}
		} else if this.curChar < 128 {
			// l := uint64(1) << uint((this.curChar)&077)
			for {
				i--
				switch this.jjstateSet[i] {
				case 4, 0:
					kind = 17
					this.jjCheckNAddStates(42, 44)
					break
				default:
					break
				}
				if i == startsAt {
					break
				}
			}
		} else {
			i2 := (uint(this.curChar) & 0xff) >> 6
			l2 := uint64(1) << (uint(this.curChar) & 077)
			for {
				i--
				switch this.jjstateSet[i] {
				case 4, 0:
					if (jjbitVec0[i2] & l2) == 0 {
						break
					}
					if kind > 17 {
						kind = 17
					}
					this.jjCheckNAddStates(42, 44)
					break
				default:
					break
				}
				if i == startsAt {
					break
				}
			}
		}
		if kind != 0x7fffffff {
			this.jjmatchedKind = kind
			this.jjmatchedPos = curPos
			kind = 0x7fffffff
		}
		curPos++
		i = this.jjnewStateCnt
		this.jjnewStateCnt = startsAt
		startsAt = 4 - startsAt
		if i == startsAt {
			return curPos
		}
		this.curChar = this.input_stream.readChar()
		if this.curChar == 0 {
			return curPos
		}
	}
}
func (this *LuaParserTokenManager) jjMoveStringLiteralDfa0_6() int {
	return this.jjMoveNfa_6(6, 0)
}
func (this *LuaParserTokenManager) jjMoveNfa_6(startState int, curPos int) int {
	startsAt := 0
	this.jjnewStateCnt = 7
	i := 1
	this.jjstateSet[0] = startState
	kind := 0x7fffffff
	for {
		this.jjround++
		if this.jjround == 0x7fffffff {
			this.ReInitRounds()
		}
		if this.curChar < 64 {
			// long l = 1L << uint(uint(this.curChar));
			for {
				i--
				switch this.jjstateSet[i] {
				case 0, 1:
					if this.curChar == 61 {
						this.jjCheckNAddTwoStates(1, 2)
					}
					break
				case 3:
					if this.curChar == 61 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 0
					}
					break
				case 4:
					if this.curChar == 61 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 3
					}
					break
				case 5:
					if this.curChar == 61 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 4
					}
					break
				default:
					break
				}
				if i == startsAt {
					break
				}
			}
		} else if this.curChar < 128 {
			// long l = 1L << (uint(this.curChar) & 077);
			for {
				i--
				switch this.jjstateSet[i] {
				case 2:
					if this.curChar == 93 && kind > 22 {
						kind = 22
					}
					break
				case 6:
					if this.curChar == 93 {
						v := this.jjnewStateCnt
						this.jjnewStateCnt++
						this.jjstateSet[v] = 5
					}
					break
				default:
					break
				}
				if i == startsAt {
					break
				}
			}
		} else {
			// int i2 = (uint(this.curChar) & 0xff) >> 6;
			// long l2 = 1L << (uint(this.curChar) & 077);
			for {
				i--
				switch this.jjstateSet[i] {
				default:
					break
				}
				if i == startsAt {
					break
				}
			}
		}
		if kind != 0x7fffffff {
			this.jjmatchedKind = kind
			this.jjmatchedPos = curPos
			kind = 0x7fffffff
		}
		curPos++
		i = this.jjnewStateCnt
		this.jjnewStateCnt = startsAt
		startsAt = 7 - startsAt
		if i == startsAt {
			return curPos
		}
		this.curChar = this.input_stream.readChar()
		if this.curChar == 0 {
			return curPos
		}
	}
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa0_5() int {
	switch this.curChar {
	case 93:
		return this.jjMoveStringLiteralDfa1_5(0x200000)
	default:
		return 1
	}
}
func (this *LuaParserTokenManager) jjMoveStringLiteralDfa1_5(active0 int64) int {
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 1
	}
	switch this.curChar {
	case 61:
		return this.jjMoveStringLiteralDfa2_5(active0, 0x200000)
	default:
		return 2
	}
}
func (this *LuaParserTokenManager) jjMoveStringLiteralDfa2_5(old0 int64, active0 int64) int {
	active0 &= old0
	if active0 == 0 {
		return 2
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 2
	}
	switch this.curChar {
	case 61:
		return this.jjMoveStringLiteralDfa3_5(active0, 0x200000)
	default:
		return 3
	}
}
func (this *LuaParserTokenManager) jjMoveStringLiteralDfa3_5(old0 int64, active0 int64) int {
	active0 &= old0
	if active0 == 0 {
		return 3
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 3
	}
	switch this.curChar {
	case 61:
		return this.jjMoveStringLiteralDfa4_5(active0, 0x200000)
	default:
		return 4
	}
}
func (this *LuaParserTokenManager) jjMoveStringLiteralDfa4_5(old0 int64, active0 int64) int {
	active0 &= old0
	if active0 == 0 {
		return 4
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 4
	}
	switch this.curChar {
	case 93:
		if (active0 & 0x200000) != 0 {
			return this.jjStopAtPos(4, 21)
		}
		break
	default:
		return 5
	}
	return 5
}
func (this *LuaParserTokenManager) jjMoveStringLiteralDfa0_4() int {
	switch this.curChar {
	case 93:
		return this.jjMoveStringLiteralDfa1_4(0x100000)
	default:
		return 1
	}
}
func (this *LuaParserTokenManager) jjMoveStringLiteralDfa1_4(active0 int64) int {
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 1
	}
	switch this.curChar {
	case 61:
		return this.jjMoveStringLiteralDfa2_4(active0, 0x100000)
	default:
		return 2
	}
}

func (this *LuaParserTokenManager) jjMoveStringLiteralDfa2_4(old0 int64, active0 int64) int {
	active0 &= old0
	if active0 == 0 {
		return 2
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 2
	}
	switch this.curChar {
	case 61:
		return this.jjMoveStringLiteralDfa3_4(active0, 0x100000)
	default:
		return 3
	}
}
func (this *LuaParserTokenManager) jjMoveStringLiteralDfa3_4(old0 int64, active0 int64) int {
	active0 &= old0
	if active0 == 0 {
		return 3
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 3
	}
	switch this.curChar {
	case 93:
		if (active0 & 0x100000) != 0 {
			return this.jjStopAtPos(3, 20)
		}
		break
	default:
		return 4
	}
	return 4
}
func (this *LuaParserTokenManager) jjMoveStringLiteralDfa0_3() int {
	switch this.curChar {
	case 93:
		return this.jjMoveStringLiteralDfa1_3(0x80000)
	default:
		return 1
	}
}
func (this *LuaParserTokenManager) jjMoveStringLiteralDfa1_3(active0 int64) int {
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 1
	}
	switch this.curChar {
	case 61:
		return this.jjMoveStringLiteralDfa2_3(active0, 0x80000)
	default:
		return 2
	}
}
func (this *LuaParserTokenManager) jjMoveStringLiteralDfa2_3(old0 int64, active0 int64) int {
	active0 &= old0
	if active0 == 0 {
		return 2
	}
	this.curChar = this.input_stream.readChar()
	if this.curChar == 0 {
		return 2
	}
	switch this.curChar {
	case 93:
		if (active0 & 0x80000) != 0 {
			return this.jjStopAtPos(2, 19)
		}
		break
	default:
		return 3
	}
	return 3
}

var jjnextStates []int
var jjstrLiteralImages []string

func init() {
	jjnextStates = []int{
		62, 63, 65, 32, 49, 50, 51, 36, 37, 38, 26, 27, 29, 22, 36, 37,
		38, 46, 36, 47, 37, 38, 49, 50, 51, 59, 49, 60, 50, 51, 20, 25,
		23, 24, 33, 34, 39, 40, 45, 52, 53, 58, 0, 1, 3,
	}

	/** Token literal values. */
	jjstrLiteralImages = []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "\141\156\144", "\142\162\145\141\153", "\144\157", "\145\154\163\145",
		"\145\154\163\145\151\146", "\145\156\144", "\146\141\154\163\145", "\146\157\162",
		"\146\165\156\143\164\151\157\156", "\147\157\164\157", "\151\146", "\151\156", "\154\157\143\141\154",
		"\156\151\154", "\156\157\164", "\157\162", "\162\145\164\165\162\156",
		"\162\145\160\145\141\164", "\164\150\145\156", "\164\162\165\145", "\165\156\164\151\154",
		"\167\150\151\154\145", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "\072\072", "", "", "", "\043", "\073", "\075", "\054", "\056", "\072", "\050",
		"\051", "\133", "\135", "\056\056\056", "\173", "\175", "\053", "\055", "\052", "\057",
		"\136", "\045", "\056\056", "\074", "\074\075", "\076", "\076\075", "\075\075", "\176\075",
	}

}

func (this *LuaParserTokenManager) jjFillToken() *Token {
	var t *Token
	var curTokenImage string
	var beginLine int
	var endLine int
	var beginColumn int
	var endColumn int
	// fmt.Println("1", this.jjmatchedPos, this.image, "end")
	if this.jjmatchedPos < 0 {
		if this.image == nil {
			curTokenImage = ""
		} else {
			curTokenImage = this.image.String()
		}
		beginLine = this.input_stream.getEndLine()
		endLine = beginLine
		beginColumn = this.input_stream.getEndColumn()
		endColumn = beginColumn
	} else {
		im := jjstrLiteralImages[this.jjmatchedKind]
		if im == "" {
			curTokenImage = this.input_stream.GetImage()
		} else {
			curTokenImage = im
		}
		beginLine = this.input_stream.getBeginLine()
		beginColumn = this.input_stream.getBeginColumn()
		endLine = this.input_stream.getEndLine()
		endColumn = this.input_stream.getEndColumn()
	}
	// fmt.Println("2", this.jjmatchedKind, curTokenImage, "end")
	t = newToken2(this.jjmatchedKind, curTokenImage)

	t.BeginLine = beginLine
	t.EndLine = endLine
	t.BeginColumn = beginColumn
	t.EndColumn = endColumn

	return t
}

/** Get the next Token. */
func (this *LuaParserTokenManager) getNextToken() (*Token, error) {
	var specialToken *Token
	var matchedToken *Token
	curPos := 0

EOFLoop:
	for {
		this.curChar = this.input_stream.BeginToken()
		if this.curChar == 0 {
			this.jjmatchedKind = 0
			this.jjmatchedPos = -1
			matchedToken = this.jjFillToken()
			matchedToken.SpecialToken = specialToken
			return matchedToken, nil
		}
		this.image = this.jjimage
		this.image.Reset()
		this.jjimageLen = 0

		for {
			switch this.curLexState {
			case 0:
				this.input_stream.backup(0)
				for this.curChar <= 32 && (0x100003600&(int64(1)<<uint(this.curChar))) != 0 {
					this.curChar = this.input_stream.BeginToken()
				}
				if this.curChar == 0 {
					continue EOFLoop
				}
				this.jjmatchedKind = 0x7fffffff
				this.jjmatchedPos = 0
				curPos = this.jjMoveStringLiteralDfa0_0()
				break
			case 1:
				this.jjmatchedKind = 17
				this.jjmatchedPos = -1
				curPos = 0
				curPos = this.jjMoveStringLiteralDfa0_1()
				break
			case 2:
				this.jjmatchedKind = 0x7fffffff
				this.jjmatchedPos = 0
				curPos = this.jjMoveStringLiteralDfa0_2()
				if this.jjmatchedPos == 0 && this.jjmatchedKind > 28 {
					this.jjmatchedKind = 28
				}
				break
			case 3:
				this.jjmatchedKind = 0x7fffffff
				this.jjmatchedPos = 0
				curPos = this.jjMoveStringLiteralDfa0_3()
				if this.jjmatchedPos == 0 && this.jjmatchedKind > 28 {
					this.jjmatchedKind = 28
				}
				break
			case 4:
				this.jjmatchedKind = 0x7fffffff
				this.jjmatchedPos = 0
				curPos = this.jjMoveStringLiteralDfa0_4()
				if this.jjmatchedPos == 0 && this.jjmatchedKind > 28 {
					this.jjmatchedKind = 28
				}
				break
			case 5:
				this.jjmatchedKind = 0x7fffffff
				this.jjmatchedPos = 0
				curPos = this.jjMoveStringLiteralDfa0_5()
				if this.jjmatchedPos == 0 && this.jjmatchedKind > 28 {
					this.jjmatchedKind = 28
				}
				break
			case 6:
				this.jjmatchedKind = 0x7fffffff
				this.jjmatchedPos = 0
				curPos = this.jjMoveStringLiteralDfa0_6()
				if this.jjmatchedPos == 0 && this.jjmatchedKind > 28 {
					this.jjmatchedKind = 28
				}
				break
			case 7:
				this.jjmatchedKind = 0x7fffffff
				this.jjmatchedPos = 0
				curPos = this.jjMoveStringLiteralDfa0_7()
				if this.jjmatchedPos == 0 && this.jjmatchedKind > 28 {
					this.jjmatchedKind = 28
				}
				break
			case 8:
				this.jjmatchedKind = 0x7fffffff
				this.jjmatchedPos = 0
				curPos = this.jjMoveStringLiteralDfa0_8()
				if this.jjmatchedPos == 0 && this.jjmatchedKind > 28 {
					this.jjmatchedKind = 28
				}
				break
			case 9:
				this.jjmatchedKind = 0x7fffffff
				this.jjmatchedPos = 0
				curPos = this.jjMoveStringLiteralDfa0_9()
				if this.jjmatchedPos == 0 && this.jjmatchedKind > 28 {
					this.jjmatchedKind = 28
				}
				break
			case 10:
				this.jjmatchedKind = 0x7fffffff
				this.jjmatchedPos = 0
				curPos = this.jjMoveStringLiteralDfa0_10()
				if this.jjmatchedPos == 0 && this.jjmatchedKind > 28 {
					this.jjmatchedKind = 28
				}
				break
			case 11:
				this.jjmatchedKind = 0x7fffffff
				this.jjmatchedPos = 0
				curPos = this.jjMoveStringLiteralDfa0_11()
				if this.jjmatchedPos == 0 && this.jjmatchedKind > 28 {
					this.jjmatchedKind = 28
				}
				break
			}
			if this.jjmatchedKind != 0x7fffffff {
				if this.jjmatchedPos+1 < curPos {
					this.input_stream.backup(curPos - this.jjmatchedPos - 1)
				}
				if (jjtoToken[this.jjmatchedKind>>6] & (int64(1) << uint(this.jjmatchedKind&077))) != 0 {
					matchedToken = this.jjFillToken()
					matchedToken.SpecialToken = specialToken
					if jjnewLexState[this.jjmatchedKind] != -1 {
						this.curLexState = jjnewLexState[this.jjmatchedKind]
					}
					return matchedToken, nil
				} else if (jjtoSkip[this.jjmatchedKind>>6] & (int64(1) << uint(this.jjmatchedKind&077))) != 0 {
					if (jjtoSpecial[this.jjmatchedKind>>6] & (int64(1) << uint(this.jjmatchedKind&077))) != 0 {
						matchedToken = this.jjFillToken()
						if specialToken == nil {
							specialToken = matchedToken
						} else {
							matchedToken.SpecialToken = specialToken
							specialToken = matchedToken
							specialToken.Next = matchedToken
						}
						this.SkipLexicalActions(matchedToken)
					} else {
						this.SkipLexicalActions(nil)
					}
					if jjnewLexState[this.jjmatchedKind] != -1 {
						this.curLexState = jjnewLexState[this.jjmatchedKind]
					}
					continue EOFLoop
				}
				this.jjimageLen += this.jjmatchedPos + 1
				if jjnewLexState[this.jjmatchedKind] != -1 {
					this.curLexState = jjnewLexState[this.jjmatchedKind]
				}
				curPos = 0
				this.jjmatchedKind = 0x7fffffff
				this.curChar = this.input_stream.readChar()
				if this.curChar != 0 {
					continue
				}
			}

			error_line := this.input_stream.getEndLine()
			error_column := this.input_stream.getEndColumn()
			error_after := ""
			EOFSeen := false
			c := this.input_stream.readChar()
			if c != 0 {
				this.input_stream.backup(1)
			} else {
				EOFSeen = true
				if curPos > 1 {
					error_after = this.input_stream.GetImage()
				}
				if this.curChar == '\n' || this.curChar == '\r' {
					error_line++
					error_column = 0
				} else {
					error_column++
				}
			}
			if !EOFSeen {
				this.input_stream.backup(1)
				if curPos > 1 {
					error_after = this.input_stream.GetImage()
				}
			}
			return nil, newTokenMgrErrorAll(EOFSeen, this.curLexState, error_line, error_column, error_after, this.curChar, LEXICAL_ERROR)
		}
	}
}

func (this *LuaParserTokenManager) SkipLexicalActions(matchedToken *Token) {
	switch this.jjmatchedKind {
	default:
		break
	}
}

func (this *LuaParserTokenManager) jjCheckNAdd(state int) {
	if this.jjrounds[state] != this.jjround {
		v := this.jjnewStateCnt
		this.jjnewStateCnt++
		this.jjstateSet[v] = state
		this.jjrounds[state] = this.jjround
	}
}

func (this *LuaParserTokenManager) jjAddStates(start int, end int) {
	for {
		v := this.jjnewStateCnt
		this.jjnewStateCnt++
		this.jjstateSet[v] = jjnextStates[start]
		v2 := start
		start++
		if v2 == end {
			break
		}
	}
}

func (this *LuaParserTokenManager) jjCheckNAddTwoStates(state1 int, state2 int) {
	this.jjCheckNAdd(state1)
	this.jjCheckNAdd(state2)
}

func (this *LuaParserTokenManager) jjCheckNAddStates(start int, end int) {
	for {
		this.jjCheckNAdd(jjnextStates[start])
		v2 := start
		start++
		if v2 == end {
			break
		}
	}
}

/** Reinitialise parser. */
func (this *LuaParserTokenManager) ReInit(stream *SimpleCharStream) {
	this.jjmatchedPos = 0
	this.jjnewStateCnt = 0
	this.curLexState = this.defaultLexState
	this.input_stream = stream
	this.ReInitRounds()
}

func (this *LuaParserTokenManager) ReInitRounds() {
	this.jjround = 0x80000001
	for i := 66; i > 0; i-- {
		this.jjrounds[i] = 0x80000000
	}
}

func (this *LuaParserTokenManager) ReInit2(stream *SimpleCharStream, lexState int) {
	this.ReInit(stream)
	this.SwitchTo(lexState)
}

func (this *LuaParserTokenManager) SwitchTo(lexState int) {
	// if (lexState >= 12 || lexState < 0)
	//   throw new TokenMgrError("Error: Ignoring invalid lexical state : " + lexState + ". State unchanged.", TokenMgrError.INVALID_LEXICAL_STATE);
	// else
	this.curLexState = lexState
}
