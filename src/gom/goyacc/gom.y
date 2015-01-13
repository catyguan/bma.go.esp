%{
package goyacc

%}

%token OBJECT
%token SERVICE
%token STRUCT
%token TRUE
%token FALSE
%token NIL

%token NUMBER
%token STRING
%token NAME

%%
GOM:
	Body { endGOM(yylex, &$$) }

Body:
	DefList	

DefList:
	Def
	| DefList Def

Def:
	DefValue
	| DefStruct
	| DefService
	| DefObject

FullName:
	NAME
	| FullName '.' NAME { nameAppend(yylex, &$$, &$1, &$3) }

AnnoList:
	Anno { annoAppend(yylex, &$$, &$1, nil) }
	| AnnoList Anno { annoAppend(yylex, &$$, &$1, &$2) }

Anno:
	'@' FullName '=' Value { defineAnnotation(yylex, &$$, &$2, &$4) }

DefValue:
	DefSValue { commitValue(yylex, &$$, &$1, nil) }
	| AnnoList DefSValue { commitValue(yylex, &$$, &$2, &$1) }

DefSValue:
    | FullName '=' Value { op2(yylex, &$$, OP_VALUE, &$1, &$3) }

Value:
	NIL 		{ opValue(yylex, &$$) }
	| TRUE 		{ opValue(yylex, &$$) }
	| FALSE 	{ opValue(yylex, &$$) }
	| NUMBER 	{ opValue(yylex, &$$) }
	| STRING	{ opValue(yylex, &$$) }
	| Table
	| Array

DefStruct:
	DefSStruct { commitStruct(yylex, &$$, &$1, nil) }
	| AnnoList DefSStruct { commitStruct(yylex, &$$, &$2, &$1) }

DefSStruct:
	STRUCT FullName StructBody { op2(yylex, &$$, OP_STRUCT, &$2, &$3) }

StructBody:
	'{' '}' { opN(yylex, &$$, OP_STRUCT_BODY, nil) }
	| '{' StructFieldList '}' { opN(yylex, &$$, OP_STRUCT_BODY, &$2) }
	| '{' StructFieldList ',' '}' { opN(yylex, &$$, OP_STRUCT_BODY, &$2) }

StructFieldList:
	StructField { nodeAppend(yylex, &$$, &$1, nil) }
	| StructFieldList ',' StructField { nodeAppend(yylex, &$$, &$1, &$3) }

StructField:
	SStructField { commitStructField(yylex, &$$, &$1, nil) }
	| AnnoList SStructField { commitStructField(yylex, &$$, &$2, &$1) }

SStructField:
	FullName ':' ValType { op2(yylex, &$$, OP_SFIELD, &$1, &$3) }

ValType:
	NAME { op2(yylex, &$$, OP_TYPE, &$1, nil) }
	| NAME '<' ValType '>' { op2(yylex, &$$, OP_TYPE, &$1, &$3) }
	| NAME '<' ValType ',' ValType '>' { 		
		op2(yylex, &$3, OP_TYPE, &$3, &$5)
		op2(yylex, &$$, OP_TYPE, &$1, &$3)
	}

DefService:
	DefSService { commitService(yylex, &$$, &$1, nil) }
	| AnnoList DefSService { commitService(yylex, &$$, &$2, &$1) }

DefSService:
	SERVICE FullName ServiceBody { op2(yylex, &$$, OP_SERVICE, &$2, &$3) }

ServiceBody:
	'{' '}' { opN(yylex, &$$, OP_SERVICE_BODY, nil) }
	| '{' ServiceMethodList '}' { opN(yylex, &$$, OP_SERVICE_BODY, &$2) }
	| '{' ServiceMethodList ',' '}' { opN(yylex, &$$, OP_SERVICE_BODY, &$2) }

ServiceMethodList:
	ServiceMethod { nodeAppend(yylex, &$$, &$1, nil) }
	| ServiceMethodList ',' ServiceMethod { nodeAppend(yylex, &$$, &$1, &$3) }

ServiceMethod:
	SServiceMethod { commitServiceMethod(yylex, &$$, &$1, nil) }
	| AnnoList SServiceMethod { commitServiceMethod(yylex, &$$, &$2, &$1) }

SServiceMethod:
	NAME ParamDefList ':' ValType { op3(yylex, &$$, OP_SMETHOD, &$1, &$2, &$4) }

ParamDefList:	
	'(' ParDefList ')' { opN(yylex, &$$, OP_SM_PARAMS, &$2) }

ParDefList:
	ParDef { nodeAppend(yylex, &$$, &$1, nil) }
	| ParDefList ',' ParDef { nodeAppend(yylex, &$$, &$1, &$3) }
	| { nodeAppend(yylex, &$$, nil, nil) }
	;

ParDef:
	SParDef { commitMethodParam(yylex, &$$, &$1, nil) }
	| AnnoList SParDef { commitMethodParam(yylex, &$$, &$2, &$1) }

SParDef:
	NAME ':' ValType { op2(yylex, &$$, OP_SM_PARAM, &$1, &$3) }

DefObject:
	DefSObject { commitObject(yylex, &$$, &$1, nil) }
	| AnnoList DefSObject { commitObject(yylex, &$$, &$2, &$1) }

DefSObject:
	OBJECT FullName ObjectBody { op2(yylex, &$$, OP_OBJECT, &$2, &$3) }

ObjectBody:
	'{' '}' { opN(yylex, &$$, OP_OBJECT_BODY, nil) }
	| '{' ObjectAttrList '}' { opN(yylex, &$$, OP_SERVICE_BODY, &$2) }
	| '{' ObjectAttrList ',' '}' { opN(yylex, &$$, OP_OBJECT_BODY, &$2) }

ObjectAttrList:
	ObjectAttr { nodeAppend(yylex, &$$, &$1, nil) }
	| ObjectAttrList ',' ObjectAttr { nodeAppend(yylex, &$$, &$1, &$3) }

ObjectAttr:
	ServiceMethod
	| StructField

Array:
	'[' ']' { defineArray(yylex, &$$, nil) }
	| '[' ValueList ']' { defineArray(yylex, &$$, &$2) }
	| '[' ValueList ',' ']' { defineArray(yylex, &$$, &$2) }

ValueList:
	Value { beNode(yylex, &$$, &$1) }
	| ValueList ',' Value { nodeAppend(yylex, &$$, &$1, &$3) }

Table:
	'{' '}' { defineTable(yylex, &$$, nil) }
	| '{' FieldList '}' { defineTable(yylex, &$$, &$2) }
	| '{' FieldList ',' '}' { defineTable(yylex, &$$, &$2) }

FieldList:
	Field { nodeAppend(yylex, &$$, &$1, nil) }
	| FieldList ',' Field { nodeAppend(yylex, &$$, &$1, &$3) }

Field:
	STRING ':' Value { defineField(yylex, &$$, &$1, &$3) }

%%