%{
package goyacc
%}

%token NUMBER
%token NAME
%token EOF

%%

Chunk:
	Block
	EOF

Block:
	'[' NUMBER ']'
	;

%%