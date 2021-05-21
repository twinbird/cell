# Tutorial for the programming language cell

cell is a command and language for reading and writing Excel files (xlsx format).

Influenced by awk, perl, etc.

## Install

You can get from the [GitHub](https://github.com/twinbird/cell/releases).

## Hello, world

The following command will create greeting.xlsx with "Hello, world" in A1 cell.

```
$ cell -to greeting.xlsx '["A1"]="Hello, world";'
```

Or, if you want to greet people via standard output, do the following

```
$ cell 'puts("Hello, world");' # => "Hello, world"
```

cell requires a semicolon or line break at the end of the sentence.

However, since a new line is automatically inserted at the end of the program, the above command is equivalent to the following

```
$ cell 'puts("Hello, world")'  # => Hello, world
```

## Let's get some input

If you want to input from an Excel file, you can use the 'from' option.

The following will display the contents of cell A1 of users.xlsx to the console.

```
$ cell -from users.xlsx 'puts(["A1"])'
```

If you want to get the standard input, you can use gets().

```
$ echo "Hello, world" | cell -to greeting.xlsx '["A1"]=gets()'
```

## Value/Variable/Expression

cell can use only string type and 64-bit floating point number type.

Variables do not need to be declared.

When it appears in the program, it is initialized and prepared with an empty string.

There are only two scopes for variables: global and function.

There are no arrays or other data structures.

You can use the Excel data structure instead.

### Accessing Excel data

You can access the A1 cell of the active sheet in an open Excel book by enclosing the string in \[ and \], as in \["A1"\].

The active sheet can be get and set with the @special variable.

```
$ cell 'puts(@)' #=> "Sheet1"
```

```
$ cell '@="Sheet2";puts(@)' #=> "Sheet2"
```

If you set a sheet name that has not been created, that sheet will be created.

### Expression/Operator

cell depends on the operator to interpret the value.

```
one = 1
two = 2
puts(one + two) # => 3(add number)
puts(one . two)  # => 12(string concat)
```

Most of the operators are all common in other languages.

For details on operators, see [Operators](#Operators).

However, the increment operator may seem rather quirky.

These operators interpret the value of the variable as a number and increment it.

However, if it can be interpreted as a column number in Excel, it will be incremented as a string of column numbers.

```
cell 'v=1;v++;puts(v);' # => 2
cell 'v="a";v++;puts(v);' #=> B
cell 'v="Z";v++;puts(v);' #=> AA
cell 'v="string";v++;puts(v);' #=> 1
```

## Branch

In cell, the if statement is used for branching.

As in many other languages, you can also use else.

```
$ cell 'if (gets()){ puts("true"); }else{ puts("false");}'
```

As is common in languages whose ancestor is C, the {} in the block is not necessary when the branch is a single statement.

```
$ cell 'if (gets()) puts("true"); else puts("false");'
```

The if statement will be false if the value in "()" is an empty string or a numeric value of 0, and true otherwise.

## Loop

There are three types of loop structures in cell: "while", "do-while", and "for".

As with the "if" statement, if the content of the block body is a single statement, the {} in the block is unnecessary.

### while

```
$ cell 'i=0;while(i<3){puts(i);i++;}'
# => 0
# => 1
# => 2
```

### do-while

```
$ cell 'do { puts("output this text"); } while(0);'
# => output this text
```

### for

```
$ cell 'for(i=0; i<3;i++) { puts(i); }'
# => 0
# => 1
# => 2
```

## Function definition

The function must be defined before it is called.

Within a function, the scope of a variable is different.

```
# return x + y
function add(x, y) {
  return x + y;
}
puts(add(1, 2));
# => 3
```

## Comment

\# to the end of the line is a comment.

## Example

Seeing is believing, right?

### Create an Excel table with a list of users

Lists the users and home directories of the system.

```
cat /etc/passwd | cell -n -F ":" -to users.xlsx '["A".NR]=$1;["B".NR]=$6'
```

### Register product information to the PostgreSQL

Create SQL from an Excel table with a header in the first row.

```
cell -from items.xlsx -N -s 2 'puts("INSERT INTO items(name, value) (" . ["A".NER] . "," . ["B".NER] .");")' | psql mydb
```

### Make business cards for members of your department.

In Japan, business cards are sometimes created using Excel.

Excel is often used as a tool to create forms.(believe it?)

Let's use the Excel book template.xlsx that contains the template sheet (template) and put it into the D5 cell of the sheet where we will write the name.

```
cat member.txt | cell -from template.xlsx -to business_cards.xlsx -n 'copy("template", $1);@=$1;["D5"]=$1;'
```

## Quick Reference

### Operators

#### Assignment/Reference Operators

| Operator | Feature |
| --------|------|
| \= | Assigning to a variable |
| \[string\] | Refers to the value of a cell |
| \[string\] \= | Set the value to a cell |

#### Numerical operator

| Operator | Feature |
| --------|------|
| + | Addition | 
| - | Subtraction |
| * | Multiplication | 
| / | Division |
| % | Modulo |
| ** | Power |
| +(unary ) | Interpret strings as numbers | 
| -(unary) | Interpret strings as numbers and reverse sign |
| += | Add and assignment |
| -= | Subtract and assignment |
| /= | Division and assignment |
| *= | Multiplication and assignment  |
| %= | Modulo and assignment |
| **= | Power and assignment |

#### String operator

| Operator | Feature |
| --------|------|
| . | Concat |

#### Numeric comparison operator

| Operator | Feature |
| --------|------|
| < | Less than |
| > | Greater than |
| <= | Less than or equal |
| >= | Greater than or equal |
| == | Equal |
| != | Not equal |

#### String comparison operator

| Operator | Feature |
| --------|------|
| eq | Equal |
| ne | Not equal |
| ~ | Match |
| !~ | Not match |

#### Column Number comparison operator

| Operator | Feature |
| --------|------|
| lt | Interpreted as Excel column numbers. Less than |
| le | Interpreted as Excel column numbers. Less than or equal |
| gt | Interpreted as Excel column numbers. Greater than |
| ge | Interpreted as Excel column numbers. Greater than or equal |

#### Logical operator

| Operator | Feature |
| --------|------|
| && | Logical and |
| \|\| | Logical add |
| ! | Logical not |

#### Increment/Decrement operators

The increment/decrement operator is an operator that is applied to variables.

There is a prefix operator and a postfix operator.

These operators behave differently depending on the variable and the value of the variable.

 * Increment/decrement for numeric values
 * If it is a string, it will be interpreted as a number and incremented/decremented.
 * In the case of a column number string, change to the next/previous column number
 * In the case of "@" variables, change to the next/previous sheet

### Special Variables

| Variable | Feature |
| -----|-----|
| @ | Active sheet name |
| FS | Field separator for standard input(default is space or tab) |
| OFS | Field separator for standard output(default is space) |
| RS | Record separator for standard input(default is \\n) |
| ORS | Record separator for standard output(default is \\n) |
| NR | Number of lines imported from standard input |
| NER | Number of Excel rows shown in the loop process(When using Option N) |
| SER | Number of rows to start an Excel loop process(When using Option N) |
| LR | Last row number of active sheet |
| LC | Last col number of active sheet |
| LCC | Last col character number of active sheet |
| $0 | The previous standard input obtained by gets() |
| $1 | The first field of $0 split by field separator |
| $n | The nth field of $0 split by field separator |
| $_0 | A string matched by the match operator(~) |
| $_1 | The first string captured when matched with match operator(~) |
| $_n | The nth string captured when matched with match operator(~) |

### Command line options

| Option | Feature |
| --------------|------|
| -to | Specify the path of the processed Excel file that will be saved |
| -from | Specify the Excel file to be processed. No overwriting will be done. The default is an empty book containing only Sheet1. |
| -f | Read the Cell program source from the file program-file, instead of from the first command line argument. |
| -F | Use fs for the input field separator (the value of the FS predefined variable). |
| -n | Wrap your script inside while(gets()){... ;} loop |
| -N | Wrap your script inside for(NER = SER; NER <= LR; NER++){... ;} loop (NER and SER, LR are predefined variables) |
| -s | Specify the special variable SER(Start Excel Row) (default 1) |
| -S | Specify default active sheet by name |
| -V | Print version information. |
| -h | Show this help |


### Builtin functions

#### exit(n)

"exit" exits the program with exit code n.

If the "to" option is specified, the Excel file being processed will be output upon exit.

#### abort(n)

"abort" terminates the program with exit code n.

Unlike exit, it does not output the Excel file on exit even if the to option is specified.

#### gets()

"gets" reads and returns a single line of text from standard input up to a newline.

If the end of the input is reached, an empty string is returned.

"gets" sets the value of the special variable $0, $1... $n.

The $0 field contains the entire string read, and the $1 field contains the first field delimited by FS. The same applies to $2 and after.

It also increments the special variable NR each time it reads.

#### puts(s...)

"puts" outputs the string s to the standard output.

When called with no arguments, the contents of $0 are output.

When multiple arguments are specified, the output will be a string concatenated with OFS.

The output string will be suffixed with ORS.

#### head()

Set the special variable "@" to the first sheet of the currently opened Excel book.

#### tail()

Set the special variable @ to the last sheet of the currently opened Excel book.

#### rename(old, new)

Change the sheet name "old" to "new".

#### exist(sheetname)

Returns 1 if the currently opened Excel book has a sheet with the sheet name, and 0 otherwise.

#### count()

Returns the number of sheets in the currently open Excel book.

#### delete(sheetname)

Deletes the "sheetname" sheet.

#### copy(from, to)

Copy the "from" sheet with the name "to".

#### srand(n)

Sets the seed of the pseudo-random number.

Use "n" as the new seed for the random number generator.  If no "n" is provided, use the current time.

#### rand()

Returns a pseudo-random number between 0 and 1.

#### floor(n)

Returns the value of n with the decimal point truncated.

#### ceil(n)

Returns the value of n rounded up to the nearest whole number.

#### round(n)

Returns the value of n rounded to the nearest whole number.

## In the end

Thank you DeepL.

