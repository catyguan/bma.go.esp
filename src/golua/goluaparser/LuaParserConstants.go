package goluaparser

const (
	/** End of File. */
	EOF = 0
	/** RegularExpression Id. */
	COMMENT = 17
	/** RegularExpression Id. */
	LONGCOMMENT0 = 18
	/** RegularExpression Id. */
	LONGCOMMENT1 = 19
	/** RegularExpression Id. */
	LONGCOMMENT2 = 20
	/** RegularExpression Id. */
	LONGCOMMENT3 = 21
	/** RegularExpression Id. */
	LONGCOMMENTN = 22
	/** RegularExpression Id. */
	LONGSTRING0 = 23
	/** RegularExpression Id. */
	LONGSTRING1 = 24
	/** RegularExpression Id. */
	LONGSTRING2 = 25
	/** RegularExpression Id. */
	LONGSTRING3 = 26
	/** RegularExpression Id. */
	LONGSTRINGN = 27
	/** RegularExpression Id. */
	AND = 29
	/** RegularExpression Id. */
	BREAK = 30
	/** RegularExpression Id. */
	DO = 31
	/** RegularExpression Id. */
	ELSE = 32
	/** RegularExpression Id. */
	ELSEIF = 33
	/** RegularExpression Id. */
	END = 34
	/** RegularExpression Id. */
	FALSE = 35
	/** RegularExpression Id. */
	FOR = 36
	/** RegularExpression Id. */
	FUNCTION = 37
	/** RegularExpression Id. */
	GOTO = 38
	/** RegularExpression Id. */
	IF = 39
	/** RegularExpression Id. */
	IN = 40
	/** RegularExpression Id. */
	LOCAL = 41
	/** RegularExpression Id. */
	NIL = 42
	/** RegularExpression Id. */
	NOT = 43
	/** RegularExpression Id. */
	OR = 44
	/** RegularExpression Id. */
	RETURN = 45
	/** RegularExpression Id. */
	REPEAT = 46
	/** RegularExpression Id. */
	THEN = 47
	/** RegularExpression Id. */
	TRUE = 48
	/** RegularExpression Id. */
	UNTIL = 49
	/** RegularExpression Id. */
	WHILE = 50
	/** RegularExpression Id. */
	NAME = 51
	/** RegularExpression Id. */
	NUMBER = 52
	/** RegularExpression Id. */
	FLOAT = 53
	/** RegularExpression Id. */
	FNUM = 54
	/** RegularExpression Id. */
	DIGIT = 55
	/** RegularExpression Id. */
	EXP = 56
	/** RegularExpression Id. */
	HEX = 57
	/** RegularExpression Id. */
	HEXNUM = 58
	/** RegularExpression Id. */
	HEXDIGIT = 59
	/** RegularExpression Id. */
	HEXEXP = 60
	/** RegularExpression Id. */
	STRING = 61
	/** RegularExpression Id. */
	CHARSTRING = 62
	/** RegularExpression Id. */
	QUOTED = 63
	/** RegularExpression Id. */
	DECIMAL = 64
	/** RegularExpression Id. */
	DBCOLON = 65
	/** RegularExpression Id. */
	UNICODE = 66
	/** RegularExpression Id. */
	CHAR = 67
	/** RegularExpression Id. */
	LF = 68

	/** Lexical state. */
	DEFAULT = 0
	/** Lexical state. */
	IN_COMMENT = 1
	/** Lexical state. */
	IN_LC0 = 2
	/** Lexical state. */
	IN_LC1 = 3
	/** Lexical state. */
	IN_LC2 = 4
	/** Lexical state. */
	IN_LC3 = 5
	/** Lexical state. */
	IN_LCN = 6
	/** Lexical state. */
	IN_LS0 = 7
	/** Lexical state. */
	IN_LS1 = 8
	/** Lexical state. */
	IN_LS2 = 9
	/** Lexical state. */
	IN_LS3 = 10
	/** Lexical state. */
	IN_LSN = 11
)

var tokenImage []string
var TokenImage []string

func init() {
	tokenImage = []string{
		"<EOF>",
		"\" \"",
		"\"\\t\"",
		"\"\\n\"",
		"\"\\r\"",
		"\"\\f\"",
		"\"--[[\"",
		"\"--[=[\"",
		"\"--[==[\"",
		"\"--[===[\"",
		"<token of kind 10>",
		"\"[[\"",
		"\"[=[\"",
		"\"[==[\"",
		"\"[===[\"",
		"<token of kind 15>",
		"\"--\"",
		"<COMMENT>",
		"\"]]\"",
		"\"]=]\"",
		"\"]==]\"",
		"\"]===]\"",
		"<LONGCOMMENTN>",
		"\"]]\"",
		"\"]=]\"",
		"\"]==]\"",
		"\"]===]\"",
		"<LONGSTRINGN>",
		"<token of kind 28>",
		"\"and\"",
		"\"break\"",
		"\"do\"",
		"\"else\"",
		"\"elseif\"",
		"\"end\"",
		"\"false\"",
		"\"for\"",
		"\"function\"",
		"\"goto\"",
		"\"if\"",
		"\"in\"",
		"\"local\"",
		"\"nil\"",
		"\"not\"",
		"\"or\"",
		"\"return\"",
		"\"repeat\"",
		"\"then\"",
		"\"true\"",
		"\"until\"",
		"\"while\"",
		"<NAME>",
		"<NUMBER>",
		"<FLOAT>",
		"<FNUM>",
		"<DIGIT>",
		"<EXP>",
		"<HEX>",
		"<HEXNUM>",
		"<HEXDIGIT>",
		"<HEXEXP>",
		"<STRING>",
		"<CHARSTRING>",
		"<QUOTED>",
		"<DECIMAL>",
		"\"::\"",
		"<UNICODE>",
		"<CHAR>",
		"<LF>",
		"\"#\"",
		"\";\"",
		"\"=\"",
		"\",\"",
		"\".\"",
		"\":\"",
		"\"(\"",
		"\")\"",
		"\"[\"",
		"\"]\"",
		"\"...\"",
		"\"{\"",
		"\"}\"",
		"\"+\"",
		"\"-\"",
		"\"*\"",
		"\"/\"",
		"\"^\"",
		"\"%\"",
		"\"..\"",
		"\"<\"",
		"\"<=\"",
		"\">\"",
		"\">=\"",
		"\"==\"",
		"\"~=\"",
	}
	TokenImage = tokenImage
}
