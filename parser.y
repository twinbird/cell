%{
package main
%}
%union {
  num float64
  expr *Expression
}
%type<expr> program expr 
%token<num> NUMBER 
%token<token> LF

%%
program
  : expr LF { yylex.(*Lexer).ast = $1 }

expr
  : NUMBER { $$ = NewNumberExpression($1) }
%%

