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
  params *ParamList
}
%type<stmts>  program stmts
%type<stmt>   stmt
%type<expr>   expr funcCall 
%type<args>   argList
%type<params> paramList
%token<num>   NUMBER 
%token<str>   STRING
%token<token> LF '[' ']' '(' ')' ',' '=' NUMEQ NUMNE '<' NUMLE '>' NUMGE STREQ STRNE COLLT COLLE COLGT COLGE '.' '+' '-' '/' '*' '%' POW AND OR '!' ADD_ASSIGN SUB_ASSIGN MUL_ASSIGN DIV_ASSIGN MOD_ASSIGN POW_ASSIGN '~' NOT_MATCH IF ELSE '{' '}' WHILE CONCAT_ASSIGN BREAK CONTINUE INC DEC DO FOR FUNCTION RETURN
%token<ident> IDENT
%left '=' ADD_ASSIGN SUB_ASSIGN MUL_ASSIGN DIV_ASSIGN MOD_ASSIGN POW_ASSIGN CONCAT_ASSIGN
%left AND OR '!'
%left NUMEQ NUMNE '<' NUMLE '>' NUMGE STREQ STRNE COLLT COLLE COLGT COLGE
%left '.' '+' '-'
%left '/' '*' '%'
%left '~' NOT_MATCH
%right MINUS PLUS
%right POW
%left INC DEC
%right PREINC PREDEC
%left '(' ')'
%nonassoc THEN
%nonassoc ELSE

%%
program
  : stmts { yylex.(*Lexer).ast = $$ }

stmts
  : stmt { $$ = NewStatements($1) }
  | stmts stmt { $$ = $1.appendStatement($2) }

stmt
  : LF { $$ = NewBlankStatement() }
  | expr LF { $$ = NewExpressionStatement($1) }
  | IF '(' expr ')' stmt %prec THEN { $$ = NewIfStatement($3, $5) }
  | IF '(' expr ')' stmt ELSE stmt { $$ = NewIfElseStatement($3, $5, $7) }
  | '{' stmts '}' { $$ = NewBlockStatement($2) }
  | WHILE '(' expr ')' stmt { $$ = NewWhileStatement($3, $5) }
  | DO stmt WHILE '(' expr ')' LF { $$ = NewDoWhileStatement($2, $5) }
  | FOR '(' expr LF expr LF expr ')' stmt { $$ = NewForStatement($3, $5, $7, $9) }
  | BREAK LF { $$ = NewBreakStatement() }
  | CONTINUE LF { $$ = NewContinueStatement() }
  | FUNCTION IDENT '(' paramList ')' stmt { $$ = NewFunctionDefineStatement($2, $4, $6) }
  | RETURN LF { $$ = NewReturnStatement(NewStringExpression("")) }
  | RETURN expr LF { $$ = NewReturnStatement($2) }

