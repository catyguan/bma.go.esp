%{
package goyacc

%}

%token AND
%token BREAK
%token DO
%token ELSE
%token ELSEIF
%token END
%token FALSE
%token FOR
%token FUNCTION
%token GOTO
%token IF
%token IN
%token LOCAL
%token NIL
%token NOT
%token OR
%token RETURN
%token REPEAT
%token THEN
%token TRUE
%token UNTIL
%token WHILE

%token NUMBER
%token STRING
%token SLTEQ
%token SGTEQ
%token SEQ
%token SNOTEQ
%token NAME
%token MORE
%token STRADD

%left  OR
%left  AND
%left  '<' SLTEQ '>' SGTEQ SEQ SNOTEQ
%right STRADD
%left  '+' '-'
%left  '*' '/' '%'
%right NOT '~'
%right '^'

%%
Chunk:
	Block { endChunk(yylex, &$$) }

Block:
	StatList
	| StatList LastStat { opAppend(yylex, &$$, &$1, &$2) }

UBlock:
	Block UNTIL Exp { op2(yylex, &$$, OP_UNTIL, &$1, &$3) }

StatList:
	Stat
	| StatList Stat { opAppend(yylex, &$$, &$1, &$2) }
	|

Stat:
	Binding
	| DO Block END { $$ = $2 }
	| WHILE Exp DO Block END { op2(yylex, &$$, OP_WHILE, &$2, &$4) }
	| Repetition DO Block END
	| REPEAT UBlock { $$ = $2 }
	| IF Conds END
	| FUNCTION FuncName FuncBody
	| SetList '=' ExpList1 { op2(yylex, &$$, OP_ASSIGN, &$1, &$3) }
	| FuncCall

Repetition:
	FOR NAME '=' ExpList23
	| FOR NameList IN ExpList1

Conds:
 	CondList
 	| CondList ELSE Block

CondList:
	Cond
	| CondList ELSEIF Cond

Cond:
	Exp THEN Block

LastStat:
	BREAK
	| RETURN { op1(yylex, &$$, OP_RETURN, nil) }
	| RETURN ExpList1 { op1(yylex, &$$,OP_RETURN, &$2) }

Binding:
    LOCAL NameList { opLocal(yylex, &$$, &$2, nil) }
    | LOCAL NameList '=' ExpList1 { opLocal(yylex, &$$, &$2, &$4) }
	| LOCAL FUNCTION NAME FuncBody

SetList:
	Var
	| SetList ',' Var

FuncName:
	DottedName
	| DottedName ':' NAME

DottedName:
	NAME
	| DottedName ',' NAME

Exp:
	NIL { opValue(yylex, &$$) }
	| TRUE { opValue(yylex, &$$) }
	| FALSE { opValue(yylex, &$$) }
	| NUMBER { opValue(yylex, &$$) }
	| STRING { opValue(yylex, &$$) }
	| MORE
	| FuncDef
	| PrefixExp
	| Tableconstructor
	| Arrayconstructor
	| NOT Exp { op1(yylex, &$$, OP_NOT, &$2) }
	| '#' Exp { op1(yylex, &$$, OP_LEN, &$2) }
	| '-' Exp %prec'*' { op1(yylex, &$$, OP_NSIGN, &$2) }
	| Exp OR Exp { op2(yylex, &$$, OP_OR, &$1, &$3) }
	| Exp AND Exp { op2(yylex, &$$, OP_AND, &$1, &$3) }
	| Exp LogicOp Exp { op2(yylex, &$$, $2.op, &$1, &$3) }
	| Exp STRADD Exp { op2(yylex, &$$, OP_STRADD, &$1, &$3) }
	| Exp '-' Exp { op2(yylex, &$$, OP_SUB, &$1, &$3) }
	| Exp '+' Exp { op2(yylex, &$$, OP_ADD, &$1, &$3) }
	| Exp '*' Exp { op2(yylex, &$$, OP_MUL, &$1, &$3) }
	| Exp '/' Exp { op2(yylex, &$$, OP_DIV, &$1, &$3) }
	| Exp '^' Exp { op2(yylex, &$$, OP_PMUL, &$1, &$3) }
	| Exp '%' Exp { op2(yylex, &$$, OP_MOD, &$1, &$3) }

LogicOp:
	'<' { opFlag(&$$, OP_LT) }
	| '>' { opFlag(&$$, OP_GT) }
	| SLTEQ { opFlag(&$$, OP_LTEQ) }
	| SGTEQ { opFlag(&$$, OP_GTEQ) }
	| SEQ { opFlag(&$$, OP_EQ) }
	| SNOTEQ { opFlag(&$$, OP_NOTEQ) }

Var:
	NAME { opVar(&$$, &$1) }
	| PrefixExp '[' Exp ']' { op2(yylex, &$$, OP_MEMBER, &$1, &$3) }
	| PrefixExp '.' NAME { 
		opValueExt(&$3, $3.token.image)
		op2(yylex, &$$, OP_MEMBER, &$1, &$3)
	}

PrefixExp:
	Var
	| FuncCall

FuncCall:
	PrefixExp Args
	| PrefixExp ':' NAME Args

Args:
	'(' ')'
	| '(' ExpList1 ')'

FuncDef:
	FUNCTION FuncBody

FuncBody:
	ParamDefList Block END

ParamDefList:
	'(' ParDefList ')'

ParDefList:
	NameList
	| MORE
	| NameList "," MORE
	|
	;

NameList:
	NAME { nameAppend(yylex, &$$, &$1, nil) }
	| NameList "," NAME { nameAppend(yylex, &$$, &$1, &$3) }

Arrayconstructor:
	'[' ']'
	| '[' ExpList1 ']'

ExpList1:
	Exp
	| ExpList1 ',' Exp

ExpList23:
	Exp ',' Exp
	| Exp ',' Exp ',' Exp

Tableconstructor:
	'{' '}'
	| '{' FieldList '}'

FieldList:
	Field
	| FieldList FieldSP Field

FieldSP:
	','
	| ';'
            
Field:
	Exp
	| NAME '=' Exp
	| '[' Exp ']' '=' Exp

%%