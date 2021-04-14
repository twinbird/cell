%{
package main
%}
%union {
  num  float64
  expr *Expression
  str  string
}
%type<expr> program expr 
%token<num> NUMBER 
%token<str> STRING
%token<token> LF

%%
program
  : expr LF { yylex.(*Lexer).ast = $1 }

expr
  : NUMBER { $$ = NewNumberExpression($1) }
  | STRING { $$ = NewStringExpression($1) }
%%

