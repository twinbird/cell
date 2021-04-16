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
  prim  Primitive
  ident string
}
%type<prim>   primitive
%type<stmts>  program stmts
%type<stmt>   stmt
%type<expr>   expr 
%token<num>   NUMBER 
%token<str>   STRING
%token<token> LF '[' ']'
%token<ident> IDENT

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
  : primitive { $$ = NewLiteralExpression($1) }
  | '[' expr ']' { $$ = NewCellReferExpression($2) }
  | '[' expr ']' '=' expr { $$ = NewCellAssignExpression($2, $5) }
  | IDENT { $$ = NewVarReferExpression($1) }
  | IDENT '=' expr { $$ = NewVarAssignExpression($1, $3) }

primitive
  : NUMBER { $$ = NewNumberValue($1) }
  | STRING { $$ = NewStringValue($1) }
%%
