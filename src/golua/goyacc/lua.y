%{
package goyacc

%}

%token AND
%token BREAK
%token DO
%token ELSEIF
%token ELSE
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
%token CLOSURE
%token CONTINUE

%token NUMBER
%token STRING
%token SLTEQ
%token SGTEQ
%token SEQ
%token SNOTEQ
%token NAME
%token MORE
%token STRADD

%token ANNOTATION

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
	| { $$.value = nil }

Stat:
	Binding
	| ANNOTATION { defineAnnotation(yylex, &$$) }
	| DO Block END { $$ = $2 }
	| WHILE Exp DO Block END { op2(yylex, &$$, OP_WHILE, &$2, &$4) }
	| Repetition DO Block END { opForBind(yylex, &$$, &$1, &$3) }
	| REPEAT UBlock { $$ = $2 }
	| IF Conds END { $$ = $2 }
	| FUNCTION FuncName FuncBody {
		bindFuncName(yylex, &$3, &$2, "")
		op2(yylex, &$$, OP_ASSIGN, &$2, &$3) 
	}
	| SetList '=' ExpList { op2(yylex, &$$, OP_ASSIGN, &$1, &$3) }
	| FuncCall

Repetition:
	FOR NAME '=' ExpList { opFor(yylex, &$$, OP_FOR, &$2, &$4)	}
	| FOR Name2 IN ExpList { opFor(yylex, &$$, OP_FORIN, &$2, &$4)	}

Conds:
 	CondList
 	| CondList ELSE Block { opIf(yylex, &$$, nil, &$1, &$3) }

CondList:
	Cond
	| CondList ELSEIF Cond { opIf(yylex, &$$, nil, &$1, &$3) }

Cond:
	Exp THEN Block { opIf(yylex, &$$, &$1, &$3, nil) }

LastStat:
	BREAK { op0(yylex, &$$, OP_BREAK, &$1) }
	| CONTINUE { op0(yylex, &$$, OP_CONTINUE, &$1) }
	| RETURN { op1(yylex, &$$, OP_RETURN, nil) }
	| RETURN ExpList { op1(yylex, &$$,OP_RETURN, &$2) }

Binding:
    LOCAL NameList { opLocal(yylex, &$$, &$2, nil) }
    | LOCAL NameList '=' ExpList { opLocal(yylex, &$$, &$2, &$4) }
	| LOCAL FUNCTION NAME FuncBody {
		bindFuncName(yylex, &$4, nil, $3.token.image)
		var tmp yySymType
		nameAppend(yylex, &tmp, &$3, nil)		
		opLocal(yylex, &$$, &tmp, &$4)
	}
	| CLOSURE '(' NameList ')' { opClosure(yylex, &$$, &$3) }
	| CLOSURE NameList { opClosure(yylex, &$$, &$2) }

SetList:
	Var
	| SetList ',' Var { opExpList(yylex, &$$, &$1, &$3) }

FuncName:
	DottedName

DottedName:
	NAME { opVar(&$$, &$1) }
	| DottedName '.' NAME { 
		opValueExt(&$3, $3.token.image)
		op2(yylex, &$$, OP_MEMBER, &$1, &$3)
	}

Exp:
	NIL { opValue(yylex, &$$) }
	| TRUE { opValue(yylex, &$$) }
	| FALSE { opValue(yylex, &$$) }
	| NUMBER { opValue(yylex, &$$) }
	| STRING { opValue(yylex, &$$) }
	| MORE { opVar(&$$, &$1) }
	| FuncDef
	| PrefixExp
	| Tableconstructor
	| Arrayconstructor
	| NOT Exp { op1(yylex, &$$, OP_NOT, &$2) }
	| '#' Exp { op1(yylex, &$$, OP_LEN, &$2) }
	| '-' Exp %prec'*' { op1(yylex, &$$, OP_NSIGN, &$2) }
	| Exp OR Exp { op2(yylex, &$$, OP_OR, &$1, &$3) }
	| Exp AND Exp { op2(yylex, &$$, OP_AND, &$1, &$3) }	
	| Exp '<' Exp { op2(yylex, &$$, OP_LT, &$1, &$3) }
	| Exp '>' Exp { op2(yylex, &$$, OP_GT, &$1, &$3) }
	| Exp SLTEQ Exp { op2(yylex, &$$, OP_LTEQ, &$1, &$3) }
	| Exp SGTEQ Exp { op2(yylex, &$$, OP_GTEQ, &$1, &$3) }
	| Exp SEQ Exp { op2(yylex, &$$, OP_EQ, &$1, &$3) }
	| Exp SNOTEQ Exp { op2(yylex, &$$, OP_NOTEQ, &$1, &$3) }
	| Exp STRADD Exp { op2(yylex, &$$, OP_STRADD, &$1, &$3) }
	| Exp '-' Exp { op2(yylex, &$$, OP_SUB, &$1, &$3) }
	| Exp '+' Exp { op2(yylex, &$$, OP_ADD, &$1, &$3) }
	| Exp '*' Exp { op2(yylex, &$$, OP_MUL, &$1, &$3) }
	| Exp '/' Exp { op2(yylex, &$$, OP_DIV, &$1, &$3) }
	| Exp '^' Exp { op2(yylex, &$$, OP_PMUL, &$1, &$3) }
	| Exp '%' Exp { op2(yylex, &$$, OP_MOD, &$1, &$3) }

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
	PrefixExp Args { op2(yylex, &$$, OP_CALL, &$1, &$2) }

Args:
	'(' ')'
	| '(' ExpList ')' { $$ = $2 }

FuncDef:
	FUNCTION FuncBody { $$ = $2 }

FuncBody:
	ParamDefList Block END { opFunc(yylex, &$$, &$1, &$2) }

ParamDefList:
	'(' ParDefList ')' { $$ = $2 }

ParDefList:
	NameList
	| MORE { nameAppend(yylex, &$$, &$1, nil) }
	| NameList "," MORE { nameAppend(yylex, &$$, &$1, &$3) }
	| { nameAppend(yylex, &$$, nil, nil) }
	;

Name2:
	NAME { nameAppend(yylex, &$$, &$1, nil) }
	| NAME "," NAME { nameAppend(yylex, &$$, &$1, &$3) }

NameList:
	NAME { nameAppend(yylex, &$$, &$1, nil) }
	| NameList "," NAME { nameAppend(yylex, &$$, &$1, &$3) }

Arrayconstructor:
	'[' ']' { op1(yylex, &$$, OP_ARRAY, nil) }
	| '[' ExpList ']' { op1(yylex, &$$, OP_ARRAY, &$2) }
	| '[' ExpList ',' ']' { op1(yylex, &$$, OP_ARRAY, &$2) }

ExpList:
	Exp { opExpList(yylex, &$$, &$1, nil) }
	| ExpList ',' Exp { opExpList(yylex, &$$, &$1, &$3) }

Tableconstructor:
	'{' '}' { op1(yylex, &$$, OP_TABLE, nil) }
	| '{' FieldList '}' { op1(yylex, &$$, OP_TABLE, &$2) }
	| '{' FieldList ',' '}' { op1(yylex, &$$, OP_TABLE, &$2) }

FieldList:
	Field
	| FieldList ',' Field { opAppend(yylex, &$$, &$1, &$3) }

Field:
	NAME '=' Exp { 
		opValueExt(&$1, $1.token.image)
		op2(yylex, &$$, OP_FIELD, &$1, &$3)
	}
	| '[' Exp ']' '=' Exp { op2(yylex, &$$, OP_FIELD, &$2, &$5) }

%%