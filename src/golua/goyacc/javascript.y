%{
package goyacc

import (
     "fmt"
)
func _keepfmt() {
     fmt.Println("hi")
}
%}

%token FUNCTION
%token IF
%token ELSE
%token WHILE
%token FOR
%token IN
%token BREAK
%token CONTINUE
%token VAR
%token NEW
%token DELETE
%token THIS
%token TRUE
%token FALSE
%token NULL
%token WITH
%token RETURN

%token NUMBER
%token STRING
%token NAME
%token SLTEQ
%token SGTEQ
%token SEQ
%token SNOTEQ
%token SLSHIFT
%token SRSHIFT
%token SEQ3
%token SADD2
%token SSUB2
%token SADDASS
%token SSUBASS
%token SMULASS
%token SDIVASS
%token SOR
%token SAND

%left SOR
%left SAND
%left '<' '>' SLTEQ SGTEQ SEQ SNOTEQ
%left '+' '-' '*' '/' '^' '%'
%right '!' '~'

%%
Program:
     Element
     | Program Element

Element:
     FUNCTION NAME  ParameterListOpt CompoundStatement
     | Statement

ParameterListOpt:
     '(' ')'
     | '(' ParameterList ')'

ParameterList:
     NAME
     | ParameterList ',' NAME

CompoundStatement:
     '{' '}'
     | '{' Statements '}'

Statements:
     Statement
     | Statements Statement

Statement:
     ';'
     | IF Condition Statement
     | IF Condition Statement ELSE Statement
     | WHILE Condition Statement
     | ForParen ';' Expression ';' Expression ')' Statement
     | ForBegin ';' Expression ';' Expression ')' Statement
     | ForBegin IN Expression ')' Statement
     | BREAK ';'
     | CONTINUE ';'
     | WITH '(' Expression ')' Statement
     | RETURN ';'
     | RETURN Expression ';'
     | CompoundStatement
     | VariablesOrExpression ';'

Condition:
     '(' Expression ')'

ForParen:
     FOR '('

ForBegin:
     ForParen VariablesOrExpression

VariablesOrExpression:
     VAR Variables
     | Expression

Variables:
     Variable
     | Variables ',' Variable

Variable:
     NAME
     | NAME '=' AssignmentExpression

Expression:
     AssignmentExpression
     | Expression ',' AssignmentExpression 

AssignmentExpression:
     ConditionalExpression
     | AssignmentExpression AssignmentOperator ConditionalExpression

ConditionalExpression:
     OrExpression
     | OrExpression '?' AssignmentExpression ':' AssignmentExpression

OrExpression:
     AndExpression
     | AndExpression SOR OrExpression

AndExpression:
     BitwiseOrExpression
     | BitwiseOrExpression SAND AndExpression

BitwiseOrExpression:
     BitwiseXorExpression
     | BitwiseXorExpression '|' BitwiseOrExpression

BitwiseXorExpression:
     BitwiseAndExpression
     | BitwiseAndExpression '^' BitwiseXorExpression

BitwiseAndExpression:
     EqualityExpression
     | EqualityExpression '&' BitwiseAndExpression

EqualityExpression:
     RelationalExpression
     | RelationalExpression EqualityualityOperator EqualityExpression

RelationalExpression:
     ShiftExpression
     | RelationalExpression RelationalationalOperator ShiftExpression

ShiftExpression:
     AdditiveExpression
     | ShiftExpression ShiftOperator AdditiveExpression

AdditiveExpression:
     MultiplicativeExpression
     | MultiplicativeExpression '+' AdditiveExpression
     | MultiplicativeExpression '-' AdditiveExpression

MultiplicativeExpression:
     | UnaryExpression
     | MultiplicativeExpression MultiplicativeOperator UnaryExpression

UnaryExpression:
     MemberExpression
     | UnaryOperator UnaryExpression
     | '-' UnaryExpression
     | IncrementOperator MemberExpression
     | MemberExpression IncrementOperator
     | NEW Constructor
     | DELETE MemberExpression

Constructor:
     THIS '.' ConstructorCall
     | ConstructorCall

ConstructorCall:
     NAME
     | NAME '(' ')'
     | NAME '(' ArgumentList ')'
     | NAME '.' ConstructorCall

MemberExpression:
     PrimaryExpression
     | PrimaryExpression '.' MemberExpression
     | PrimaryExpression '[' Expression ']'
     | PrimaryExpression '(' ')'
     | PrimaryExpression '(' ArgumentList ')'

ArgumentList:
     AssignmentExpression
     | AssignmentExpression ',' ArgumentList

PrimaryExpression:
     '(' Expression ')'
     | NAME
     | NUMBER { opValue(yylex, &$$) }
     | STRING { opValue(yylex, &$$) }
     | FALSE { opValue(yylex, &$$) }
     | TRUE { opValue(yylex, &$$) }
     | NULL { opValue(yylex, &$$) }
     | THIS
     | TableNew
     | ArrayNew
     | FUNCTION ParameterListOpt CompoundStatement

RelationalationalOperator:
     '<'
     | SLTEQ
     | '>'
     | SGTEQ
     | SEQ
     | SNOTEQ

ShiftOperator:
     SLSHIFT
     | SRSHIFT

UnaryOperator:
     '!'
     | '~'

MultiplicativeOperator:
     '*'
     | '/'
     | '%'
     | '^'

EqualityualityOperator:
     SEQ3

AssignmentOperator:
     '='
     | SADDASS
     | SSUBASS
     | SMULASS
     | SDIVASS

IncrementOperator:
     SADD2
     | SSUB2

ArrayNew:
     '[' ']'
     | '[' ExpList ']'

ExpList:
     Expression
     | ExpList ',' Expression

TableNew:
     '{' '}'
     | '{' FieldList '}' { fmt.Println("TableNew") }

FieldList:
     Field
     | FieldList ',' Field

Field:
     NAME ':' Expression

%%