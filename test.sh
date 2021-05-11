#!/usr/bin/bash

./cell -to 'thisFileIsCreated.xlsx' 'exit(1)'
if [[ $? -ne 1 ]]; then
  echo "exit status is wrong. want 1, but got $?"
  exit 1
fi
if [[ ! -e 'thisFileIsCreated.xlsx' ]]; then
  echo 'file is not created when aborted by exit()'
  exit 1
fi

./cell -to 'thisFileIsNotCreated.xlsx' 'abort(2)'
if [[ $? -ne 2 ]]; then
  echo "exit status is wrong. want 2, but got $?"
  exit 1
fi
if [[ -e 'thisFileIsNotCreated.xlsx' ]]; then
  echo 'the file specified by 'to' option exists even through program aborted by abort().'
fi

./cell -f test/prog.cell
if [[ $? -ne 4 ]]; then
  echo "-f option could not working"
  exit 1
fi

./cell 'gets();ret+=$1;gets();ret+=$1;exit(ret);' test/data1.txt test/data2.txt
if [[ $? -ne 6 ]]; then
  echo 'file args could not working'
  exit 1
fi

./cell -F ":" 'exit(FS eq ":")'
if [[ $? -ne 1 ]]; then
  echo 'option -F could not working'
  exit 1
fi

./cell -s 3 'exit(SER)'
if [[ $? -ne 3 ]]; then
  echo 'option -s could not working'
  exit 1
fi

./cell -S Sheet2 'exit(@ eq "Sheet2")'
if [[ $? -ne 1 ]]; then
  echo 'option -S could not working'
  exit 1
fi
