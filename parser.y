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
  ident string
  args *ArgList
}
%type<stmts>  program stmts
%type<stmt>   stmt
%type<expr>   expr funcCall 
%type<args>   argList
%token<num>   NUMBER 
%token<str>   STRING
%token<token> LF '[' ']' '(' ')' ','
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
  : NUMBER { $$ = NewNumberExpression($1) }
  | STRING { $$ = NewStringExpression($1) }
  | '[' expr ']' { $$ = NewCellReferExpression($2) }
  | '[' expr ']' '=' expr { $$ = NewCellAssignExpression($2, $5) }
  | IDENT { $$ = NewVarReferExpression($1) }
  | IDENT '=' expr { $$ = NewVarAssignExpression($1, $3) }
  | funcCall

funcCall
  : IDENT '(' ')' { $$ = NewFuncCallExpression($1, NewEmptyArgList()) }
  | IDENT '(' argList ')' { $$ = NewFuncCallExpression($1, $3) }

argList
  : expr { $$ = NewArgList($1) }
  | expr ',' argList { $$ = $3.appendArg($1) }
%%
