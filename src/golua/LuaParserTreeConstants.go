package golua

const (
	JJTCHUNK            = 0
	JJTBLOCK            = 1
	JJTSTAT             = 2
	JJTRETURNSTAT       = 3
	JJTLABEL            = 4
	JJTEXPRSTAT         = 5
	JJTASSIGN           = 6
	JJTVAREXP           = 7
	JJTFUNCNAME         = 8
	JJTPREFIXEXP        = 9
	JJTPARENEXP         = 10
	JJTPRIMARYEXP       = 11
	JJTPOSTFIXOP        = 12
	JJTFIELDOP          = 13
	JJTFUNCOP           = 14
	JJTFUNCARGS         = 15
	JJTNAMELIST         = 16
	JJTEXPLIST          = 17
	JJTSIMPLEEXP        = 18
	JJTSTR              = 19
	JJTEXP              = 20
	JJTSUBEXP           = 21
	JJTFUNCTIONCALL     = 22
	JJTFUNCBODY         = 23
	JJTPARLIST          = 24
	JJTTABLECONSTRUCTOR = 25
	JJTFIELDLIST        = 26
	JJTFIELD            = 27
	JJTFIELDSEP         = 28
	JJTBINOP            = 29
	JJTUNOP             = 30

	JJTTOKEN = 31
)

var JJT_NODE_NAME []string

func init() {
	JJT_NODE_NAME = []string{
		"Chunk",
		"Block",
		"Stat",
		"ReturnStat",
		"Label",
		"ExprStat",
		"Assign",
		"VarExp",
		"FuncName",
		"PrefixExp",
		"ParenExp",
		"PrimaryExp",
		"PostfixOp",
		"FieldOp",
		"FuncOp",
		"FuncArgs",
		"NameList",
		"ExpList",
		"SimpleExp",
		"Str",
		"Exp",
		"SubExp",
		"FunctionCall",
		"FuncBody",
		"ParList",
		"TableConstructor",
		"FieldList",
		"Field",
		"FieldSep",
		"Binop",
		"Unop",
		"Token",
	}
}
