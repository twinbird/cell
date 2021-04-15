%{
package main
%}
%union {
  num  float64
  expr *Expression
  str  string
  axis string
  stmt *Statement
}
%type<stmt> program stmt
%type<expr> expr 
%token<num> NUMBER 
%token<str> STRING
%token<token> LF '[' ']'

%%
program
  : stmt { yylex.(*Lexer).ast = $1 }

stmt
  : expr LF { $$ = NewExpressionStatement($1) }

expr
  : NUMBER { $$ = NewNumberExpression($1) }
  | STRING { $$ = NewStringExpression($1) }
  | '[' expr ']' { $$ = NewCellReferExpression($2) }
%%

