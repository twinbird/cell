%{
package main
%}
%union {
  num   float64
  expr  *Expression
  str   string
  axis  string
  stmt  *Statement
  stmts *Statements
}
%type<stmts> program stmts
%type<stmt> stmt
%type<expr> expr 
%token<num> NUMBER 
%token<str> STRING
%token<token> LF '[' ']'

%%
program
  : stmts { yylex.(*Lexer).ast = $$ }

stmts
  : stmt { $$ = NewStatements($1) }
  | stmts stmt { $$ = $1.appendStatement($2) }

stmt
  : LF { $$ = NewBlankStatement() }
  | expr LF { $$ = NewExpressionStatement($1) }

expr
  : NUMBER { $$ = NewNumberExpression($1) }
  | STRING { $$ = NewStringExpression($1) }
  | '[' expr ']' { $$ = NewCellReferExpression($2) }
  | '[' expr ']' '=' expr { $$ = NewCellAssignExpression($2, $5) }
%%