expr
  : NUMBER { $$ = NewNumberExpression($1) }
  | STRING { $$ = NewStringExpression($1) }
  | '[' expr ']' { $$ = NewCellReferExpression($2) }
  | '[' expr ']' '=' expr { $$ = NewCellAssignExpression($2, $5) }
  | '[' expr ']' ADD_ASSIGN expr { $$ = NewAddCellAssignExpression($2, $5) }
  | '[' expr ']' SUB_ASSIGN expr { $$ = NewSubCellAssignExpression($2, $5) }
  | '[' expr ']' MUL_ASSIGN expr { $$ = NewMulCellAssignExpression($2, $5) }
  | '[' expr ']' DIV_ASSIGN expr { $$ = NewDivCellAssignExpression($2, $5) }
  | '[' expr ']' MOD_ASSIGN expr { $$ = NewModCellAssignExpression($2, $5) }
  | '[' expr ']' POW_ASSIGN expr { $$ = NewPowCellAssignExpression($2, $5) }
  | '[' expr ']' CONCAT_ASSIGN expr { $$ = NewConcatCellAssignExpression($2, $5) }
  | '[' expr ']' INC { $$ = NewIncrementCellExpression($2) }
  | INC '[' expr ']' %prec PREINC { $$ = NewPreIncrementCellExpression($3) }
  | '[' expr ']' DEC { $$ = NewDecrementCellExpression($2) }
  | DEC '[' expr ']' %prec PREDEC { $$ = NewPreDecrementCellExpression($3) }
  | IDENT { $$ = NewVarReferExpression($1) }
  | IDENT '=' expr { $$ = NewVarAssignExpression($1, $3) }
  | IDENT ADD_ASSIGN expr { $$ = NewAddAssignExpression($1, $3) }
  | IDENT SUB_ASSIGN expr { $$ = NewSubAssignExpression($1, $3) }
  | IDENT MUL_ASSIGN expr { $$ = NewMulAssignExpression($1, $3) }
  | IDENT DIV_ASSIGN expr { $$ = NewDivAssignExpression($1, $3) }
  | IDENT MOD_ASSIGN expr { $$ = NewModAssignExpression($1, $3) }
  | IDENT POW_ASSIGN expr { $$ = NewPowAssignExpression($1, $3) }
  | IDENT CONCAT_ASSIGN expr { $$ = NewConcatAssignExpression($1, $3) }
  | IDENT INC { $$ = NewIncrementExpression($1) }
  | INC IDENT %prec PREINC { $$ = NewPreIncrementExpression($2) }
  | IDENT DEC { $$ = NewDecrementExpression($1) }
  | DEC IDENT %prec PREDEC { $$ = NewPreDecrementExpression($2) }
  | funcCall
  | expr NUMEQ expr { $$ = NewNumberEQExpression($1, $3) }
  | expr NUMNE expr { $$ = NewNumberNEExpression($1, $3) }
  | expr '<' expr { $$ = NewNumberLTExpression($1, $3) }
  | expr NUMLE expr { $$ = NewNumberLEExpression($1, $3) }
  | expr '>' expr { $$ = NewNumberGTExpression($1, $3) }
  | expr NUMGE expr { $$ = NewNumberGEExpression($1, $3) }
  | expr STREQ expr { $$ = NewStringEQExpression($1, $3) }
  | expr STRNE expr { $$ = NewStringNEExpression($1, $3) }
  | expr COLLT expr { $$ = NewColNumberLTExpression($1, $3) }
  | expr COLLE expr { $$ = NewColNumberLEExpression($1, $3) }
  | expr COLGT expr { $$ = NewColNumberGTExpression($1, $3) }
  | expr COLGE expr { $$ = NewColNumberGEExpression($1, $3) }
  | expr '.' expr { $$ = NewStringConcatExpression($1, $3) }
  | expr '+' expr { $$ = NewNumberAddExpression($1, $3) }
  | expr '-' expr { $$ = NewNumberSubExpression($1, $3) }
  | expr '*' expr { $$ = NewNumberMulExpression($1, $3) }
  | expr '/' expr { $$ = NewNumberDivExpression($1, $3) }
  | expr '%' expr { $$ = NewNumberModuloExpression($1, $3) }
  | expr '~' expr { $$ = NewStringMatchExpression($1, $3) }
  | expr NOT_MATCH expr { $$ = NewStringNotMatchExpression($1, $3) }
  | expr POW expr { $$ = NewNumberPowerExpression($1, $3) }
  | expr AND expr { $$ = NewLogicalAndExpression($1, $3) }
  | expr OR expr { $$ = NewLogicalOrExpression($1, $3) }
  | '!' expr { $$ = NewLogicalNotExpression($2) }
  | '(' expr ')' { $$ = $2 }
  | '-' expr %prec MINUS { $$ = NewMinusExpression($2) }
  | '+' expr %prec PLUS { $$ = NewPlusExpression($2) }

funcCall
  : IDENT '(' ')' { $$ = NewFuncCallExpression($1, NewEmptyArgList()) }
  | IDENT '(' argList ')' { $$ = NewFuncCallExpression($1, $3) }

argList
  : expr { $$ = NewArgList($1) }
  | expr ',' argList { $$ = $3.appendArg($1) }

paramList
  : { $$ = NewEmptyParamList() }
  | IDENT { $$ = NewParamList($1) }
  | IDENT ',' paramList { $$ = $3.appendParam($1) }
%%
