@echo off
cell.exe -to "thisFileIsCreated.xlsx" "exit(1)"
if %ERRORLEVEL% neq 1 (
  @echo on
  echo "exit status is wrong. want 1, but got " %ERRORLEVEL%
  exit /B 1
)

cell.exe -to "thisFileIsNotCreated.xlsx" "abort(2)"
if %ERRORLEVEL% neq 2 (
  @echo on
  echo "exit status is wrong. want 2, but got" %ERRORLEVEL%
  exit /B 1
)

if exist "thisFileIsNotCreated.xlsx" (
  @echo on
  echo "the file specified by 'to' option exists even through program aborted by abort()."
  exit /B 1
)
 
cell.exe -f test/prog.cell
if %ERRORLEVEL% neq 4 (
  @echo on
  echo "-f option could not working"
  exit /B 1
)

cell.exe "gets();ret+=$1;gets();ret+=$1;exit(ret);" test/data1.txt test/data2.txt
if %ERRORLEVEL% neq 6 (
  @echo on
  echo "file args could not working"
  exit /B 1
)

rem cell.exe -F ":" "exit(FS eq ':')"
rem if %ERRORLEVEL neq 1 (
rem   @echo on
rem   echo "option -F could not working"
rem   exit /B 1
rem )

cell.exe -s 3 "exit(SER)"
if %ERRORLEVEL% neq 3 (
  @echo on
  echo "option -s could not working"
  exit /B 1
)

rem cell.exe -S Sheet2 "exit(@ eq 'Sheet2')"
rem if %ERRORLEVEL% neq 1 (
rem   @echo on
rem   echo "option -S could not working"
rem   exit /B 1
rem )

cell.exe -f test/read.cell test/data1.txt test/data2.txt
if %ERRORLEVEL% neq 6 (
  @echo on
  echo "option -f with file specify could not working"
  exit /B 1
)
